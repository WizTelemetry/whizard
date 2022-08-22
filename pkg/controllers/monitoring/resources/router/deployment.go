package router

import (
	"fmt"
	"path/filepath"

	"github.com/kubesphere/whizard/pkg/constants"
	"github.com/kubesphere/whizard/pkg/controllers/monitoring/resources"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func (r *Router) deployment() (runtime.Object, resources.Operation, error) {
	var d = &appsv1.Deployment{ObjectMeta: r.meta(r.name())}

	if r.router == nil {
		return d, resources.OperationDelete, nil
	}

	d.Spec = appsv1.DeploymentSpec{
		Replicas: r.router.Replicas,
		Selector: &metav1.LabelSelector{
			MatchLabels: r.labels(),
		},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: r.labels(),
			},
			Spec: corev1.PodSpec{
				NodeSelector: r.router.NodeSelector,
				Tolerations:  r.router.Tolerations,
				Affinity:     r.router.Affinity,
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

	var container = corev1.Container{
		Name:      "receive",
		Image:     r.router.Image,
		Args:      []string{"receive"},
		Resources: r.router.Resources,
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
		LivenessProbe:  resources.DefaultLivenessProbe(),
		ReadinessProbe: resources.DefaultReadinessProbe(),
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
	if r.router.LogLevel != "" {
		container.Args = append(container.Args, "--log.level="+r.router.LogLevel)
	}
	if r.router.LogFormat != "" {
		container.Args = append(container.Args, "--log.format="+r.router.LogFormat)
	}
	container.Args = append(container.Args, fmt.Sprintf("--label=%s=\"$(POD_NAME)\"", constants.ReceiveReplicaLabelName))
	container.Args = append(container.Args, "--receive.hashrings-file="+filepath.Join(configDir, hashringsFile))
	if r.router.ReplicationFactor != nil {
		container.Args = append(container.Args, fmt.Sprintf("--receive.replication-factor=%d", *r.router.ReplicationFactor))
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

	for name, value := range r.router.Flags {
		if name == "receive.local-endpoint" {
			// ignoring these flags to make receiver run with router mode
			continue
		}
		container.Args = append(container.Args, fmt.Sprintf("--%s=%s", name, value))
	}

	d.Spec.Template.Spec.Containers = append(d.Spec.Template.Spec.Containers, container)

	return d, resources.OperationCreateOrUpdate, nil
}
