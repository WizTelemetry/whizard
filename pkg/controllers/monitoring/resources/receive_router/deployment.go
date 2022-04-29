package receive_router

import (
	"fmt"
	"path/filepath"

	"github.com/kubesphere/paodin/pkg/controllers/monitoring/resources"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func (r *ReceiveRouter) deployment() (runtime.Object, resources.Operation, error) {
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
				Name:          resources.ThanosGRPCPortName,
				ContainerPort: resources.ThanosGRPCPort,
			},
			{
				Protocol:      corev1.ProtocolTCP,
				Name:          resources.ThanosHTTPPortName,
				ContainerPort: resources.ThanosHTTPPort,
			},
			{
				Protocol:      corev1.ProtocolTCP,
				Name:          resources.ThanosRemoteWritePortName,
				ContainerPort: resources.ThanosRemoteWritePort,
			},
		},
		LivenessProbe:  resources.ThanosDefaultLivenessProbe(),
		ReadinessProbe: resources.ThanosDefaultReadinessProbe(),
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
	container.Args = append(container.Args, `--label=thanos_receive_replica="$(POD_NAME)"`)
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

	d.Spec.Template.Spec.Containers = append(d.Spec.Template.Spec.Containers, container)

	return d, resources.OperationCreateOrUpdate, nil
}
