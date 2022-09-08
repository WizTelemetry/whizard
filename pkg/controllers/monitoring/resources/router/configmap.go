package router

import (
	"encoding/json"

	monitoringv1alpha1 "github.com/kubesphere/whizard/pkg/api/monitoring/v1alpha1"
	"github.com/kubesphere/whizard/pkg/constants"
	"github.com/kubesphere/whizard/pkg/controllers/monitoring/resources"
	"github.com/kubesphere/whizard/pkg/controllers/monitoring/resources/ingester"
	"github.com/kubesphere/whizard/pkg/util"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (r *Router) hashringsConfigMap() (runtime.Object, resources.Operation, error) {

	var cm = &corev1.ConfigMap{ObjectMeta: r.meta(r.name("hashrings-config"))}

	if r.router == nil {
		return cm, resources.OperationDelete, nil
	}

	type HashringConfig struct {
		Hashring  string   `json:"hashring,omitempty"`
		Tenants   []string `json:"tenants,omitempty"`
		Endpoints []string `json:"endpoints"`
	}
	var hashrings []HashringConfig
	var softHashring = HashringConfig{
		Hashring:  "softs",
		Endpoints: []string{},
	}

	var ingesterList monitoringv1alpha1.IngesterList
	if err := r.Client.List(r.Context, &ingesterList,
		client.MatchingLabels(util.ManagedLabelByService(r.Service))); err != nil {

		r.Log.WithValues("ingesterlist", "").Error(err, "")
		return nil, resources.OperationCreateOrUpdate, err
	}

	for _, item := range ingesterList.Items {
		ingester := ingester.New(r.BaseReconciler, &item)
		if len(item.Spec.Tenants) == 0 {
			// the ingester in the "deleting" state will not be added to the soft hash ring
			if v, ok := item.ObjectMeta.Labels[constants.LabelNameIngesterState]; !ok || v != constants.IngesterStateDeleting {
				softHashring.Endpoints = append(softHashring.Endpoints, ingester.GrpcAddrs()...)
			}
			continue
		}
		hashrings = append(hashrings, HashringConfig{
			Hashring:  item.Namespace + "/" + item.Name,
			Tenants:   item.Spec.Tenants,
			Endpoints: ingester.GrpcAddrs(),
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

	return cm, resources.OperationCreateOrUpdate, ctrl.SetControllerReference(r.router, cm, r.Scheme)
}
