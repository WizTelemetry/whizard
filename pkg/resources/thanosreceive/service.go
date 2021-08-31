package thanosreceive

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (r *ThanosReceive) Service() corev1.Service {
	return corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: r.Instance.Namespace,
			Name:      r.getReceiveOperatedServiceName(),
			Labels:    r.labels(),
		},
		Spec: corev1.ServiceSpec{
			ClusterIP: corev1.ClusterIPNone,
			Type:      corev1.ServiceTypeClusterIP,
			Selector:  r.labels(),
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
				{
					Protocol: corev1.ProtocolTCP,
					Name:     "remote-write",
					Port:     remoteWritePort,
				},
			},
		},
	}
}
