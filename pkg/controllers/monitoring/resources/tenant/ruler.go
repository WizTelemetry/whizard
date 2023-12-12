package tenant

import (
	"fmt"
	"strings"

	monitoringv1alpha1 "github.com/kubesphere/whizard/pkg/api/monitoring/v1alpha1"
	"github.com/kubesphere/whizard/pkg/constants"
	"github.com/kubesphere/whizard/pkg/util"
	"github.com/prometheus-operator/prometheus-operator/pkg/k8sutil"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
	"k8s.io/utils/pointer"
)

type RuleLevel string

const (
	RuleLevelNamesapce RuleLevel = "namespace"
	RuleLevelCluster   RuleLevel = "cluster"
	RuleLevelGlobal    RuleLevel = "global"

	// label keys in PrometheusRule.metadata.labels
	PrometheusRuleResourceLabelKeyOwnerNamespace = "alerting.kubesphere.io/owner_namespace"
	PrometheusRuleResourceLabelKeyOwnerCluster   = "alerting.kubesphere.io/owner_cluster"
	PrometheusRuleResourceLabelKeyRuleLevel      = "alerting.kubesphere.io/rule_level"
)

func (t *Tenant) ruler() error {

	ruler := &monitoringv1alpha1.Ruler{}
	if t.tenant.Status.Ruler != nil {
		err := t.Client.Get(t.Context, types.NamespacedName{
			Namespace: t.tenant.Status.Ruler.Namespace,
			Name:      t.tenant.Status.Ruler.Name,
		}, ruler)
		if err != nil {
			if apierrors.IsNotFound(err) {
				klog.V(3).Infof("Cannot find ruler [%s] for tenant [%s], create one", t.tenant.Status.Ruler, t.tenant.Name)
			} else {
				return err
			}
		}
	}

	// when tenant.Labels don't contain Service, remove the bindings to ingester and ruler
	if v, ok := t.tenant.Labels[constants.ServiceLabelKey]; !ok || v == "" {
		klog.V(3).Infof("Tenant [%s]'s Service is empty. ruler does not need to be created", t.tenant.Name)
		if t.tenant.Status.Ruler != nil && ruler != nil {
			if err := t.Client.Delete(t.Context, ruler); err != nil {
				return err
			}
			t.tenant.Status.Ruler = nil
			return t.Client.Status().Update(t.Context, t.tenant)
		}
		return nil
	}
	ruler = t.createOrUpdateRulerinstance()
	t.tenant.Status.Ruler = &monitoringv1alpha1.ObjectReference{
		Namespace: ruler.Namespace,
		Name:      ruler.Name,
	}
	err := util.CreateOrUpdate(t.Context, t.Client, ruler)
	if err != nil {
		return err
	}
	return t.Client.Status().Update(t.Context, t.tenant)
}

func (t *Tenant) createOrUpdateRulerinstance() *monitoringv1alpha1.Ruler {

	label := make(map[string]string, 2)
	label[constants.ServiceLabelKey] = t.tenant.Labels[constants.ServiceLabelKey]

	serviceNamespacedName := strings.Split(t.tenant.Labels[constants.ServiceLabelKey], ".")

	var ruleSelectors []*metav1.LabelSelector
	// add default rule selectors. (mainly used to select recording rules)
	ruleSelectors = append(ruleSelectors, t.Options.Ruler.RuleSelectors...)
	if t.Options.Ruler.DisableAlertingRulesAutoSelection == nil ||
		!*t.Options.Ruler.DisableAlertingRulesAutoSelection {
		// select alerting rules associated with this tenant(cluster)
		ruleSelectors = append(ruleSelectors, &metav1.LabelSelector{
			MatchLabels: map[string]string{
				PrometheusRuleResourceLabelKeyOwnerCluster: t.tenant.Spec.Tenant,
			},
			MatchExpressions: []metav1.LabelSelectorRequirement{{
				Key:      PrometheusRuleResourceLabelKeyRuleLevel,
				Operator: metav1.LabelSelectorOpIn,
				Values:   []string{string(RuleLevelCluster), string(RuleLevelNamesapce)},
			}},
		})
	}

	rulerName := t.tenant.Name

	if len(rulerName) > 30 {
		rn := k8sutil.NewResourceNamerWithPrefix("")
		out, _ := rn.UniqueDNS1123Label(rulerName) // out is a string with 8-chars hash-based suffix and ignore the return err which is always nil
		hashSuffix := out[len(out)-8:]

		name := rulerName[:30]
		name = strings.Trim(name, "-")

		rulerName = fmt.Sprintf("%s-%s", name, hashSuffix)
	}

	return &monitoringv1alpha1.Ruler{ObjectMeta: metav1.ObjectMeta{
		Name:      rulerName,
		Namespace: serviceNamespacedName[0],
		Labels:    label,
		OwnerReferences: []metav1.OwnerReference{
			{
				APIVersion: t.tenant.APIVersion,
				Kind:       t.tenant.Kind,
				Name:       t.tenant.Name,
				UID:        t.tenant.UID,
				Controller: pointer.BoolPtr(true),
			},
		},
	},
		Spec: monitoringv1alpha1.RulerSpec{
			Tenant: t.tenant.Spec.Tenant,
			// Only set RuleSelectors. The RuleNamespaceSelector is nil and will use the ruler's namespace
			RuleSelectors: ruleSelectors,
		},
	}
}
