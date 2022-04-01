package receive

import (
	"fmt"
	"path/filepath"

	"github.com/kubesphere/paodin-monitoring/pkg/resources"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func (r *receiveRouter) deployment() (runtime.Object, resources.Operation, error) {
	var d = &appsv1.Deployment{ObjectMeta: r.meta(r.name())}

	if r.del {
		return d, resources.OperationDelete, nil
	}

	d.Spec = appsv1.DeploymentSpec{
		Replicas: r.Router.Replicas,
		Selector: &metav1.LabelSelector{
			MatchLabels: r.labels(),
		},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: r.labels(),
			},
			Spec: corev1.PodSpec{
				NodeSelector: r.Router.NodeSelector,
				Tolerations:  r.Router.Tolerations,
				Affinity:     r.Router.Affinity,
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
		Image:     r.Router.Image,
		Args:      []string{"receive"},
		Resources: r.Router.Resources,
		Ports: []corev1.ContainerPort{
			{
				Name:          "grpc",
				ContainerPort: 10901,
				Protocol:      corev1.ProtocolTCP,
			},
			{
				Name:          "http",
				ContainerPort: 10902,
				Protocol:      corev1.ProtocolTCP,
			},
		},
		LivenessProbe: &corev1.Probe{
			FailureThreshold: 4,
			PeriodSeconds:    30,
			Handler: corev1.Handler{
				HTTPGet: &corev1.HTTPGetAction{
					Scheme: "HTTP",
					Path:   "/-/healthy",
					Port:   intstr.FromString("http"),
				},
			},
		},
		ReadinessProbe: &corev1.Probe{
			FailureThreshold: 20,
			PeriodSeconds:    5,
			Handler: corev1.Handler{
				HTTPGet: &corev1.HTTPGetAction{
					Scheme: "HTTP",
					Path:   "/-/ready",
					Port:   intstr.FromString("http"),
				},
			},
		},
		VolumeMounts: []corev1.VolumeMount{{
			Name:      hashringsVol.Name,
			MountPath: configDir,
			ReadOnly:  true,
		}},
	}
	if r.Router.LogLevel != "" {
		container.Args = append(container.Args, "--log.level="+r.Router.LogLevel)
	}
	if r.Router.LogFormat != "" {
		container.Args = append(container.Args, "--log.format="+r.Router.LogFormat)
	}
	container.Args = append(container.Args, `--label=receive_replica="$(NAME)"`)
	container.Args = append(container.Args, "--receive.hashrings-file="+filepath.Join(configDir, hashringsFile))
	if r.Router.ReplicationFactor != nil {
		container.Args = append(container.Args, fmt.Sprintf("--receive.replication-factor=%d", *r.Router.ReplicationFactor))
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
