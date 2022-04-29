package gateway

import (
	"fmt"
	"path/filepath"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/kubesphere/paodin/pkg/controllers/monitoring/resources"
	"github.com/kubesphere/paodin/pkg/controllers/monitoring/resources/query"
	"github.com/kubesphere/paodin/pkg/controllers/monitoring/resources/receive_router"
)

func (g *Gateway) deployment() (runtime.Object, resources.Operation, error) {
	var d = &appsv1.Deployment{ObjectMeta: g.meta(g.name())}

	if g.gateway == nil {
		return d, resources.OperationDelete, nil
	}

	d.Spec = appsv1.DeploymentSpec{
		Replicas: g.gateway.Replicas,
		Selector: &metav1.LabelSelector{
			MatchLabels: g.labels(),
		},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: g.labels(),
			},
			Spec: corev1.PodSpec{
				NodeSelector: g.gateway.NodeSelector,
				Tolerations:  g.gateway.Tolerations,
				Affinity:     g.gateway.Affinity,
			},
		},
	}

	var container = corev1.Container{
		Name:      "gateway",
		Image:     g.gateway.Image,
		Args:      []string{},
		Resources: g.gateway.Resources,
		Ports: []corev1.ContainerPort{
			{
				Name:          "http",
				ContainerPort: 9090,
				Protocol:      corev1.ProtocolTCP,
			},
		},
		// LivenessProbe: &corev1.Probe{
		// 	FailureThreshold: 4,
		// 	PeriodSeconds:    30,
		// 	ProbeHandler: corev1.ProbeHandler{
		// 		HTTPGet: &corev1.HTTPGetAction{
		// 			Scheme: "HTTP",
		// 			Path:   "/-/healthy",
		// 			Port:   intstr.FromString("http"),
		// 		},
		// 	},
		// },
		// ReadinessProbe: &corev1.Probe{
		// 	FailureThreshold: 20,
		// 	PeriodSeconds:    5,
		// 	ProbeHandler: corev1.ProbeHandler{
		// 		HTTPGet: &corev1.HTTPGetAction{
		// 			Scheme: "HTTP",
		// 			Path:   "/-/ready",
		// 			Port:   intstr.FromString("http"),
		// 		},
		// 	},
		// },
	}

	if g.gateway.ServerCertificate != "" {
		serverCertVol := corev1.Volume{
			Name: "server-certificate",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: g.gateway.ServerCertificate,
				},
			},
		}
		d.Spec.Template.Spec.Volumes = append(d.Spec.Template.Spec.Volumes, serverCertVol)
		container.VolumeMounts = append(container.VolumeMounts, corev1.VolumeMount{
			Name:      serverCertVol.Name,
			MountPath: filepath.Join(secretsDir, g.gateway.ServerCertificate),
		})
		container.Args = append(container.Args, "--server-tls-key="+filepath.Join(secretsDir, g.gateway.ServerCertificate, "tls.key"))
		container.Args = append(container.Args, "--server-tls-crt="+filepath.Join(secretsDir, g.gateway.ServerCertificate, "tls.crt"))
	}

	if g.gateway.ClientCACertificate != "" {
		clientCaCertVol := corev1.Volume{
			Name: "client-ca-certificate",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: g.gateway.ClientCACertificate,
				},
			},
		}
		d.Spec.Template.Spec.Volumes = append(d.Spec.Template.Spec.Volumes, clientCaCertVol)
		container.VolumeMounts = append(container.VolumeMounts, corev1.VolumeMount{
			Name:      clientCaCertVol.Name,
			MountPath: filepath.Join(secretsDir, g.gateway.ClientCACertificate),
		})
		container.Args = append(container.Args, "--server-tls-client-ca="+filepath.Join(secretsDir, g.gateway.ClientCACertificate, "tls.crt"))
	}

	if g.gateway.LogLevel != "" {
		container.Args = append(container.Args, "--log.level="+g.gateway.LogLevel)
	}
	if g.gateway.LogFormat != "" {
		container.Args = append(container.Args, "--log.format="+g.gateway.LogFormat)
	}

	if g.Service.Spec.TenantHeader != "" {
		container.Args = append(container.Args, "--tenant.header="+g.Service.Spec.TenantHeader)
	}
	if g.Service.Spec.TenantLabelName != "" {
		container.Args = append(container.Args, "--tenant.label-name="+g.Service.Spec.TenantLabelName)
	}

	if thanos := g.Service.Spec.Thanos; thanos != nil {
		if thanos.Query != nil {
			q := query.New(g.ServiceBaseReconciler)
			container.Args = append(container.Args, fmt.Sprintf("--query.address=%s", q.HttpAddr()))
		}
		if thanos.ReceiveRouter != nil {
			r := receive_router.New(g.ServiceBaseReconciler)
			container.Args = append(container.Args, fmt.Sprintf("--remote-write.address=%s", r.RemoteWriteAddr()))
		}
	}

	d.Spec.Template.Spec.Containers = append(d.Spec.Template.Spec.Containers, container)

	return d, resources.OperationCreateOrUpdate, nil
}
