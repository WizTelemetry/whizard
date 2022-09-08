package store

import (
	"github.com/kubesphere/whizard/pkg/constants"
	"github.com/kubesphere/whizard/pkg/controllers/monitoring/resources"
	"github.com/kubesphere/whizard/pkg/util"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (r *Store) service() (runtime.Object, resources.Operation, error) {
	var s = &corev1.Service{ObjectMeta: r.meta(r.name(constants.ServiceNameSuffix))}

	if err := r.Client.Get(r.Context, client.ObjectKeyFromObject(s), s); err != nil {
		if !util.IsNotFound(err) {
			return nil, "", err
		}
	}

	s.Spec.ClusterIP = corev1.ClusterIPNone
	s.Spec.Selector = r.labels()

	ports := []corev1.ServicePort{
		{
			Protocol:   corev1.ProtocolTCP,
			Name:       constants.GRPCPortName,
			Port:       constants.GRPCPort,
			TargetPort: intstr.FromInt(constants.GRPCPort),
		},
		{
			Protocol:   corev1.ProtocolTCP,
			Name:       constants.HTTPPortName,
			Port:       constants.HTTPPort,
			TargetPort: intstr.FromInt(constants.HTTPPort),
		},
	}

	for i := 0; i < len(ports); i++ {
		replaced := util.ReplaceInSlice(s.Spec.Ports, func(v interface{}) bool {
			port := v.(corev1.ServicePort)
			return port.Name == ports[i].Name
		}, ports[i])

		if !replaced {
			s.Spec.Ports = append(s.Spec.Ports, ports[i])
		}
	}

	return s, resources.OperationCreateOrUpdate, ctrl.SetControllerReference(r.store, s, r.Scheme)
}
