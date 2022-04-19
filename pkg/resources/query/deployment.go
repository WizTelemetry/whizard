package query

import (
	"path/filepath"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/kubesphere/paodin-monitoring/pkg/resources"
)

func (q *Query) deployment() (runtime.Object, resources.Operation, error) {
	var d = &appsv1.Deployment{ObjectMeta: q.meta(q.name())}

	if q.query == nil {
		return d, resources.OperationDelete, nil
	}

	d.Spec = appsv1.DeploymentSpec{
		Replicas: q.query.Replicas,
		Selector: &metav1.LabelSelector{
			MatchLabels: q.labels(),
		},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: q.labels(),
			},
			Spec: corev1.PodSpec{
				NodeSelector: q.query.NodeSelector,
				Tolerations:  q.query.Tolerations,
				Affinity:     q.query.Affinity,
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
		Image:     q.query.Image,
		Args:      []string{"query"},
		Resources: q.query.Resources,
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
			Name:      storesConfigVol.Name,
			MountPath: configDir,
			ReadOnly:  true,
		}},
	}

	if q.query.LogLevel != "" {
		queryContainer.Args = append(queryContainer.Args, "--log.level="+q.query.LogLevel)
	}
	if q.query.LogFormat != "" {
		queryContainer.Args = append(queryContainer.Args, "--log.format="+q.query.LogFormat)
	}
	queryContainer.Args = append(queryContainer.Args, "--store.sd-files="+filepath.Join(configDir, storesFile))
	for _, ln := range q.query.ReplicaLabelNames {
		queryContainer.Args = append(queryContainer.Args, "--query.replica-label="+ln)
	}

	var envoyContainer = corev1.Container{
		Name:  "proxy",
		Image: q.query.Envoy.Image,
		Args: []string{
			"-c",
			filepath.Join(envoyConfigDir, envoyConfigFile),
			// "-l",
			// "debug",
		},
		Resources: q.query.Envoy.Resources,
		VolumeMounts: []corev1.VolumeMount{{
			Name:      proxyConfigVol.Name,
			MountPath: envoyConfigDir,
			ReadOnly:  true,
		}},
	}

	for _, store := range q.query.Stores {
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

	return d, resources.OperationCreateOrUpdate, nil
}
