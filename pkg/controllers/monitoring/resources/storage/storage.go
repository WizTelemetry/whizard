package storage

import (
	"fmt"

	monitoringv1alpha1 "github.com/kubesphere/whizard/pkg/api/monitoring/v1alpha1"
	"github.com/kubesphere/whizard/pkg/constants"
	"github.com/kubesphere/whizard/pkg/controllers/monitoring/resources"
	"github.com/kubesphere/whizard/pkg/util"
	corev1 "k8s.io/api/core/v1"
)

type Storage struct {
	storage *monitoringv1alpha1.Storage
	resources.BaseReconciler
}

func New(reconciler resources.BaseReconciler, storage *monitoringv1alpha1.Storage) *Storage {
	return &Storage{
		storage:        storage,
		BaseReconciler: reconciler,
	}
}

func (s *Storage) Reconcile() error {
	return s.ReconcileResources([]resources.Resource{
		s.updateHashAnnotation,
	})
}

func (s *Storage) String() (string, error) {
	body, err := s.parseObjStoreConfig()
	return string(body), err
}

func (s *Storage) VolumesAndVolumeMounts() ([]corev1.Volume, []corev1.VolumeMount) {
	if s.storage.Spec.S3 == nil {
		return nil, nil
	}

	tls := s.storage.Spec.S3.HTTPConfig.TLSConfig
	var volumes []corev1.Volume
	var volumeMounts []corev1.VolumeMount
	volumes, volumeMounts = appendVolumesAndVolumeMounts(volumes, volumeMounts, tls.CA)
	volumes, volumeMounts = appendVolumesAndVolumeMounts(volumes, volumeMounts, tls.Cert)
	volumes, volumeMounts = appendVolumesAndVolumeMounts(volumes, volumeMounts, tls.Key)

	return volumes, volumeMounts
}

func appendVolumesAndVolumeMounts(volumes []corev1.Volume, volumeMounts []corev1.VolumeMount, ref *corev1.SecretKeySelector) ([]corev1.Volume, []corev1.VolumeMount) {
	if ref == nil {
		return volumes, volumeMounts
	}

	mode := corev1.SecretVolumeSourceDefaultMode
	volume := corev1.Volume{
		Name: ref.Name,
		VolumeSource: corev1.VolumeSource{
			Secret: &corev1.SecretVolumeSource{
				SecretName:  ref.Name,
				DefaultMode: &mode,
			},
		},
	}

	replaced := util.ReplaceInSlice(volumes, func(v interface{}) bool {
		return v.(corev1.Volume).Name == volume.Name
	}, volume)

	if !replaced {
		volumes = append(volumes, volume)
	}

	volumeMount := corev1.VolumeMount{
		Name:      ref.Name,
		ReadOnly:  true,
		MountPath: fmt.Sprintf("%s%s/", constants.ConfigPath, ref.Name),
	}

	replaced = util.ReplaceInSlice(volumeMounts, func(v interface{}) bool {
		return v.(corev1.VolumeMount).Name == volumeMount.Name
	}, volumeMount)

	if !replaced {
		volumeMounts = append(volumeMounts, volumeMount)
	}

	return volumes, volumeMounts
}
