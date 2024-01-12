package storage

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

func (s *Storage) service() (runtime.Object, resources.Operation, error) {
	var svc = &corev1.Service{ObjectMeta: s.meta(s.name(constants.ServiceNameSuffix))}

	if s.storage.Spec.BlockManager == nil ||
		s.storage.Spec.BlockManager.Enable == nil ||
		*s.storage.Spec.BlockManager.Enable == false {
		return svc, resources.OperationDelete, nil
	}

	if err := s.Client.Get(s.Context, client.ObjectKeyFromObject(svc), svc); err != nil {
		if !util.IsNotFound(err) {
			return nil, "", err
		}
	}

	svc.Spec.Selector = s.labels()

	port := corev1.ServicePort{
		Protocol:   corev1.ProtocolTCP,
		Name:       constants.HTTPPortName,
		Port:       constants.HTTPPort,
		TargetPort: intstr.FromInt(constants.HTTPPort),
	}

	if s.storage.Spec.BlockManager.NodePort != 0 {
		svc.Spec.Type = corev1.ServiceTypeNodePort
		port.NodePort = s.storage.Spec.BlockManager.NodePort
	} else {
		svc.Spec.Type = corev1.ServiceTypeClusterIP
	}

	replaced := util.ReplaceInSlice(svc.Spec.Ports, func(v interface{}) bool {
		port := v.(corev1.ServicePort)
		return port.Name == port.Name
	}, port)

	if !replaced {
		svc.Spec.Ports = append(svc.Spec.Ports, port)
	}

	return svc, resources.OperationCreateOrUpdate, ctrl.SetControllerReference(s.storage, svc, s.Scheme)
}
