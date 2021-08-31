package thanosstorage

import (
	corev1 "k8s.io/api/core/v1"
)

func (s *ThanosStorage) gatewayVolumes() []corev1.Volume {
	var (
		objStorageConfig = s.Instance.Spec.ObjectStorageConfig

		volumes []corev1.Volume
	)

	if objStorageConfig != nil && objStorageConfig.Name != "" {
		volumes = append(volumes, corev1.Volume{
			Name: volumeNamePrefixSecret + objStorageConfig.Name,
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: objStorageConfig.Name,
					Items: []corev1.KeyToPath{{
						Key:  objStorageConfig.Key,
						Path: objStorageConfig.Key,
					}},
				},
			},
		})
	}

	return volumes
}

func (s *ThanosStorage) compactVolumes() []corev1.Volume {
	var (
		objStorageConfig = s.Instance.Spec.ObjectStorageConfig

		volumes []corev1.Volume
	)

	if objStorageConfig != nil && objStorageConfig.Name != "" {
		volumes = append(volumes, corev1.Volume{
			Name: volumeNamePrefixSecret + objStorageConfig.Name,
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: objStorageConfig.Name,
					Items: []corev1.KeyToPath{{
						Key:  objStorageConfig.Key,
						Path: objStorageConfig.Key,
					}},
				},
			},
		})
	}

	return volumes
}
