package apis

import (
	"github.com/WhizardTelemetry/whizard/pkg/api/monitoring/v1alpha1"
)

func init() {
	AddToSchemes = append(AddToSchemes, v1alpha1.SchemeBuilder.AddToScheme)
}
