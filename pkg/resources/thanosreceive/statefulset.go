package thanosreceive

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (r *ThanosReceive) statefulSet() appsv1.StatefulSet {
	var replicas int32
	if r.Instance.Spec.Replicas == nil {
		replicas = defaultReplicas
	} else if *r.Instance.Spec.Replicas > 0 {
		replicas = *r.Instance.Spec.Replicas
	}

	s := appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: r.Instance.Namespace,
			Name:      r.getReceiveStatefulSetName(),
			Labels:    r.labels(),
		},
		Spec: appsv1.StatefulSetSpec{
			ServiceName: r.getReceiveOperatedServiceName(),
			Replicas:    &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: r.labels(),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: r.labels(),
				},
				Spec: corev1.PodSpec{
					Volumes:    r.volumes(),
					Containers: r.containers(),
				},
			},
		},
	}

	if ingestor := r.Instance.Spec.Ingestor; ingestor == nil || ingestor.DataVolume == nil ||
		(ingestor.DataVolume.EmptyDir == nil && ingestor.DataVolume.PersistentVolumeClaim == nil) {
		s.Spec.Template.Spec.Volumes = append(s.Spec.Template.Spec.Volumes, corev1.Volume{
			Name: r.getTSDBVolumeName(),
			VolumeSource: corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{},
			},
		})
	} else if ingestor.DataVolume.EmptyDir != nil {
		s.Spec.Template.Spec.Volumes = append(s.Spec.Template.Spec.Volumes, corev1.Volume{
			Name: r.getTSDBVolumeName(),
			VolumeSource: corev1.VolumeSource{
				EmptyDir: ingestor.DataVolume.EmptyDir,
			},
		})
	} else {
		pvc := ingestor.DataVolume.PersistentVolumeClaim
		pvc.Name = r.getTSDBVolumeName()
		if pvc.Spec.AccessModes == nil {
			pvc.Spec.AccessModes = []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce}
		}
		s.Spec.VolumeClaimTemplates = append(s.Spec.VolumeClaimTemplates, *ingestor.DataVolume.PersistentVolumeClaim)
	}

	return s
}
