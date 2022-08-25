package ruler

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/kubesphere/whizard/pkg/constants"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

var (
	sequenceNumberRegexp = regexp.MustCompile(`^([1-9]\d*|0)$`)
)

type Ruler struct {
	resources.BaseReconciler
	ruler                    *monitoringv1alpha1.Ruler
	reloaderConfig           options.PrometheusConfigReloaderConfig
	rulerQueryProxyConfig    options.RulerQueryProxyConfig
	shardsRuleConfigMapNames []map[string]struct{} // rule configmaps for each shard
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

func (r *Ruler) HttpAddrs() []string {
	var addrs []string
	for shardSn := 0; shardSn < int(*r.ruler.Spec.Shards); shardSn++ {
		addrs = append(addrs, fmt.Sprintf("http://%s.%s.svc:%d",
			r.name(strconv.Itoa(shardSn), constants.ServiceNameSuffix), r.ruler.Namespace, constants.HTTPPort))
	}
	return addrs
}

func (r *Ruler) GrpcAddrs() []string {
	var addrs []string
	for shardSn := 0; shardSn < int(*r.ruler.Spec.Shards); shardSn++ {
		addrs = append(addrs, fmt.Sprintf("%s.%s.svc:%d",
			r.name(strconv.Itoa(shardSn), constants.ServiceNameSuffix), r.ruler.Namespace, constants.GRPCPort))
	}
	return addrs
}

func (r *Ruler) Reconcile() error {
	var ress []resources.Resource
	ress = append(ress, r.ruleConfigMaps()...)
	ress = append(ress, r.statefulSets()...)
	ress = append(ress, r.services()...)

	return r.ReconcileResources(ress)
}
