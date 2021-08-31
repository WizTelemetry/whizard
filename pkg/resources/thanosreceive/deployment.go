package thanosreceive

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (r *ThanosReceive) deployment() appsv1.Deployment {
	var replicas int32
	if r.Instance.Spec.Replicas == nil {
		replicas = defaultReplicas
	} else if *r.Instance.Spec.Replicas > 0 {
		replicas = *r.Instance.Spec.Replicas
	}

	d := appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: r.Instance.Namespace,
			Name:      r.getReceiveDeploymentName(),
			Labels:    r.labels(),
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
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

	d.Spec.Template.Spec.Volumes = append(d.Spec.Template.Spec.Volumes, corev1.Volume{
		Name: r.getTSDBVolumeName(),
		VolumeSource: corev1.VolumeSource{
			EmptyDir: &corev1.EmptyDirVolumeSource{},
		},
	})

	return d
}
