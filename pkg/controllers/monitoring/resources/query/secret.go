package query

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/kubesphere/whizard/pkg/constants"
	"github.com/kubesphere/whizard/pkg/controllers/monitoring/resources"
)

func (q *Query) webConfigSecret() (runtime.Object, resources.Operation, error) {
	var secret = &corev1.Secret{ObjectMeta: q.meta(q.name("web-config"))}

	if q.query == nil {
		return secret, resources.OperationDelete, nil
	}

	if q.query.Spec.WebConfig == nil {
		return secret, resources.OperationDelete, nil
	}

	body, err := q.BaseReconciler.CreateWebConfig(q.query.Namespace, q.query.Spec.WebConfig)
	if err != nil {
		return nil, resources.OperationDelete, err
	}

	secret.Data = map[string][]byte{
		constants.WhizardWebConfigFile: body,
	}

	return secret, resources.OperationCreateOrUpdate, ctrl.SetControllerReference(q.query, secret, q.Scheme)
}
