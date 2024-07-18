package store

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/WhizardTelemetry/whizard/pkg/api/monitoring/v1alpha1"
	"github.com/WhizardTelemetry/whizard/pkg/constants"
	"github.com/WhizardTelemetry/whizard/pkg/controllers/resources"
	"github.com/WhizardTelemetry/whizard/pkg/util"
)

func (r *Store) services() (retResources []resources.Resource) {
	timeRanges := r.store.Spec.TimeRanges
	if len(timeRanges) == 0 {
		timeRanges = append(timeRanges, v1alpha1.TimeRange{
			MinTime: r.store.Spec.MinTime,
			MaxTime: r.store.Spec.MaxTime,
		})
	}
	// for expected services
	var expectNames = make(map[string]struct{}, len(timeRanges))
	for i := range timeRanges {
		partitionSn := i
		partitionName := r.partitionName(i, constants.ServiceNameSuffix)
		expectNames[partitionName] = struct{}{}
		retResources = append(retResources, func() (runtime.Object, resources.Operation, error) {
			return r.service(partitionName, partitionSn)
		})
	}

	var svcList corev1.ServiceList
	ls := r.BaseLabels()
	ls[constants.LabelNameAppName] = constants.AppNameStore
	ls[constants.LabelNameAppManagedBy] = r.store.Name
	err := r.Client.List(r.Context, &svcList, client.InNamespace(r.store.Namespace), &client.ListOptions{
		LabelSelector: labels.SelectorFromSet(ls),
	})
	if err != nil {
		return errResourcesFunc(err)
	}
	// check services to be deleted.
	for i := range svcList.Items {
		svc := svcList.Items[i]
		if _, ok := expectNames[svc.Name]; !ok {
			retResources = append(retResources, func() (runtime.Object, resources.Operation, error) {
				return &svc, resources.OperationDelete, nil
			})
		}
	}
	return
}

func (r *Store) service(name string, partitionSn int) (runtime.Object, resources.Operation, error) {
	var s = &corev1.Service{ObjectMeta: r.meta(name, partitionSn)}

	if err := r.Client.Get(r.Context, client.ObjectKeyFromObject(s), s); err != nil {
		if !util.IsNotFound(err) {
			return nil, "", err
		}
	}

	s.Spec.Type = corev1.ServiceTypeClusterIP
	s.Spec.Selector = r.labels(partitionSn)

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
