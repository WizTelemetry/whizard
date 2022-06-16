package gateway

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/kubesphere/paodin/pkg/controllers/monitoring/resources"
)

func (g *Gateway) service() (runtime.Object, resources.Operation, error) {
	var s = &corev1.Service{ObjectMeta: g.meta(g.name(resources.ServiceNameSuffixOperated))}

	if g.gateway == nil {
		return s, resources.OperationDelete, nil
	}
	s.Spec = corev1.ServiceSpec{
		ClusterIP: "None",
		Selector:  g.labels(),
		Ports: []corev1.ServicePort{
			{
				Protocol: corev1.ProtocolTCP,
				Name:     "http",
				Port:     9090,
			},
		},
	}
	return s, resources.OperationCreateOrUpdate, nil
}