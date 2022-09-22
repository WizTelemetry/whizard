package query

import (
	"fmt"
	"path/filepath"
	"sort"
	"strconv"

	monitoringv1alpha1 "github.com/kubesphere/whizard/pkg/api/monitoring/v1alpha1"
	"github.com/kubesphere/whizard/pkg/constants"
	"github.com/kubesphere/whizard/pkg/controllers/monitoring/resources"
	"github.com/kubesphere/whizard/pkg/controllers/monitoring/resources/ingester"
	"github.com/kubesphere/whizard/pkg/util"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	// repeatableArgs is the args that can be set repeatedly.
	// An error will occur if a non-repeatable arg is set repeatedly.
	repeatableArgs = []string{
		"--query.replica-label",
		"--selector-label",
		"--endpoint",
		"--endpoint-strict",
		"--store.sd-files",
		"--enable-feature",
	}
	// unsupportedArgs is the args that are not allowed to be set by the user.
	unsupportedArgs = []string{
		// Deprecation
		"--log.request.decision",
		// Deprecation
		"--metadata",
		// Deprecation
		"--rule",
		// Deprecation
		"--store",
		// Deprecation
		"--exemplar",
		// Deprecation
		"--target",
		// Deprecation
		"--store-strict",
		"--grpc-address",
		"--http-address",
	}
)

func (q *Query) deployment() (runtime.Object, resources.Operation, error) {
	var d = &appsv1.Deployment{ObjectMeta: q.meta(q.name())}

	if q.query == nil {
		return d, resources.OperationDelete, nil
	}

	d.Spec = appsv1.DeploymentSpec{
		Replicas: q.query.Spec.Replicas,
		Selector: &metav1.LabelSelector{
			MatchLabels: q.labels(),
		},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: q.labels(),
			},
			Spec: corev1.PodSpec{
				NodeSelector: q.query.Spec.NodeSelector,
				Tolerations:  q.query.Spec.Tolerations,
				Affinity:     q.query.Spec.Affinity,
			},
		},
	}

	proxyConfigVol := corev1.Volume{
		Name: "proxy-config",
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: q.name("proxy-config"),
				},
			},
		},
	}
	d.Spec.Template.Spec.Volumes = append(d.Spec.Template.Spec.Volumes, proxyConfigVol)
	storesConfigVol := corev1.Volume{
		Name: "stores-config",
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: q.name("stores-config"),
				},
			},
		},
	}
	d.Spec.Template.Spec.Volumes = append(d.Spec.Template.Spec.Volumes, storesConfigVol)

	var queryContainer = corev1.Container{
		Name:      "query",
		Image:     q.query.Spec.Image,
		Args:      []string{"query"},
		Resources: q.query.Spec.Resources,
		Ports: []corev1.ContainerPort{
			{
				Protocol:      corev1.ProtocolTCP,
				Name:          constants.GRPCPortName,
				ContainerPort: constants.GRPCPort,
			},
			{
				Protocol:      corev1.ProtocolTCP,
				Name:          constants.HTTPPortName,
				ContainerPort: constants.HTTPPort,
			},
		},
		LivenessProbe:  q.DefaultLivenessProbe(),
		ReadinessProbe: q.DefaultReadinessProbe(),
		VolumeMounts: []corev1.VolumeMount{{
			Name:      storesConfigVol.Name,
			MountPath: configDir,
			ReadOnly:  true,
		}},
	}

	if q.query.Spec.LogLevel != "" {
		queryContainer.Args = append(queryContainer.Args, "--log.level="+q.query.Spec.LogLevel)
	}
	if q.query.Spec.LogFormat != "" {
		queryContainer.Args = append(queryContainer.Args, "--log.format="+q.query.Spec.LogFormat)
	}
	queryContainer.Args = append(queryContainer.Args, "--store.sd-files="+filepath.Join(configDir, storesFile))
	for _, labelName := range q.query.Spec.ReplicaLabelNames {
		queryContainer.Args = append(queryContainer.Args, "--query.replica-label="+labelName)
	}

	var ingesterList monitoringv1alpha1.IngesterList
	if err := q.Client.List(q.Context, &ingesterList,
		client.MatchingLabels(util.ManagedLabelByService(q.Service))); err != nil {

		q.Log.WithValues("ingesterlist", "").Error(err, "")
		return nil, resources.OperationCreateOrUpdate, err
	}
	for _, item := range ingesterList.Items {
		ingesterInstance, err := ingester.New(q.BaseReconciler, &item, nil)
		if err != nil {
			return nil, "", err
		}
		for _, endpoint := range ingesterInstance.GrpcAddrs() {
			queryContainer.Args = append(queryContainer.Args, "--endpoint="+endpoint)
		}
	}

	var storeList monitoringv1alpha1.StoreList
	if err := q.Client.List(q.Context, &storeList,
		client.MatchingLabels(util.ManagedLabelByService(q.Service))); err != nil {

		q.Log.WithValues("storelist", "").Error(err, "")
		return nil, resources.OperationCreateOrUpdate, err
	}
	for _, item := range storeList.Items {
		storeSvcName := util.Join("-", constants.AppNameStore, item.Name, constants.ServiceNameSuffix)
		endpoint := fmt.Sprintf("%s.%s.svc:%d", storeSvcName, item.Namespace, constants.GRPCPort)
		queryContainer.Args = append(queryContainer.Args, "--endpoint="+endpoint)
	}

	// add ruler endpoint to query args
	var rulerList monitoringv1alpha1.RulerList
	if err := q.Client.List(q.Context, &rulerList,
		client.MatchingLabels(util.ManagedLabelByService(q.Service))); err != nil {

		q.Log.WithValues("rulerlist", "").Error(err, "")
		return nil, resources.OperationCreateOrUpdate, err
	}
	for _, item := range rulerList.Items {
		// cancatenate the address instead of calling ruler.GrpcAddrs() to avoid interdependent collisions
		// should be consitent with the logic of ruler.GrpcAddrs()
		var shards int32 = 1
		if item.Spec.Shards != nil && *item.Spec.Shards > 1 {
			shards = *item.Spec.Shards
		}
		for shardSn := 0; shardSn < int(shards); shardSn++ {
			addr := fmt.Sprintf("%s.%s.svc:%d",
				q.QualifiedName(constants.AppNameRuler, item.Name, strconv.Itoa(shardSn), constants.ServiceNameSuffix),
				item.Namespace, constants.GRPCPort)
			queryContainer.Args = append(queryContainer.Args, "--endpoint="+addr)
		}
	}

	for _, flag := range q.query.Spec.Flags {
		arg := util.GetArgName(flag)
		if util.Contains(unsupportedArgs, arg) {
			klog.V(3).Infof("ignore the unsupported flag %s", arg)
			continue
		}

		if util.Contains(repeatableArgs, arg) {
			queryContainer.Args = append(queryContainer.Args, flag)
			continue
		}

		replaced := util.ReplaceInSlice(queryContainer.Args, func(v interface{}) bool {
			return util.GetArgName(v.(string)) == util.GetArgName(flag)
		}, flag)
		if !replaced {
			queryContainer.Args = append(queryContainer.Args, flag)
		}
	}

	sort.Strings(queryContainer.Args[1:])

	var envoyContainer = corev1.Container{
		Name:  "proxy",
		Image: q.query.Spec.Envoy.Image,
		Args: []string{
			"-c",
			filepath.Join(envoyConfigDir, envoyConfigFile),
			// "-l",
			// "debug",
		},
		Resources: q.query.Spec.Envoy.Resources,
		VolumeMounts: []corev1.VolumeMount{{
			Name:      proxyConfigVol.Name,
			MountPath: envoyConfigDir,
			ReadOnly:  true,
		}},
	}

	for _, store := range q.query.Spec.Stores {
		if store.CASecret == nil {
			continue
		}
		secretVol := corev1.Volume{
			Name: "secret-" + store.CASecret.Name,
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: store.CASecret.Name,
				},
			},
		}
		d.Spec.Template.Spec.Volumes = append(d.Spec.Template.Spec.Volumes, secretVol)
		envoyContainer.VolumeMounts = append(envoyContainer.VolumeMounts, corev1.VolumeMount{
			Name:      secretVol.Name,
			ReadOnly:  true,
			SubPath:   store.CASecret.Key,
			MountPath: filepath.Join(envoySecretsDir, store.CASecret.Name, store.CASecret.Key),
		})
	}

	d.Spec.Template.Spec.Containers = append(d.Spec.Template.Spec.Containers, queryContainer)
	d.Spec.Template.Spec.Containers = append(d.Spec.Template.Spec.Containers, envoyContainer)

	return d, resources.OperationCreateOrUpdate, ctrl.SetControllerReference(q.query, d, q.Scheme)
}
