package ruler

import (
	"fmt"

	"github.com/kubesphere/whizard/pkg/constants"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/pointer"

	monitoringv1alpha1 "github.com/kubesphere/whizard/pkg/api/monitoring/v1alpha1"
	"github.com/kubesphere/whizard/pkg/controllers/monitoring/options"
	"github.com/kubesphere/whizard/pkg/controllers/monitoring/resources"
)

const (
	configDir  = "/etc/whizard"
	rulesDir   = configDir + "/rules"
	storageDir = "/whizard"
)

type Ruler struct {
	resources.BaseReconciler
	ruler                 *monitoringv1alpha1.Ruler
	reloaderConfig        options.PrometheusConfigReloaderConfig
	rulerQueryProxyConfig options.RulerQueryProxyConfig
}

func New(reconciler resources.BaseReconciler, ruler *monitoringv1alpha1.Ruler,
	reloaderConfig options.PrometheusConfigReloaderConfig, rulerQueryProxyConfig options.RulerQueryProxyConfig) *Ruler {

	return &Ruler{
		BaseReconciler:        reconciler,
		ruler:                 ruler,
		reloaderConfig:        reloaderConfig,
		rulerQueryProxyConfig: rulerQueryProxyConfig,
	}
}

func (r *Ruler) labels() map[string]string {
	labels := r.BaseLabels()
	labels[constants.LabelNameAppName] = constants.AppNameRuler
	labels[constants.LabelNameAppManagedBy] = r.ruler.Name
	return labels
}

func (r *Ruler) name(nameSuffix ...string) string {
	return resources.QualifiedName(constants.AppNameRuler, r.ruler.Name, nameSuffix...)
}

func (r *Ruler) meta(name string) metav1.ObjectMeta {

	return metav1.ObjectMeta{
		Name:            name,
		Namespace:       r.ruler.Namespace,
		Labels:          r.labels(),
		OwnerReferences: r.OwnerReferences(),
	}
}

func (r *Ruler) OwnerReferences() []metav1.OwnerReference {
	return []metav1.OwnerReference{
		{
			APIVersion: r.ruler.APIVersion,
			Kind:       r.ruler.Kind,
			Name:       r.ruler.Name,
			UID:        r.ruler.UID,
			Controller: pointer.BoolPtr(true),
		},
	}
}

func (r *Ruler) HttpAddr() string {
	return fmt.Sprintf("http://%s.%s.svc:%d",
		r.name(constants.ServiceNameSuffix), r.ruler.Namespace, constants.HTTPPort)
}

func (r *Ruler) Reconcile() error {
	createOrUpdateCms, deleteCms, useCms, err := r.ruleConfigMaps()
	if err != nil {
		return err
	}

	var ress []resources.Resource
	for _, cm := range createOrUpdateCms {
		ruleCm := cm
		ress = append(ress, func() (runtime.Object, resources.Operation, error) {
			return &ruleCm, resources.OperationCreateOrUpdate, nil
		})
	}
	for _, cm := range deleteCms {
		ruleCm := cm
		ress = append(ress, func() (runtime.Object, resources.Operation, error) {
			return &ruleCm, resources.OperationDelete, nil
		})
	}
	var ruleConfigMapNames []string
	for _, cm := range useCms {
		ruleConfigMapNames = append(ruleConfigMapNames, cm.Name)
	}

	ress = append(ress, func() (runtime.Object, resources.Operation, error) {
		return r.statefulSet(ruleConfigMapNames)
	})

	return r.ReconcileResources(append(
		ress,
		r.service))
}
