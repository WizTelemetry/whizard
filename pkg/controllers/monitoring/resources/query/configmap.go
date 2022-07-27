package query

import (
	"fmt"

	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/discovery/targetgroup"
	"gopkg.in/yaml.v3"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/kubesphere/paodin/pkg/controllers/monitoring/resources"
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
	data, err := envoyConfigFiles(q.Service.Namespace+"/"+q.Service.Name, *stores)
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

	out, err := yaml.Marshal([]targetgroup.Group{{Targets: targets}})
	if err != nil {
		return nil, resources.OperationCreateOrUpdate, err
	}

	cm.Data = map[string]string{
		storesFile: string(out),
	}

	return cm, resources.OperationCreateOrUpdate, nil
}
