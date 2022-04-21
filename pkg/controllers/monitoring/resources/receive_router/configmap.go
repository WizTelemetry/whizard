package receive_router

import (
	"encoding/json"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	monitoringv1alpha1 "github.com/kubesphere/paodin-monitoring/pkg/api/monitoring/v1alpha1"
	"github.com/kubesphere/paodin-monitoring/pkg/controllers/monitoring/resources"
	"github.com/kubesphere/paodin-monitoring/pkg/controllers/monitoring/resources/receive_ingestor"
)

func (r *ReceiveRouter) hashringsConfigMap() (runtime.Object, resources.Operation, error) {

	var cm = &corev1.ConfigMap{ObjectMeta: r.meta(r.name("hashrings-config"))}

	if r.router == nil {
		return cm, resources.OperationDelete, nil
	}

	type HashringConfig struct {
		Hashring  string   `json:"hashring,omitempty"`
		Tenants   []string `json:"tenants,omitempty"`
		Endpoints []string `json:"endpoints"`
	}
	var hashrings = []HashringConfig{}
	var softHashring = HashringConfig{
		Hashring:  "softs",
		Endpoints: []string{},
	}

	var ingestorList monitoringv1alpha1.ThanosReceiveIngestorList
	if err := r.Client.List(r.Context, &ingestorList,
		client.MatchingLabels(monitoringv1alpha1.ManagedLabelByService(r.Service))); err != nil {

		r.Log.WithValues("thanosreceiveingestorlist", "").Error(err, "")
		return nil, resources.OperationCreateOrUpdate, err
	}

	for _, item := range ingestorList.Items {
		ingestor := receive_ingestor.New(r.BaseReconciler, &item)
		if len(item.Spec.Tenants) == 0 {
			softHashring.Endpoints = append(softHashring.Endpoints, ingestor.GrpcAddrs()...)
			continue
		}
		hashrings = append(hashrings, HashringConfig{
			Hashring:  item.Namespace + "/" + item.Name,
			Tenants:   item.Spec.Tenants,
			Endpoints: ingestor.GrpcAddrs(),
		})
	}
	hashrings = append(hashrings, softHashring) // put soft tenants at the end

	hashringBytes, err := json.MarshalIndent(hashrings, "", "\t")
	if err != nil {
		return nil, resources.OperationCreateOrUpdate, err
	}
	cm.Data = map[string]string{
		hashringsFile: string(hashringBytes),
	}

	return cm, resources.OperationCreateOrUpdate, nil
}
