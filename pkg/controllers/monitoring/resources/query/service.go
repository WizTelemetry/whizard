package query

import (
	"github.com/kubesphere/whizard/pkg/constants"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/kubesphere/whizard/pkg/controllers/monitoring/resources"
)

func (q *Query) service() (runtime.Object, resources.Operation, error) {
	var s = &corev1.Service{ObjectMeta: q.meta(q.name(constants.ServiceNameSuffix))}

	if q.query == nil {
		return s, resources.OperationDelete, nil
	}
	s.Spec = corev1.ServiceSpec{
		ClusterIP: "None",
		Selector:  q.labels(),
		Ports: []corev1.ServicePort{
			{
				Protocol: corev1.ProtocolTCP,
				Name:     constants.GRPCPortName,
				Port:     constants.GRPCPort,
			},
			{
				Protocol: corev1.ProtocolTCP,
				Name:     constants.HTTPPortName,
				Port:     constants.HTTPPort,
			},
		},
	}
	return s, resources.OperationCreateOrUpdate, nil
}
