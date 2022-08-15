package compactor

import (
	"github.com/kubesphere/paodin/pkg/controllers/monitoring/resources"
	"github.com/kubesphere/paodin/pkg/util"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (r *Compactor) service() (runtime.Object, resources.Operation, error) {
	var s = &corev1.Service{ObjectMeta: r.meta(util.Join("-", r.compactor.Name, resources.ServiceNameSuffixOperated))}
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
			Name:       resources.ThanosHTTPPortName,
			Port:       resources.ThanosHTTPPort,
			TargetPort: intstr.FromInt(resources.ThanosHTTPPort),
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

	return s, resources.OperationCreateOrUpdate, ctrl.SetControllerReference(r.compactor, s, r.Scheme)
}
