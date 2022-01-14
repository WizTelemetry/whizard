package thanosquery

import (
	"path/filepath"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func (q *ThanosQuery) containers() []corev1.Container {
	return []corev1.Container{
		q.queryContainer(),
		q.envoyContainer(),
	}
}

func (q *ThanosQuery) queryContainer() corev1.Container {
	var (
		image           = q.Instance.Spec.Image
		imagePullPolicy = q.Instance.Spec.ImagePullPolicy
	)
	if image == "" {
		image = q.Cfg.ThanosDefaultImage
	}

	var c = corev1.Container{
		Name:            componentName,
		Image:           image,
		ImagePullPolicy: imagePullPolicy,
		Resources:       q.Instance.Spec.Resources,
		Args:            q.args(),
	}

	c.Ports = []corev1.ContainerPort{
		{
			Name:          "grpc",
			ContainerPort: grpcPort,
			Protocol:      corev1.ProtocolTCP,
		},
		{
			Name:          "http",
			ContainerPort: httpPort,
			Protocol:      corev1.ProtocolTCP,
		},
	}

	c.LivenessProbe = &corev1.Probe{
		FailureThreshold: 4,
		PeriodSeconds:    30,
		Handler: corev1.Handler{
			HTTPGet: &corev1.HTTPGetAction{
				Scheme: "HTTP",
				Path:   "/-/healthy",
				Port: intstr.IntOrString{
					Type:   intstr.Int,
					IntVal: httpPort,
				},
			},
		},
	}

	c.ReadinessProbe = &corev1.Probe{
		FailureThreshold: 20,
		PeriodSeconds:    5,
		Handler: corev1.Handler{
			HTTPGet: &corev1.HTTPGetAction{
				Scheme: "HTTP",
				Path:   "/-/ready",
				Port: intstr.IntOrString{
					Type:   intstr.Int,
					IntVal: httpPort,
				},
			},
		},
	}

	c.Env = []corev1.EnvVar{
		{
			Name: "POD_IP",
			ValueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{
					FieldPath: "status.podIP",
				},
			},
		},
	}

	c.VolumeMounts = append(c.VolumeMounts, corev1.VolumeMount{
		Name:      volumeNamePrefixConfigMap + q.getStoreSDConfigMapName(),
		MountPath: thanosConfigDir,
		ReadOnly:  true,
	})

	return c
}

func (q *ThanosQuery) envoyContainer() corev1.Container {
	var (
		image           = q.Cfg.EnvoyDefaultImage
		imagePullPolicy = corev1.PullIfNotPresent
		resources       corev1.ResourceRequirements
	)
	if envoySpec := q.Instance.Spec.Envoy; envoySpec != nil {
		if envoySpec.Image != "" {
			image = envoySpec.Image
		}
		if envoySpec.ImagePullPolicy != "" {
			imagePullPolicy = envoySpec.ImagePullPolicy
		}
		resources = envoySpec.Resources
	}

	var c = corev1.Container{
		Name:            "envoy-sidecar",
		Image:           image,
		ImagePullPolicy: imagePullPolicy,
		Resources:       resources,
		Args: []string{
			"-c",
			filepath.Join(envoyConfigDir, envoyConfigFileName),
			"-l",
			"debug",
		},
	}

	c.VolumeMounts = append(c.VolumeMounts, corev1.VolumeMount{
		Name:      volumeNamePrefixConfigMap + q.getEnvoyConfigMapName(),
		MountPath: envoyConfigDir,
		ReadOnly:  true,
	})

	var secrets = make(map[string]struct{})
	for _, store := range q.Instance.Spec.Stores {
		if !storeRequireProxy(store) {
			continue
		}
		if _, ok := secrets[store.SecretName]; ok {
			continue
		}
		c.VolumeMounts = append(c.VolumeMounts, corev1.VolumeMount{
			Name:      volumeNamePrefixSecret + store.SecretName,
			MountPath: filepath.Join(envoySecretsDir, store.SecretName),
			ReadOnly:  true,
		})
		secrets[store.SecretName] = struct{}{}
	}

	return c
}
