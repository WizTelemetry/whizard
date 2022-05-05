package apis

import (
	"github.com/kubesphere/paodin/pkg/api/monitoring/v1alpha1"
)

func init() {
	AddToSchemes = append(AddToSchemes, v1alpha1.SchemeBuilder.AddToScheme)
}
