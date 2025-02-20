package resources

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"

	monitoringv1alpha1 "github.com/WhizardTelemetry/whizard/pkg/api/monitoring/v1alpha1"
	"github.com/WhizardTelemetry/whizard/pkg/constants"
)

// Apply defaults to the service
func ApplyDefaults(service *monitoringv1alpha1.Service) *monitoringv1alpha1.Service {
	var whizardDefaultReplicas int32 = 2
	var whizardCompactorReplicas int32 = 1
	var whizardRulerReplicas int32 = 1
	var whizardQueryReplicas int32 = 3
	var whizardRouterReplicationFactor uint64 = 1

	if service.Spec.CompactorTemplateSpec.Image == "" {
		service.Spec.CompactorTemplateSpec.Image = constants.DefaultWhizardBaseImage
	}
	if service.Spec.CompactorTemplateSpec.Replicas == nil {
		service.Spec.CompactorTemplateSpec.Replicas = &whizardCompactorReplicas
	}
	if service.Spec.CompactorTemplateSpec.Resources.Size() == 0 {
		service.Spec.CompactorTemplateSpec.Resources = constants.ComponentResourcePresetMedium
	}

	if service.Spec.GatewayTemplateSpec.Image == "" {
		service.Spec.GatewayTemplateSpec.Image = constants.DefaultWhizardMonitoringGatewayImage
	}
	if service.Spec.GatewayTemplateSpec.Replicas == nil {
		service.Spec.GatewayTemplateSpec.Replicas = &whizardDefaultReplicas
	}
	if service.Spec.GatewayTemplateSpec.Resources.Size() == 0 {
		service.Spec.GatewayTemplateSpec.Resources = constants.ComponentResourcePresetMedium
	}

	if service.Spec.IngesterTemplateSpec.Image == "" {
		service.Spec.IngesterTemplateSpec.Image = constants.DefaultWhizardBaseImage
	}
	if service.Spec.IngesterTemplateSpec.Replicas == nil {
		service.Spec.IngesterTemplateSpec.Replicas = &whizardDefaultReplicas
	}
	if service.Spec.IngesterTemplateSpec.Resources.Size() == 0 {
		service.Spec.IngesterTemplateSpec.Resources = constants.ComponentResourcePresetLarge
	}
	if service.Spec.IngesterTemplateSpec.IngesterTSDBCleanUp.Image == "" {
		service.Spec.IngesterTemplateSpec.IngesterTSDBCleanUp.Image = constants.DefaultIngesterTSDBCleanupImage
	}
	if service.Spec.IngesterTemplateSpec.IngesterTSDBCleanUp.Resources.Size() == 0 {
		service.Spec.IngesterTemplateSpec.IngesterTSDBCleanUp.Resources = constants.ComponentResourcePresetMedium
	}
	if service.Spec.IngesterTemplateSpec.LocalTsdbRetention == "" {
		service.Spec.IngesterTemplateSpec.LocalTsdbRetention = "7d"
	}

	if service.Spec.QueryFrontendTemplateSpec.Image == "" {
		service.Spec.QueryFrontendTemplateSpec.Image = constants.DefaultWhizardBaseImage
	}
	if service.Spec.QueryFrontendTemplateSpec.Replicas == nil {
		service.Spec.QueryFrontendTemplateSpec.Replicas = &whizardDefaultReplicas
	}
	if service.Spec.QueryFrontendTemplateSpec.Resources.Size() == 0 {
		service.Spec.QueryFrontendTemplateSpec.Resources = constants.ComponentResourcePresetMedium
	}

	if service.Spec.QueryTemplateSpec.Image == "" {
		service.Spec.QueryTemplateSpec.Image = constants.DefaultWhizardBaseImage
	}
	if service.Spec.QueryTemplateSpec.Replicas == nil {
		service.Spec.QueryTemplateSpec.Replicas = &whizardQueryReplicas
	}
	if service.Spec.QueryTemplateSpec.Resources.Size() == 0 {
		service.Spec.QueryTemplateSpec.Resources = constants.ComponentResourcePresetMedium
	}
	if service.Spec.QueryTemplateSpec.Envoy.Image == "" {
		service.Spec.QueryTemplateSpec.Envoy.Image = constants.DefaultEnvoyImage
	}
	if service.Spec.QueryTemplateSpec.Envoy.Resources.Size() == 0 {
		service.Spec.QueryTemplateSpec.Envoy.Resources = constants.ComponentResourcePresetMedium
	}

	if service.Spec.RulerTemplateSpec.Image == "" {
		service.Spec.RulerTemplateSpec.Image = constants.DefaultWhizardBaseImage
	}
	if service.Spec.RulerTemplateSpec.Replicas == nil {
		service.Spec.RulerTemplateSpec.Replicas = &whizardRulerReplicas
	}
	if service.Spec.RulerTemplateSpec.Resources.Size() == 0 {
		service.Spec.RulerTemplateSpec.Resources = constants.ComponentResourcePresetMedium
	}
	if service.Spec.RulerTemplateSpec.PrometheusConfigReloader.Image == "" {
		service.Spec.RulerTemplateSpec.PrometheusConfigReloader.Image = constants.DefaultPrometheusConfigReloaderImage
	}
	if service.Spec.RulerTemplateSpec.PrometheusConfigReloader.Resources.Size() == 0 {
		service.Spec.RulerTemplateSpec.PrometheusConfigReloader.Resources = corev1.ResourceRequirements{
			Requests: corev1.ResourceList{
				corev1.ResourceCPU:    resource.MustParse("100m"),
				corev1.ResourceMemory: resource.MustParse("128Mi"),
			},
			Limits: corev1.ResourceList{
				corev1.ResourceCPU:    resource.MustParse("100m"),
				corev1.ResourceMemory: resource.MustParse("128Mi"),
			},
		}
	}
	if service.Spec.RulerTemplateSpec.RulerQueryProxy.Image == "" {
		service.Spec.RulerTemplateSpec.RulerQueryProxy.Image = constants.DefaultWhizardMonitoringGatewayImage
	}
	if service.Spec.RulerTemplateSpec.RulerQueryProxy.Resources.Size() == 0 {
		service.Spec.RulerTemplateSpec.RulerQueryProxy.Resources = constants.ComponentResourcePresetMedium
	}

	if service.Spec.RulerTemplateSpec.RulerWriteProxy.Image == "" {
		service.Spec.RulerTemplateSpec.RulerWriteProxy.Image = constants.DefaultRulerWriteProxyImage
	}
	if service.Spec.RulerTemplateSpec.RulerWriteProxy.Resources.Size() == 0 {
		service.Spec.RulerTemplateSpec.RulerWriteProxy.Resources = constants.ComponentResourcePresetMedium
	}

	if service.Spec.RouterTemplateSpec.Image == "" {
		service.Spec.RouterTemplateSpec.Image = constants.DefaultWhizardBaseImage
	}
	if service.Spec.RouterTemplateSpec.Replicas == nil {
		service.Spec.RouterTemplateSpec.Replicas = &whizardDefaultReplicas
	}
	if service.Spec.RouterTemplateSpec.Resources.Size() == 0 {
		service.Spec.RouterTemplateSpec.Resources = constants.ComponentResourcePresetMedium
	}
	if service.Spec.RouterTemplateSpec.ReplicationFactor == nil {
		service.Spec.RouterTemplateSpec.ReplicationFactor = &whizardRouterReplicationFactor
	}

	if service.Spec.StoreTemplateSpec.Image == "" {
		service.Spec.StoreTemplateSpec.Image = constants.DefaultWhizardBaseImage
	}
	if service.Spec.StoreTemplateSpec.Replicas == nil {
		service.Spec.StoreTemplateSpec.Replicas = &whizardDefaultReplicas
	}
	if service.Spec.StoreTemplateSpec.Resources.Size() == 0 {
		service.Spec.StoreTemplateSpec.Resources = constants.ComponentResourcePresetLarge
	}

	return service
}
