package thanosquery

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (q *ThanosQuery) service() corev1.Service {
	return corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: q.Instance.Namespace,
			Name:      q.getServiceName(),
			Labels:    q.labels(),
		},
		Spec: corev1.ServiceSpec{
			ClusterIP: "None",
			Selector:  q.labels(),
			Ports: []corev1.ServicePort{
				{
					Protocol: corev1.ProtocolTCP,
					Name:     "grpc",
					Port:     grpcPort,
				},
				{
					Protocol: corev1.ProtocolTCP,
					Name:     "http",
					Port:     httpPort,
				},
			},
		},
	}
}
