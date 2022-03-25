package query

import (
	"fmt"

	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/discovery/targetgroup"
	"gopkg.in/yaml.v3"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/kubesphere/paodin-monitoring/pkg/resources"
	"github.com/kubesphere/paodin-monitoring/pkg/resources/receive"
	"github.com/kubesphere/paodin-monitoring/pkg/resources/storegateway"
)

func (q *Query) proxyConfigMap() (runtime.Object, resources.Operation, error) {

	var cm = &corev1.ConfigMap{ObjectMeta: q.meta(q.name("proxy-config"))}

	if q.query == nil {
		return cm, resources.OperationDelete, nil
	}

	stores, err := q.stores()
	if err != nil {
		return nil, resources.OperationCreateOrUpdate, err
	}
	data, err := envoyConfigFiles(q.Thanos.Namespace+"/"+q.Thanos.Name, *stores)
	if err != nil {
		return nil, resources.OperationCreateOrUpdate, err
	}
	cm.Data = data

	return cm, resources.OperationCreateOrUpdate, nil
}

func (q *Query) storesConfigMap() (runtime.Object, resources.Operation, error) {
	var cm = &corev1.ConfigMap{ObjectMeta: q.meta(q.name("stores-config"))}
	if q.query == nil {
		return cm, resources.OperationDelete, nil
	}

	stores, err := q.stores()
	if err != nil {
		return nil, resources.OperationCreateOrUpdate, err
	}

	var targets []model.LabelSet
	for _, store := range stores.DirectStores {
		targets = append(targets, model.LabelSet{
			model.AddressLabel: model.LabelValue(store.Address),
		})
	}
	for _, store := range stores.ProxyStores {
		targets = append(targets, model.LabelSet{
			model.AddressLabel: model.LabelValue(fmt.Sprintf("%s:%d", store.ListenHost, store.ListenPort)),
		})
	}

	innerStores := receive.New(q.ThanosBaseReconciler).GrpcAddrs()
	innerStores = append(innerStores, storegateway.New(q.ThanosBaseReconciler).GrpcAddrs()...)
	for _, store := range innerStores {
		if store != "" {
			targets = append(targets, model.LabelSet{
				model.AddressLabel: model.LabelValue(store),
			})
		}
	}

	out, err := yaml.Marshal([]targetgroup.Group{{Targets: targets}})
	if err != nil {
		return nil, resources.OperationCreateOrUpdate, err
	}

	cm.Data = map[string]string{
		storesFile: string(out),
	}

	return cm, resources.OperationCreateOrUpdate, nil
}
