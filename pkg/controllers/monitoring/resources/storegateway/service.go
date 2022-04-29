package storegateway

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/kubesphere/paodin/pkg/controllers/monitoring/resources"
)

func (r *StoreGateway) service() (runtime.Object, resources.Operation, error) {
	var s = &corev1.Service{ObjectMeta: r.meta(r.name(resources.ServiceNameSuffixOperated))}

	if r.store == nil || r.Store.Spec.ObjectStorageConfig == nil {
		return s, resources.OperationDelete, nil
	}

	s.Spec = corev1.ServiceSpec{
		ClusterIP: corev1.ClusterIPNone,
		Selector:  r.labels(),
		Ports: []corev1.ServicePort{
			{
				Protocol: corev1.ProtocolTCP,
				Name:     resources.ThanosGRPCPortName,
				Port:     resources.ThanosGRPCPort,
			},
			{
				Protocol: corev1.ProtocolTCP,
				Name:     resources.ThanosHTTPPortName,
				Port:     resources.ThanosHTTPPort,
			},
		},
	}
	return s, resources.OperationCreateOrUpdate, nil
}
