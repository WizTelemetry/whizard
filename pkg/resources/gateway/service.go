package gateway

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/kubesphere/paodin-monitoring/pkg/resources"
)

func (g *Gateway) service() (runtime.Object, resources.Operation, error) {
	var s = &corev1.Service{ObjectMeta: g.meta(g.name("operated"))}

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
				Port:     9080,
			},
		},
	}
	return s, resources.OperationCreateOrUpdate, nil
}
