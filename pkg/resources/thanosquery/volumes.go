package thanosquery

import (
	corev1 "k8s.io/api/core/v1"
)

func (q *ThanosQuery) volumes() []corev1.Volume {
	var volumes []corev1.Volume
	volumes = append(volumes, corev1.Volume{
		Name: volumeNamePrefixConfigMap + q.getStoreSDConfigMapName(),
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: q.getStoreSDConfigMapName(),
				},
			},
		},
	})

	volumes = append(volumes, corev1.Volume{
		Name: volumeNamePrefixConfigMap + q.getEnvoyConfigMapName(),
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: q.getEnvoyConfigMapName(),
				},
			},
		},
	})

	var secrets = make(map[string]struct{})
	for _, store := range q.Instance.Spec.Stores {
		if store.SecretName == "" {
			continue
		}
		if _, ok := secrets[store.SecretName]; ok {
			continue
		}
		volume := corev1.Volume{
			Name: volumeNamePrefixSecret + store.SecretName,
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: store.SecretName,
				},
			},
		}
		volumes = append(volumes, volume)
		secrets[store.SecretName] = struct{}{}
	}

	return volumes
}
