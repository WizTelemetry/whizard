package compactor

import (
	"fmt"
	"strings"

	"github.com/kubesphere/paodin/pkg/api/monitoring/v1alpha1"
	"github.com/kubesphere/paodin/pkg/controllers/monitoring/resources"
	"github.com/kubesphere/paodin/pkg/controllers/monitoring/resources/storage"
	"github.com/kubesphere/paodin/pkg/util"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	storageDir = "/thanos"

	mainContainerName = "compactor"
	tsdbVolumeName    = "tsdb"
)

func (r *Compactor) statefulSet() (runtime.Object, resources.Operation, error) {
	var sts = &appsv1.StatefulSet{ObjectMeta: r.meta(r.compactor.Name)}
	if err := r.Client.Get(r.Context, client.ObjectKeyFromObject(sts), sts); err != nil {
		if !util.IsNotFound(err) {
			return nil, "", err
		}
	}

	sts.Spec.Selector = &metav1.LabelSelector{
		MatchLabels: r.labels(),
	}

	sts.Spec.Replicas = r.compactor.Spec.Replicas
	sts.Spec.Template.Labels = r.labels()
	sts.Spec.Template.Spec.Affinity = r.compactor.Spec.Affinity
	sts.Spec.Template.Spec.NodeSelector = r.compactor.Spec.NodeSelector

	var container *corev1.Container
	for i := 0; i < len(sts.Spec.Template.Spec.Containers); i++ {
		if sts.Spec.Template.Spec.Containers[i].Name == mainContainerName {
			container = &sts.Spec.Template.Spec.Containers[i]
		}
	}

	needToAppend := false
	if container == nil {
		container = &corev1.Container{
			Name:      mainContainerName,
			Image:     r.compactor.Spec.Image,
			Resources: r.compactor.Spec.Resources,
			Ports: []corev1.ContainerPort{
				{
					Name:          resources.ThanosHTTPPortName,
					ContainerPort: resources.ThanosHTTPPort,
					Protocol:      corev1.ProtocolTCP,
				},
			},
		}

		needToAppend = true
	}

	util.AddVolume(sts, container, r.compactor.Spec.DataVolume, tsdbVolumeName, storageDir)

	if container.LivenessProbe == nil {
		container.LivenessProbe = resources.ThanosDefaultLivenessProbe()
	}

	if container.ReadinessProbe == nil {
		container.ReadinessProbe = resources.ThanosDefaultReadinessProbe()
	}

	container.Resources = r.compactor.Spec.Resources

	if err := r.megerArgs(container); err != nil {
		return nil, "", err
	}

	if needToAppend {
		sts.Spec.Template.Spec.Containers = append(sts.Spec.Template.Spec.Containers, *container)
	}

	return sts, resources.OperationCreateOrUpdate, ctrl.SetControllerReference(r.compactor, sts, r.Scheme)
}

type relabelConfig struct {
	Action        string   `yaml:"action"`
	SourceLablels []string `yaml:"source_labels"`
	Regex         string   `yaml:"regex"`
}

func (r *Compactor) createRelabelConfig() (string, error) {

	namespacedName := strings.Split(r.compactor.Labels[v1alpha1.MonitoringPaodinService], ".")
	svc := &v1alpha1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      namespacedName[1],
			Namespace: namespacedName[0],
		},
	}
	if err := r.Client.Get(r.Context, client.ObjectKeyFromObject(svc), svc); err != nil {
		return "", err
	}

	label := svc.Spec.TenantLabelName
	if len(label) == 0 {
		label = v1alpha1.DefaultTenantLabelName
	}

	regex := ""
	for _, tenant := range r.compactor.Spec.Tenants {
		regex = fmt.Sprintf("%s|^%s$", regex, tenant)
	}

	return util.YamlMarshal([]relabelConfig{
		{Action: "keep",
			SourceLablels: []string{label},
			Regex:         strings.TrimPrefix(regex, "|"),
		},
	})
}

func (r *Compactor) megerArgs(container *corev1.Container) error {

	defaultArgs := []string{"compact", "--wait", fmt.Sprintf("--data-dir=%s", storageDir)}

	storageInstance := &v1alpha1.Storage{}
	namespaceName := strings.Split(r.compactor.Labels[v1alpha1.MonitoringPaodinStorage], ".")
	if err := r.Client.Get(r.Context, client.ObjectKey{Name: namespaceName[1], Namespace: namespaceName[0]}, storageInstance); err != nil {
		return err
	}

	objstoreConfig, err := storage.New(r.BaseReconciler, storageInstance).String()
	if err != nil {
		return err
	}
	defaultArgs = append(defaultArgs, "--objstore.config="+objstoreConfig)

	if r.compactor.Spec.LogLevel != "" {
		defaultArgs = append(defaultArgs, "--log.level="+r.compactor.Spec.LogLevel)
	}
	if r.compactor.Spec.LogFormat != "" {
		defaultArgs = append(defaultArgs, "--log.format="+r.compactor.Spec.LogFormat)
	}
	if r.compactor.Spec.DownsamplingDisable != nil {
		defaultArgs = append(defaultArgs, fmt.Sprintf("--downsampling.disable=%v", r.compactor.Spec.DownsamplingDisable))
	}
	if retention := r.compactor.Spec.Retention; retention != nil {
		if retention.RetentionRaw != "" {
			defaultArgs = append(defaultArgs, fmt.Sprintf("--retention.resolution-raw=%s", retention.RetentionRaw))
		}
		if retention.Retention5m != "" {
			defaultArgs = append(defaultArgs, fmt.Sprintf("--retention.resolution-5m=%s", retention.Retention5m))
		}
		if retention.Retention5m != "" {
			defaultArgs = append(defaultArgs, fmt.Sprintf("--retention.resolution-1h=%s", retention.Retention5m))
		}
	}
	defaultArgs = append(defaultArgs, "--deduplication.replica-label=receive_replica")

	rc, err := r.createRelabelConfig()
	if err != nil {
		return err
	}

	defaultArgs = append(defaultArgs, fmt.Sprintf("--selector.relabel-config=%s", rc))

	for name, value := range r.compactor.Spec.Flags {
		arg := fmt.Sprintf("--%s=%s", name, value)
		replaced := util.ReplaceInSlice(defaultArgs, func(v interface{}) bool {
			return util.GetArgName(v.(string)) == name
		}, arg)

		if !replaced {
			defaultArgs = append(defaultArgs, arg)
		}
	}

	for _, arg := range defaultArgs {

		replaced := util.ReplaceInSlice(container.Args, func(v interface{}) bool {
			return util.GetArgName(v.(string)) == util.GetArgName(arg)
		}, arg)

		if !replaced {
			container.Args = append(container.Args, arg)
		}
	}

	return nil
}
