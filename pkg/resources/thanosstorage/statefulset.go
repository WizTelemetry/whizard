package thanosstorage

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (s *ThanosStorage) statefulSets() []*appsv1.StatefulSet {
	var stss []*appsv1.StatefulSet
	if s.Instance.Spec.ObjectStorageConfig == nil {
		return stss
	}
	if s.Instance.Spec.Gateway != nil {
		stss = append(stss, s.gatewayStatefulSet())
	}
	if s.Instance.Spec.Compact != nil {
		stss = append(stss, s.compactStatefulSet())
	}
	return stss
}

func (s *ThanosStorage) gatewayStatefulSet() *appsv1.StatefulSet {
	var replicas int32
	if s.Instance.Spec.Gateway.Replicas == nil {
		replicas = defaultReplicas
	} else if *s.Instance.Spec.Gateway.Replicas > 0 {
		replicas = *s.Instance.Spec.Gateway.Replicas
	}

	sts := appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: s.Instance.Namespace,
			Name:      s.getGatewayStatefulSetName(),
			Labels:    s.gatewayLabels(),
		},
		Spec: appsv1.StatefulSetSpec{
			ServiceName: s.getGatewayOperatedServiceName(),
			Replicas:    &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: s.gatewayLabels(),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: s.gatewayLabels(),
				},
				Spec: corev1.PodSpec{
					Volumes:    s.gatewayVolumes(),
					Containers: []corev1.Container{s.gatewayContainer()},
				},
			},
		},
	}

	if gateway := s.Instance.Spec.Gateway; gateway == nil || gateway.DataVolume == nil ||
		(gateway.DataVolume.EmptyDir == nil && gateway.DataVolume.PersistentVolumeClaim == nil) {
		sts.Spec.Template.Spec.Volumes = append(sts.Spec.Template.Spec.Volumes, corev1.Volume{
			Name: s.getGatewayTSDBVolumeName(),
			VolumeSource: corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{},
			},
		})
	} else if gateway.DataVolume.EmptyDir != nil {
		sts.Spec.Template.Spec.Volumes = append(sts.Spec.Template.Spec.Volumes, corev1.Volume{
			Name: s.getGatewayTSDBVolumeName(),
			VolumeSource: corev1.VolumeSource{
				EmptyDir: gateway.DataVolume.EmptyDir,
			},
		})
	} else {
		pvc := gateway.DataVolume.PersistentVolumeClaim
		pvc.Name = s.getGatewayTSDBVolumeName()
		if pvc.Spec.AccessModes == nil {
			pvc.Spec.AccessModes = []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce}
		}
		sts.Spec.VolumeClaimTemplates = append(sts.Spec.VolumeClaimTemplates, *gateway.DataVolume.PersistentVolumeClaim)
	}

	return &sts
}

func (s *ThanosStorage) compactStatefulSet() *appsv1.StatefulSet {
	var replicas int32
	if s.Instance.Spec.Gateway.Replicas == nil {
		replicas = defaultReplicas
	} else if *s.Instance.Spec.Gateway.Replicas > 0 {
		replicas = *s.Instance.Spec.Gateway.Replicas
	}

	sts := appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: s.Instance.Namespace,
			Name:      s.getCompactStatefulSetName(),
			Labels:    s.compactLabels(),
		},
		Spec: appsv1.StatefulSetSpec{
			ServiceName: s.getCompactOperatedServiceName(),
			Replicas:    &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: s.compactLabels(),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: s.compactLabels(),
				},
				Spec: corev1.PodSpec{
					Volumes:    s.compactVolumes(),
					Containers: []corev1.Container{s.compactContainer()},
				},
			},
		},
	}

	if gateway := s.Instance.Spec.Gateway; gateway == nil || gateway.DataVolume == nil ||
		(gateway.DataVolume.EmptyDir == nil && gateway.DataVolume.PersistentVolumeClaim == nil) {
		sts.Spec.Template.Spec.Volumes = append(sts.Spec.Template.Spec.Volumes, corev1.Volume{
			Name: s.getCompactTSDBVolumeName(),
			VolumeSource: corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{},
			},
		})
	} else if gateway.DataVolume.EmptyDir != nil {
		sts.Spec.Template.Spec.Volumes = append(sts.Spec.Template.Spec.Volumes, corev1.Volume{
			Name: s.getCompactTSDBVolumeName(),
			VolumeSource: corev1.VolumeSource{
				EmptyDir: gateway.DataVolume.EmptyDir,
			},
		})
	} else {
		pvc := gateway.DataVolume.PersistentVolumeClaim
		pvc.Name = s.getCompactTSDBVolumeName()
		if pvc.Spec.AccessModes == nil {
			pvc.Spec.AccessModes = []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce}
		}
		sts.Spec.VolumeClaimTemplates = append(sts.Spec.VolumeClaimTemplates, *gateway.DataVolume.PersistentVolumeClaim)
	}

	return &sts
}
