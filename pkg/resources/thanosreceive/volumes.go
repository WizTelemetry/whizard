package thanosreceive

import (
	corev1 "k8s.io/api/core/v1"
)

func (r *ThanosReceive) volumes() []corev1.Volume {
	var volumes []corev1.Volume

	if m := r.GetMode(); m == RouterOnly || m == RouterIngestor {
		volumes = append(volumes, corev1.Volume{
			Name: volumeNamePrefixConfigMap + r.getHashringsConfigMapName(),
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: r.getHashringsConfigMapName(),
					},
				},
			},
		})
	}

	if r.Instance.Spec.Ingestor != nil {
		if configSecret := r.Instance.Spec.Ingestor.ObjectStorageConfig; configSecret != nil && configSecret.Name != "" {
			volumes = append(volumes, corev1.Volume{
				Name: volumeNamePrefixSecret + configSecret.Name,
				VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						SecretName: configSecret.Name,
						Items: []corev1.KeyToPath{{
							Key:  configSecret.Key,
							Path: configSecret.Key,
						}},
					},
				},
			})
		}
	}

	return volumes
}
