package thanosreceive

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"path/filepath"
)

func (r *ThanosReceive) containers() []corev1.Container {
	return []corev1.Container{
		r.receiveContainer(),
	}
}

func (r *ThanosReceive) receiveContainer() corev1.Container {
	var (
		image           = r.Instance.Spec.Image
		imagePullPolicy = r.Instance.Spec.ImagePullPolicy
	)
	if image == "" {
		image = r.Cfg.ThanosDefaultImage
	}

	c := corev1.Container{
		Name:            componentName,
		Image:           image,
		ImagePullPolicy: imagePullPolicy,
		Resources:       r.Instance.Spec.Resources,
		Args:            r.args(),
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
		{
			Name:          "remote-write",
			ContainerPort: remoteWritePort,
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
		Name:      r.getTSDBVolumeName(),
		MountPath: storageDir,
	})

	if m := r.GetMode(); m == RouterOnly || m == RouterIngestor {
		c.VolumeMounts = append(c.VolumeMounts, corev1.VolumeMount{
			Name:      volumeNamePrefixConfigMap + r.getHashringsConfigMapName(),
			MountPath: thanosConfigDir,
			ReadOnly:  true,
		})
	}

	if r.Instance.Spec.Ingestor != nil {
		if configSecret := r.Instance.Spec.Ingestor.ObjectStorageConfig; configSecret != nil && configSecret.Name != "" {
			c.VolumeMounts = append(c.VolumeMounts, corev1.VolumeMount{
				Name:      volumeNamePrefixSecret + configSecret.Name,
				MountPath: filepath.Join(mountDirSecrets, configSecret.Name),
				ReadOnly:  true,
			})
		}
	}

	return c
}
