package queryfrontend

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/WhizardTelemetry/whizard/pkg/constants"
	"github.com/WhizardTelemetry/whizard/pkg/controllers/resources"
)

func (q *QueryFrontend) webConfigSecret() (runtime.Object, resources.Operation, error) {
	var secret = &corev1.Secret{ObjectMeta: q.meta(q.name("web-config"))}

	if q.queryFrontend == nil {
		return secret, resources.OperationDelete, nil
	}

	if q.queryFrontend.Spec.WebConfig == nil {
		return secret, resources.OperationDelete, nil
	}

	body, err := q.BaseReconciler.CreateWebConfig(q.queryFrontend.Namespace, q.queryFrontend.Spec.WebConfig)
	if err != nil {
		return nil, resources.OperationDelete, err
	}

	secret.Data = map[string][]byte{
		constants.WhizardWebConfigFile: body,
	}

	return secret, resources.OperationCreateOrUpdate, ctrl.SetControllerReference(q.queryFrontend, secret, q.Scheme)
}
