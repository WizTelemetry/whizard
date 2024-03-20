package router

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"path/filepath"
	"sort"

	"github.com/prometheus-operator/prometheus-operator/pkg/k8sutil"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/kubesphere/whizard/pkg/constants"
	"github.com/kubesphere/whizard/pkg/controllers/monitoring/resources"
	"github.com/kubesphere/whizard/pkg/util"
)

var (
	// repeatableArgs is the args that can be set repeatedly.
	// An error will occur if a non-repeatable arg is set repeatedly.
	repeatableArgs = []string{
		"--label",
	}
	// unsupportedArgs is the args that are not allowed to be set by the user.
	unsupportedArgs = []string{
		"--receive.local-endpoint",
		"--http-address",
		"--grpc-address",
	}
)

func (r *Router) deployment() (runtime.Object, resources.Operation, error) {
	var d = &appsv1.Deployment{ObjectMeta: r.meta(r.name())}

	if r.router == nil {
		return d, resources.OperationDelete, nil
	}

	d.Spec = appsv1.DeploymentSpec{
		Replicas: r.router.Spec.Replicas,
		Selector: &metav1.LabelSelector{
			MatchLabels: r.labels(),
		},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: r.labels(),
			},
			Spec: corev1.PodSpec{
				NodeSelector:    r.router.Spec.NodeSelector,
				Tolerations:     r.router.Spec.Tolerations,
				Affinity:        r.router.Spec.Affinity,
				SecurityContext: r.router.Spec.SecurityContext,
			},
		},
	}

	hashringsVol := corev1.Volume{
		Name: "hashrings-config",
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: r.name("hashrings-config"),
				},
			},
		},
	}
	d.Spec.Template.Spec.Volumes = append(d.Spec.Template.Spec.Volumes, hashringsVol)

	if r.router.Spec.ImagePullSecrets != nil && len(r.router.Spec.ImagePullSecrets) > 0 {
		d.Spec.Template.Spec.ImagePullSecrets = r.router.Spec.ImagePullSecrets
	}

	var container = corev1.Container{
		Name:      "receive",
		Image:     r.router.Spec.Image,
		Args:      []string{"receive"},
		Resources: r.router.Spec.Resources,
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
			{
				Protocol:      corev1.ProtocolTCP,
				Name:          constants.RemoteWritePortName,
				ContainerPort: constants.RemoteWritePort,
			},
		},
		LivenessProbe:  r.DefaultLivenessProbe(),
		ReadinessProbe: r.DefaultReadinessProbe(),
		VolumeMounts: []corev1.VolumeMount{{
			Name:      hashringsVol.Name,
			MountPath: configDir,
			ReadOnly:  true,
		}},
		Env: []corev1.EnvVar{
			{
				Name: "POD_NAME",
				ValueFrom: &corev1.EnvVarSource{
					FieldRef: &corev1.ObjectFieldSelector{
						FieldPath: "metadata.name",
					},
				},
			},
		},
	}

	if r.router.Spec.WebConfig != nil {
		secret, _, err := r.webConfigSecret()
		if err != nil {
			return nil, "", err
		}
		hash := md5.New()
		hash.Write(secret.(*corev1.Secret).Data[constants.WhizardWebConfigFile])
		hashStr := hex.EncodeToString(hash.Sum(nil))
		if d.Spec.Template.Annotations == nil {
			d.Spec.Template.Annotations = make(map[string]string)
		}
		d.Spec.Template.Annotations[constants.LabelNameConfigHash] = hashStr

		volumes, volumeMounts := r.BaseReconciler.CreateWebConfigVolumeMount(r.name("web-config"), r.router.Spec.WebConfig)
		d.Spec.Template.Spec.Volumes = append(d.Spec.Template.Spec.Volumes, volumes...)
		container.VolumeMounts = append(container.VolumeMounts, volumeMounts...)

		container.Args = append(container.Args, fmt.Sprintf("--http.config=%s", constants.WhizardWebConfigMountPath+constants.WhizardWebConfigFile))

		if r.router.Spec.WebConfig.HTTPServerTLSConfig != nil {
			container.LivenessProbe = r.DefaultLivenessProbeWithTLS()
			container.ReadinessProbe = r.DefaultReadinessProbeWithTLS()
		}
	}

	if r.router.Spec.LogLevel != "" {
		container.Args = append(container.Args, "--log.level="+r.router.Spec.LogLevel)
	}
	if r.router.Spec.LogFormat != "" {
		container.Args = append(container.Args, "--log.format="+r.router.Spec.LogFormat)
	}
	container.Args = append(container.Args, fmt.Sprintf("--label=%s=\"$(POD_NAME)\"", constants.ReceiveReplicaLabelName))
	container.Args = append(container.Args, "--receive.hashrings-file="+filepath.Join(configDir, hashringsFile))
	if r.router.Spec.ReplicationFactor != nil {
		container.Args = append(container.Args, fmt.Sprintf("--receive.replication-factor=%d", *r.router.Spec.ReplicationFactor))
	}

	if r.Service.Spec.TenantHeader != "" {
		container.Args = append(container.Args, "--receive.tenant-header="+r.Service.Spec.TenantHeader)
	}
	if r.Service.Spec.TenantLabelName != "" {
		container.Args = append(container.Args, "--receive.tenant-label-name="+r.Service.Spec.TenantLabelName)
	}
	if r.Service.Spec.DefaultTenantId != "" {
		container.Args = append(container.Args, "--receive.default-tenant-id="+r.Service.Spec.DefaultTenantId)
	}

	for _, flag := range r.router.Spec.Flags {
		arg := util.GetArgName(flag)
		if util.Contains(unsupportedArgs, arg) {
			klog.V(3).Infof("ignore the unsupported flag %s", arg)
			continue
		}

		if util.Contains(repeatableArgs, arg) {
			container.Args = append(container.Args, flag)
			continue
		}

		replaced := util.ReplaceInSlice(container.Args, func(v interface{}) bool {
			return util.GetArgName(v.(string)) == util.GetArgName(flag)
		}, flag)
		if !replaced {
			container.Args = append(container.Args, flag)
		}
	}

	sort.Strings(container.Args[1:])

	d.Spec.Template.Spec.Containers = append(d.Spec.Template.Spec.Containers, container)

	if len(r.router.Spec.Containers.Raw) > 0 {
		var err error
		r.router.Spec.EmbeddedContainers, err = util.DecodeRawToContainers(r.router.Spec.Containers)
		if err != nil {
			return nil, "", fmt.Errorf("failed to decode containers: %w", err)
		}
		containers, err := k8sutil.MergePatchContainers(d.Spec.Template.Spec.Containers, r.router.Spec.EmbeddedContainers)
		if err != nil {
			return nil, "", fmt.Errorf("failed to merge containers spec: %w", err)
		}
		d.Spec.Template.Spec.Containers = containers
	}

	return d, resources.OperationCreateOrUpdate, ctrl.SetControllerReference(r.router, d, r.Scheme)
}
