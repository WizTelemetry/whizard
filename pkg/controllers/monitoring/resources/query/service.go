package query

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/kubesphere/paodin/pkg/controllers/monitoring/resources"
)

func (q *Query) service() (runtime.Object, resources.Operation, error) {
	var s = &corev1.Service{ObjectMeta: q.meta(q.name("operated"))}

	if q.query == nil {
		return s, resources.OperationDelete, nil
	}
	s.Spec = corev1.ServiceSpec{
		ClusterIP: "None",
		Selector:  q.labels(),
		Ports: []corev1.ServicePort{
			{
				Protocol: corev1.ProtocolTCP,
				Name:     "grpc",
				Port:     10901,
			},
			{
				Protocol: corev1.ProtocolTCP,
				Name:     "http",
				Port:     10902,
			},
		},
	}
	return s, resources.OperationCreateOrUpdate, nil
}
