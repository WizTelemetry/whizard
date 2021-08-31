package thanosquery

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (q *ThanosQuery) deployment() appsv1.Deployment {
	var replicas int32
	if q.Instance.Spec.Replicas == nil {
		replicas = defaultReplicas
	} else if *q.Instance.Spec.Replicas > 0 {
		replicas = *q.Instance.Spec.Replicas
	}

	return appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: q.Instance.Namespace,
			Name:      q.getDeploymentName(),
			Labels:    q.labels(),
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: q.labels(),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: q.labels(),
				},
				Spec: corev1.PodSpec{
					Volumes:    q.volumes(),
					Containers: q.containers(),
					NodeSelector: q.Instance.Spec.NodeSelector,
					Tolerations: q.Instance.Spec.Tolerations,
					Affinity: q.Instance.Spec.Affinity,
				},
			},
		},
	}
}
