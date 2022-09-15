package ruler

import (
	"strconv"
	"strings"

	"github.com/kubesphere/whizard/pkg/constants"
	"github.com/kubesphere/whizard/pkg/controllers/monitoring/resources"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (r *Ruler) services() (retResources []resources.Resource) {
	// for target services
	var targetNames = make(map[string]struct{}, *r.ruler.Spec.Shards)
	for i := 0; i < int(*r.ruler.Spec.Shards); i++ {
		shardSn := i
		targetNames[r.name(strconv.Itoa(shardSn), constants.ServiceNameSuffix)] = struct{}{}
		retResources = append(retResources, func() (runtime.Object, resources.Operation, error) {
			return r.service(shardSn)
		})
	}

	var svcList corev1.ServiceList
	err := r.Client.List(r.Context, &svcList, client.InNamespace(r.ruler.Namespace))
	if err != nil {
		return errResourcesFunc(err)
	}
	// check services to be deleted.
	// the services owned by the ruler have a same name prefix and a shard sequence number suffix
	var namePrefix = r.name() + "-"
	var nameSuffix = "-" + constants.ServiceNameSuffix
	for i := range svcList.Items {
		svc := svcList.Items[i]
		if !strings.HasPrefix(svc.Name, namePrefix) || !strings.HasSuffix(svc.Name, nameSuffix) {
			continue
		}
		sn := strings.TrimSuffix(strings.TrimPrefix(svc.Name, namePrefix), nameSuffix)
		if sequenceNumberRegexp.MatchString(sn) {
			if _, ok := targetNames[svc.Name]; !ok {
				retResources = append(retResources, func() (runtime.Object, resources.Operation, error) {
					return &svc, resources.OperationDelete, nil
				})
			}
		}
	}
	return
}

func (r *Ruler) service(shardSn int) (runtime.Object, resources.Operation, error) {
	var s = &corev1.Service{ObjectMeta: r.meta(r.name(strconv.Itoa(shardSn), constants.ServiceNameSuffix))}

	if r.ruler == nil {
		return s, resources.OperationDelete, nil
	}

	ls := r.labels()
	ls[constants.LabelNameRulerShardSn] = strconv.Itoa(shardSn)

	s.Spec = corev1.ServiceSpec{
		Type:     corev1.ServiceTypeClusterIP,
		Selector: ls,
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
		},
	}
	return s, resources.OperationCreateOrUpdate, ctrl.SetControllerReference(r.ruler, s, r.Scheme)
}
