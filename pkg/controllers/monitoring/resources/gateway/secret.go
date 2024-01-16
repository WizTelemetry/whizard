package gateway

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/kubesphere/whizard/pkg/constants"
	"github.com/kubesphere/whizard/pkg/controllers/monitoring/resources"
)

const (
	TLSVersionTLS12 = "TLS12"
	TLSVersionTLS13 = "TLS13"
)

func (g *Gateway) webConfigSecret() (runtime.Object, resources.Operation, error) {
	var secret = &corev1.Secret{ObjectMeta: g.meta(g.name("web-config"))}

	if g.gateway == nil {
		return secret, resources.OperationDelete, nil
	}

	if g.gateway.Spec.WebConfig == nil {
		return secret, resources.OperationDelete, nil
	}

	body, err := g.BaseReconciler.CreateWebConfig(g.gateway.Namespace, g.gateway.Spec.WebConfig)
	if err != nil {
		return nil, resources.OperationDelete, err
	}

	secret.Data = map[string][]byte{
		constants.WhizardWebConfigFile: body,
	}

	return secret, resources.OperationCreateOrUpdate, ctrl.SetControllerReference(g.gateway, secret, g.Scheme)
}
