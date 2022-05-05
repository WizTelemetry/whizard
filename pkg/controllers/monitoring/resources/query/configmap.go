package query

import (
	"fmt"

	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/discovery/targetgroup"
	"gopkg.in/yaml.v3"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	monitoringv1alpha1 "github.com/kubesphere/paodin/pkg/api/monitoring/v1alpha1"
	"github.com/kubesphere/paodin/pkg/controllers/monitoring/resources"
	"github.com/kubesphere/paodin/pkg/controllers/monitoring/resources/receive_ingestor"
	"github.com/kubesphere/paodin/pkg/controllers/monitoring/resources/storegateway"
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

	var innerStores []string

	var ingestorList monitoringv1alpha1.ThanosReceiveIngestorList
	if err := q.Client.List(q.Context, &ingestorList,
		client.MatchingLabels(monitoringv1alpha1.ManagedLabelByService(q.Service))); err != nil {

		q.Log.WithValues("thanosreceiveingestorlist", "").Error(err, "")
		return nil, resources.OperationCreateOrUpdate, err
	}
	for _, item := range ingestorList.Items {
		ingestor := receive_ingestor.New(q.BaseReconciler, &item)
		innerStores = append(innerStores, ingestor.GrpcAddrs()...)
	}

	var storeList monitoringv1alpha1.StoreList
	if err := q.Client.List(q.Context, &storeList,
		client.MatchingLabels(monitoringv1alpha1.ManagedLabelByService(q.Service))); err != nil {

		q.Log.WithValues("storelist", "").Error(err, "")
		return nil, resources.OperationCreateOrUpdate, err
	}
	for _, item := range storeList.Items {
		storeGateway := storegateway.New(resources.StoreBaseReconciler{
			BaseReconciler: q.BaseReconciler,
			Store:          &item})
		innerStores = append(innerStores, storeGateway.GrpcAddrs()...)
	}

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
