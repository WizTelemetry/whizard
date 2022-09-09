package gateway

import (
	"github.com/kubesphere/whizard/pkg/constants"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/kubesphere/whizard/pkg/controllers/monitoring/resources"
)

func (g *Gateway) service() (runtime.Object, resources.Operation, error) {
	var s = &corev1.Service{ObjectMeta: g.meta(g.name(constants.ServiceNameSuffix))}

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
	return s, resources.OperationCreateOrUpdate, ctrl.SetControllerReference(g.gateway, s, g.Scheme)
}
