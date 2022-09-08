package ingester

import (
	"github.com/kubesphere/whizard/pkg/constants"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/kubesphere/whizard/pkg/controllers/monitoring/resources"
)

func (r *Ingester) service() (runtime.Object, resources.Operation, error) {
	var s = &corev1.Service{ObjectMeta: r.meta(r.name(constants.ServiceNameSuffix))}

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
			{
				Protocol: corev1.ProtocolTCP,
				Name:     constants.RemoteWritePortName,
				Port:     constants.RemoteWritePort,
			},
		},
	}
	return s, resources.OperationCreateOrUpdate, ctrl.SetControllerReference(r.ingester, s, r.Scheme)
}
