package router

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/thanos-io/thanos/pkg/receive"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	monitoringv1alpha1 "github.com/WhizardTelemetry/whizard/pkg/api/monitoring/v1alpha1"
	"github.com/WhizardTelemetry/whizard/pkg/constants"
	"github.com/WhizardTelemetry/whizard/pkg/controllers/resources"
	"github.com/WhizardTelemetry/whizard/pkg/controllers/resources/ingester"
	"github.com/WhizardTelemetry/whizard/pkg/util"
)

func (r *Router) hashringsConfigMap() (runtime.Object, resources.Operation, error) {

	var cm = &corev1.ConfigMap{ObjectMeta: r.meta(r.name("hashrings-config"))}

	if r.router == nil {
		return cm, resources.OperationDelete, nil
	}

	var hashrings []receive.HashringConfig
	var softHashring = receive.HashringConfig{
		Hashring:  "softs",
		Endpoints: []receive.Endpoint{},
	}

	var ingesterList monitoringv1alpha1.IngesterList
	if err := r.Client.List(r.Context, &ingesterList,
		client.MatchingLabels(util.ManagedLabelByService(r.Service))); err != nil {

		r.Log.WithValues("ingesterlist", "").Error(err, "")
		return nil, resources.OperationCreateOrUpdate, err
	}

	for _, item := range ingesterList.Items {

		ingester, err := ingester.New(r.BaseReconciler, &item)
		if err != nil {
			return nil, "", err
		}
		if len(item.Spec.Tenants) == 0 {
			// specific ingesters can join the softHashring
			if v, ok := item.Labels[constants.SoftTenantLabelKey]; ok && v == "true" {
				for _, addr := range ingester.Address() {
					ep := receive.Endpoint{
						Address:          fmt.Sprintf("%s:%d", addr, constants.GRPCPort),
						CapNProtoAddress: fmt.Sprintf("%s:%d", addr, constants.CapNProtoPort),
					}
					softHashring.Endpoints = append(softHashring.Endpoints, ep)
				}
			}
			continue
		}

		var endpoints []receive.Endpoint
		for _, addr := range ingester.Address() {
			ep := receive.Endpoint{
				Address:          fmt.Sprintf("%s:%d", addr, constants.GRPCPort),
				CapNProtoAddress: fmt.Sprintf("%s:%d", addr, constants.CapNProtoPort),
			}
			endpoints = append(endpoints, ep)
		}

		hashrings = append(hashrings, receive.HashringConfig{
			Hashring:  item.Namespace + "/" + item.Name,
			Tenants:   item.Spec.Tenants,
			Endpoints: endpoints,
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
