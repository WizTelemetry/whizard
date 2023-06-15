package router

import (
	"encoding/json"
	"strings"

	monitoringv1alpha1 "github.com/kubesphere/whizard/pkg/api/monitoring/v1alpha1"
	"github.com/kubesphere/whizard/pkg/constants"
	"github.com/kubesphere/whizard/pkg/controllers/monitoring/resources"
	"github.com/kubesphere/whizard/pkg/controllers/monitoring/resources/ingester"
	"github.com/kubesphere/whizard/pkg/util"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
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
		r.Options.Ingester.Override(&item.Spec)
		ingester, err := ingester.New(r.BaseReconciler, &item, r.Options.Ingester)
		if err != nil {
			return nil, "", err
		}
		if len(item.Spec.Tenants) == 0 {
			// specific ingesters can join the softHashring
			if v, ok := item.Labels[constants.SoftTenantLabelKey]; ok && v == "true" {
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

func (r *Router) envoyConfigMap(data map[string]string) error {
	var cm = &corev1.ConfigMap{ObjectMeta: r.meta(r.name("envoy-config"))}

	var buff strings.Builder
	tmpl := util.EnvoyStaticConfigTemplate
	if err := tmpl.Execute(&buff, data); err != nil {
		return err
	}

	cm.Data = map[string]string{
		envoyConfigFile: buff.String(),
	}

	if err := ctrl.SetControllerReference(r.router, cm, r.Scheme); err != nil {
		return err
	}
	_, err := controllerutil.CreateOrPatch(r.Context, r.Client, cm, configmapDataMutate(cm, cm.Data))
	return err
}

func configmapDataMutate(cm *corev1.ConfigMap, data map[string]string) controllerutil.MutateFn {
	return func() error {
		cm.Data = data
		return nil
	}
}
