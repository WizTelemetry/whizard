package gateway

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kubesphere/whizard/pkg/constants"
	"github.com/kubesphere/whizard/pkg/controllers/monitoring/resources"
	"github.com/kubesphere/whizard/pkg/util"
)

func (g *Gateway) service() (runtime.Object, resources.Operation, error) {
	var s = &corev1.Service{ObjectMeta: g.meta(g.name(constants.ServiceNameSuffix))}
	if err := g.Client.Get(g.Context, client.ObjectKeyFromObject(s), s); err != nil {
		if !util.IsNotFound(err) {
			return nil, "", err
		}
	}

	s.Spec.Selector = g.labels()

	port := corev1.ServicePort{
		Protocol:   corev1.ProtocolTCP,
		Name:       constants.HTTPPortName,
		Port:       9090,
		TargetPort: intstr.FromInt(9090),
		NodePort:   g.gateway.Spec.NodePort,
	}
	if g.gateway.Spec.NodePort != 0 {
		s.Spec.Type = corev1.ServiceTypeNodePort
		port.NodePort = g.gateway.Spec.NodePort
	} else {
		s.Spec.Type = corev1.ServiceTypeClusterIP
	}

	replaced := util.ReplaceInSlice(s.Spec.Ports, func(v interface{}) bool {
		port := v.(corev1.ServicePort)
		return port.Name == port.Name
	}, port)

	if !replaced {
		s.Spec.Ports = append(s.Spec.Ports, port)
	}

	return s, resources.OperationCreateOrUpdate, ctrl.SetControllerReference(g.gateway, s, g.Scheme)
}
