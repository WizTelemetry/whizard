package gateway

import (
	"encoding/json"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/kubesphere/whizard/pkg/api/monitoring/v1alpha1"
	"github.com/kubesphere/whizard/pkg/constants"
	"github.com/kubesphere/whizard/pkg/controllers/monitoring/resources"
	monitoringgateway "github.com/kubesphere/whizard/pkg/monitoring-gateway"
)

const (
	tenantsAdmissionConfigFile = "tenants-admission.yaml"
	webConfigFile              = "web-config.yaml"
)

func (g *Gateway) tenantsAdmissionConfigMap() (runtime.Object, resources.Operation, error) {

	var cm = &corev1.ConfigMap{ObjectMeta: g.meta(g.name("tenants-admission-config"))}

	if g.gateway == nil {
		return cm, resources.OperationDelete, nil
	}

	if !g.gateway.Spec.EnabledTenantsAdmission {
		return cm, resources.OperationDelete, nil
	}

	acConfig := monitoringgateway.AdmissionControlConfig{}
	tenantList := &v1alpha1.TenantList{}
	err := g.Client.List(g.Context, tenantList)
	if err != nil {
		return nil, resources.OperationCreateOrUpdate, err
	}

	for _, tenant := range tenantList.Items {
		if tenant.GetDeletionTimestamp().IsZero() {
			if v, ok := tenant.Labels[constants.ServiceLabelKey]; ok && g.gateway.Labels[constants.ServiceLabelKey] == v {
				acConfig.Tenants = append(acConfig.Tenants, tenant.Spec.Tenant)
			}
		}
	}

	acBytes, err := json.Marshal(acConfig)
	if err != nil {
		return nil, resources.OperationCreateOrUpdate, err
	}
	cm.Data = map[string]string{
		tenantsAdmissionConfigFile: string(acBytes),
	}

	return cm, resources.OperationCreateOrUpdate, ctrl.SetControllerReference(g.gateway, cm, g.Scheme)
}
