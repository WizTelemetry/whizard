package ingester

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/WhizardTelemetry/whizard/pkg/constants"
	"github.com/WhizardTelemetry/whizard/pkg/controllers/resources"
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
			{
				Protocol: corev1.ProtocolTCP,
				Name:     constants.CapNProtoPortName,
				Port:     constants.CapNProtoPort,
			},
		},
	}
	return s, resources.OperationCreateOrUpdate, ctrl.SetControllerReference(r.ingester, s, r.Scheme)
}
