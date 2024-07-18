package query

import (
	"fmt"
	"strings"

	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/discovery/targetgroup"
	"gopkg.in/yaml.v3"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"github.com/WhizardTelemetry/whizard/pkg/controllers/resources"
	"github.com/WhizardTelemetry/whizard/pkg/util"
)

func (q *Query) proxyConfigMap() (runtime.Object, resources.Operation, error) {

	var cm = &corev1.ConfigMap{ObjectMeta: q.meta(q.name("proxy-config"))}

	if q.query == nil {
		return cm, resources.OperationDelete, nil
	}

	stores, err := q.stores()
	if err != nil {
		return nil, resources.OperationCreateOrUpdate, err
	}
	data, err := envoyConfigFiles(q.Service.Namespace+"/"+q.Service.Name, *stores)
	if err != nil {
		return nil, resources.OperationCreateOrUpdate, err
	}
	cm.Data = data

	return cm, resources.OperationCreateOrUpdate, ctrl.SetControllerReference(q.query, cm, q.Scheme)
}

func (q *Query) storesConfigMap() (runtime.Object, resources.Operation, error) {
	var cm = &corev1.ConfigMap{ObjectMeta: q.meta(q.name("stores-config"))}
	if q.query == nil {
		return cm, resources.OperationDelete, nil
	}

	stores, err := q.stores()
	if err != nil {
		return nil, resources.OperationCreateOrUpdate, err
	}

	var targets []model.LabelSet
	for _, store := range stores.DirectStores {
		targets = append(targets, model.LabelSet{
			model.AddressLabel: model.LabelValue(store.Address),
		})
	}
	for _, store := range stores.ProxyStores {
		targets = append(targets, model.LabelSet{
			model.AddressLabel: model.LabelValue(fmt.Sprintf("%s:%d", store.ListenHost, store.ListenPort)),
		})
	}

	out, err := yaml.Marshal([]targetgroup.Group{{Targets: targets}})
	if err != nil {
		return nil, resources.OperationCreateOrUpdate, err
	}

	cm.Data = map[string]string{
		storesFile: string(out),
	}

	return cm, resources.OperationCreateOrUpdate, ctrl.SetControllerReference(q.query, cm, q.Scheme)
}

func (q *Query) envoyConfigMap(data map[string]string) error {
	var cm = &corev1.ConfigMap{ObjectMeta: q.meta(q.name("envoy-config"))}

	var buff strings.Builder
	tmpl := util.EnvoyStaticConfigTemplate
	if err := tmpl.Execute(&buff, data); err != nil {
		return err
	}

	cm.Data = map[string]string{
		envoyConfigFile: buff.String(),
	}

	if err := ctrl.SetControllerReference(q.query, cm, q.Scheme); err != nil {
		return err
	}
	_, err := controllerutil.CreateOrPatch(q.Context, q.Client, cm, configmapDataMutate(cm, cm.Data))
	return err
}

func configmapDataMutate(cm *corev1.ConfigMap, data map[string]string) controllerutil.MutateFn {
	return func() error {
		cm.Data = data
		return nil
	}
}
