package thanosstorage

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"path/filepath"
)

func (s *ThanosStorage) gatewayContainer() corev1.Container {
	var (
		image           = s.Instance.Spec.Image
		imagePullPolicy = s.Instance.Spec.ImagePullPolicy
	)
	if image == "" {
		image = s.Cfg.ThanosDefaultImage
	}

	c := corev1.Container{
		Name:            componentName,
		Image:           image,
		ImagePullPolicy: imagePullPolicy,
		Resources:       s.Instance.Spec.Gateway.Resources,
		Args:            s.gatewayArgs(),
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
		{
			Name: "NAME",
			ValueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{
					FieldPath: "metadata.name",
				},
			},
		},
	}

	c.VolumeMounts = append(c.VolumeMounts, corev1.VolumeMount{
		Name:      s.getGatewayTSDBVolumeName(),
		MountPath: storageDir,
	})

	if configSecret := s.Instance.Spec.ObjectStorageConfig; configSecret != nil && configSecret.Name != "" {
		c.VolumeMounts = append(c.VolumeMounts, corev1.VolumeMount{
			Name:      volumeNamePrefixSecret + configSecret.Name,
			MountPath: filepath.Join(mountDirSecrets, configSecret.Name),
			ReadOnly:  true,
		})
	}

	return c
}

func (s *ThanosStorage) compactContainer() corev1.Container {
	var (
		image           = s.Instance.Spec.Image
		imagePullPolicy = s.Instance.Spec.ImagePullPolicy
	)
	if image == "" {
		image = s.Cfg.ThanosDefaultImage
	}

	c := corev1.Container{
		Name:            componentName,
		Image:           image,
		ImagePullPolicy: imagePullPolicy,
		Resources:       s.Instance.Spec.Compact.Resources,
		Args:            s.compactArgs(),
	}

	c.Ports = []corev1.ContainerPort{
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
		{
			Name: "NAME",
			ValueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{
					FieldPath: "metadata.name",
				},
			},
		},
	}

	c.VolumeMounts = append(c.VolumeMounts, corev1.VolumeMount{
		Name:      s.getCompactTSDBVolumeName(),
		MountPath: storageDir,
	})

	if configSecret := s.Instance.Spec.ObjectStorageConfig; configSecret != nil && configSecret.Name != "" {
		c.VolumeMounts = append(c.VolumeMounts, corev1.VolumeMount{
			Name:      volumeNamePrefixSecret + configSecret.Name,
			MountPath: filepath.Join(mountDirSecrets, configSecret.Name),
			ReadOnly:  true,
		})
	}

	return c
}
