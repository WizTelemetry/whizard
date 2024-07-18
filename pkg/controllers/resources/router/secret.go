package router

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/WhizardTelemetry/whizard/pkg/constants"
	"github.com/WhizardTelemetry/whizard/pkg/controllers/resources"
)

func (r *Router) webConfigSecret() (runtime.Object, resources.Operation, error) {
	var secret = &corev1.Secret{ObjectMeta: r.meta(r.name("web-config"))}

	if r.router == nil {
		return secret, resources.OperationDelete, nil
	}

	if r.router.Spec.WebConfig == nil {
		return secret, resources.OperationDelete, nil
	}

	body, err := r.BaseReconciler.CreateWebConfig(r.router.Namespace, r.router.Spec.WebConfig)
	if err != nil {
		return nil, resources.OperationDelete, err
	}

	secret.Data = map[string][]byte{
		constants.WhizardWebConfigFile: body,
	}

	return secret, resources.OperationCreateOrUpdate, ctrl.SetControllerReference(r.router, secret, r.Scheme)
}
