package thanosstorage

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (s *ThanosStorage) services() []*corev1.Service {
	var svcs []*corev1.Service
	if s.Instance.Spec.ObjectStorageConfig == nil {
		return svcs
	}
	if s.Instance.Spec.Gateway != nil {
		svcs = append(svcs, s.gatewayService())
	}
	return svcs
}

func (s *ThanosStorage) gatewayService() *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: s.Instance.Namespace,
			Name:      s.getGatewayOperatedServiceName(),
			Labels:    s.gatewayLabels(),
		},
		Spec: corev1.ServiceSpec{
			ClusterIP: corev1.ClusterIPNone,
			Type:      corev1.ServiceTypeClusterIP,
			Selector:  s.gatewayLabels(),
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
