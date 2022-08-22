package ruler

import (
	"github.com/kubesphere/whizard/pkg/constants"
	"github.com/kubesphere/whizard/pkg/controllers/monitoring/resources"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func (r *Ruler) service() (runtime.Object, resources.Operation, error) {
	var s = &corev1.Service{ObjectMeta: r.meta(r.name(constants.ServiceNameSuffix))}

	if r.ruler == nil {
		return s, resources.OperationDelete, nil
	}

	s.Spec = corev1.ServiceSpec{
		ClusterIP: corev1.ClusterIPNone,
		Selector:  r.labels(),
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
