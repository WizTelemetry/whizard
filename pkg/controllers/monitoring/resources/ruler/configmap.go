package ruler

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/prometheus-operator/prometheus-operator/pkg/prometheus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	monitoringv1alpha1 "github.com/kubesphere/paodin/pkg/api/monitoring/v1alpha1"
)

var maxConfigMapDataSize = int(float64(corev1.MaxSecretSize) * 0.5)

func (r *Ruler) ruleConfigMaps() (createOrUpdates []corev1.ConfigMap, deletes []corev1.ConfigMap, uses []corev1.ConfigMap, err error) {
	var prometheusRules = make(map[string]string)
	prometheusRules, err = r.selectPrometheusRules()
	if err != nil {
		return nil, nil, nil, err
	}

	var rules = make(map[string]string)
	rules, err = r.selectRules()
	if err != nil {
		return
	}

	for k, v := range rules {
		prometheusRules[k] = v
	}

	uses, err = r.makeRulesConfigMaps(prometheusRules)
	if err != nil {
		return
	}
	usesMap := make(map[string]struct{})
	for _, useCm := range uses {
		usesMap[useCm.Name] = struct{}{}
	}

	var cmList corev1.ConfigMapList
	err = r.Client.List(r.Context, &cmList, client.InNamespace(r.ruler.Namespace), client.MatchingLabels(r.labels()))
	if err != nil {
		return nil, nil, nil, err
	}
	currentsMap := make(map[string]corev1.ConfigMap)
	namePrefix := r.name("rulefiles")
	for _, currentCm := range cmList.Items {
		if strings.HasPrefix(currentCm.Name, namePrefix) {
			cm := currentCm
			currentsMap[cm.Name] = cm
			if _, ok := usesMap[cm.Name]; !ok {
				deletes = append(deletes, cm)
			}
		}
	}

	for _, useCm := range uses {
		if currentCm, ok := currentsMap[useCm.Name]; !ok || !reflect.DeepEqual(useCm.Data, currentCm.Data) {
			cm := useCm
			createOrUpdates = append(createOrUpdates, cm)
		}
	}

	return

}

type rulesWrapper struct {
	rules []promv1.Rule
	by    func(r1, r2 *promv1.Rule) bool
}

func (rw rulesWrapper) Len() int {
	return len(rw.rules)
}

func (rw rulesWrapper) Swap(i, j int) {
	rw.rules[i], rw.rules[j] = rw.rules[j], rw.rules[i]
}

func (rw rulesWrapper) Less(i, j int) bool {
	return rw.by(&rw.rules[i], &rw.rules[j])
}

// select rule resources and combine them to PrometheusRules (defined by prometheus-operator) struct
func (r *Ruler) selectRules() (map[string]string, error) {
	rules := make(map[string]string)

	namespaces, err := r.selectNamespaces(r.ruler.Spec.RuleNamespaceSelector)
	if err != nil {
		return nil, err
	}
	ruleSelector, err := metav1.LabelSelectorAsSelector(r.ruler.Spec.RuleSelector)
	if err != nil {
		return nil, err
	}

	var defaultGroupName = "default"

	for _, ns := range namespaces {
		var ruleList monitoringv1alpha1.RuleList
		err = r.Client.List(r.Context, &ruleList,
			client.MatchingLabelsSelector{Selector: ruleSelector}, client.InNamespace(ns))
		if err != nil {
			return nil, err
		}

		var groupList monitoringv1alpha1.RuleGroupList
		err = r.Client.List(r.Context, &groupList, client.InNamespace(ns))
		if err != nil {
			return nil, err
		}
		var groups = make(map[string]*monitoringv1alpha1.RuleGroup)
		for _, group := range groupList.Items {
			groups[group.Name] = &group
		}

		// combine Rules to the RuleGroups(defined by prometheus-operator)
		var promGroups = make(map[string]*promv1.RuleGroup)
		var groupsChecked = make(map[string]struct{})
		var groupNames []string
		for _, rule := range ruleList.Items {
			if rule.Spec.Alert == "" && rule.Spec.Record == "" {
				r.Log.WithValues("rule", ns+"/"+rule.Name).V(2).Error(nil, "ignore the rule because both alert and record are empty")
				continue
			}
			var groupName = defaultGroupName
			if g := monitoringv1alpha1.RuleGroupName(&rule); g != "" {
				groupName = g
			}
			if _, ok := promGroups[groupName]; !ok {
				group, ok2 := groups[groupName]
				if groupName != defaultGroupName && !ok2 { // for the group does not exist
					if _, checked := groupsChecked[groupName]; !checked { // avoid to log not found err too many times for a same group
						groupsChecked[groupName] = struct{}{}
						r.Log.WithValues("rulegroup", ns+"/"+groupName).Error(err, "not found")
					}
					continue
				}
				promGroup := promv1.RuleGroup{Name: groupName}
				if group != nil {
					promGroup.Interval = group.Spec.Interval
					promGroup.PartialResponseStrategy = group.Spec.PartialResponseStrategy
				}
				promGroups[groupName] = &promGroup
				groupNames = append(groupNames, groupName)

			}
			promGroups[groupName].Rules = append(promGroups[groupName].Rules, promv1.Rule{
				Alert:       rule.Spec.Alert,
				Record:      rule.Spec.Record,
				Expr:        rule.Spec.Expr,
				Labels:      rule.Spec.Labels,
				Annotations: rule.Spec.Annotations,
				For:         string(rule.Spec.For),
			})
		}

		// sort rules in each group by rule name asc
		for _, groupName := range groupNames {
			sort.Sort(rulesWrapper{promGroups[groupName].Rules, func(r1, r2 *promv1.Rule) bool {
				if r1.Alert == r2.Alert { // for record rules
					return r1.Record < r2.Record
				}
				return r1.Alert < r2.Alert
			}})
		}

		// split rules in default group in which when there are too much rules
		const defaultGroupSize = 20
		if promGroup, ok := promGroups[defaultGroupName]; ok && len(promGroup.Rules) > defaultGroupSize {
			rules := promGroup.Rules
			promGroup.Rules = rules[:defaultGroupSize]
			promGroups[defaultGroupName] = promGroup
			for i := 1; ; i++ {
				g := &promv1.RuleGroup{
					Name:                    fmt.Sprintf("%s.%d", defaultGroupName, i),
					Interval:                promGroup.Interval,
					PartialResponseStrategy: promGroup.PartialResponseStrategy,
				}
				if len(rules) > defaultGroupSize*(i+1) {
					g.Rules = rules[defaultGroupSize*i : defaultGroupSize*(i+1)]
					promGroups[g.Name] = g
				} else {
					g.Rules = rules[defaultGroupSize*i:]
					promGroups[g.Name] = g
					break
				}
			}
		}

		sort.Strings(groupNames)

		// combine RuleGroups(prometheus-operator) into PrometheusRules(prometheus-operator)
		var promRule promv1.PrometheusRuleSpec
		var size, count int
		for _, groupName := range groupNames {
			if size > maxConfigMapDataSize*90/100 {
				content, err := yaml.Marshal(&promRule)
				if err != nil {
					return nil, errors.Wrap(err, "failed to marshal content")
				}
				rules[fmt.Sprintf("%s.rules.%d.yaml", ns, count)] = string(content)

				promRule = promv1.PrometheusRuleSpec{}
				size = 0
				count++
			}

			promGroup := promGroups[groupName]

			content, err := yaml.Marshal(promGroup)
			if err != nil {
				return nil, errors.Wrap(err, "failed to marshal content")
			}
			size += len(string(content))

			promRule.Groups = append(promRule.Groups, *promGroup)
			continue

		}
		if size > 0 && len(promRule.Groups) > 0 {
			content, err := yaml.Marshal(&promRule)
			if err != nil {
				return nil, errors.Wrap(err, "failed to marshal content")
			}
			rules[fmt.Sprintf("%s.rules.%d.yaml", ns, count)] = string(content)
		}

	}

	return rules, nil

}

func (r *Ruler) selectNamespaces(nsSelector *metav1.LabelSelector) ([]string, error) {
	namespaces := []string{}

	// If nsSelector is nil, only check own namespace.
	if nsSelector == nil {
		namespaces = append(namespaces, r.ruler.Namespace)
	} else {
		selector, err := metav1.LabelSelectorAsSelector(nsSelector)
		if err != nil {
			return namespaces, errors.Wrap(err, "convert rule namespace label selector to selector")
		}

		var nsList corev1.NamespaceList
		err = r.Client.List(r.Context, &nsList, client.MatchingLabelsSelector{Selector: selector})
		if err != nil {
			return nil, err
		}
		for _, ns := range nsList.Items {
			namespaces = append(namespaces, ns.Name)
		}
	}

	return namespaces, nil
}

func (r *Ruler) selectPrometheusRules() (map[string]string, error) {
	rules := make(map[string]string)

	namespaces, err := r.selectNamespaces(r.ruler.Spec.PrometheusRuleNamespaceSelector)
	if err != nil {
		return nil, err
	}
	ruleSelector, err := metav1.LabelSelectorAsSelector(r.ruler.Spec.PrometheusRuleSelector)
	if err != nil {
		return nil, err
	}

	for _, ns := range namespaces {
		var prometheusRuleList promv1.PrometheusRuleList
		err = r.Client.List(r.Context, &prometheusRuleList,
			client.MatchingLabelsSelector{Selector: ruleSelector}, client.InNamespace(ns))
		if err != nil {
			return nil, err
		}
		for _, promRule := range prometheusRuleList.Items {
			content, err := GenerateContent(promRule.Spec, r.BaseReconciler.Log)
			if err != nil {
				return nil, err
			}
			rules[fmt.Sprintf("%v-%v.yaml", promRule.Namespace, promRule.Name)] = content
		}
	}

	return rules, nil
}

// makeRulesConfigMaps takes a ThanosRuler configuration and rule files and
// returns a list of Kubernetes ConfigMaps to be later on mounted into the
// ThanosRuler instance.
// If the total size of rule files exceeds the Kubernetes ConfigMap limit,
// they are split up via the simple first-fit [1] bin packing algorithm. In the
// future this can be replaced by a more sophisticated algorithm, but for now
// simplicity should be sufficient.
// [1] https://en.wikipedia.org/wiki/Bin_packing_problem#First-fit_algorithm
//
// refer to prometheus-operator
func (r *Ruler) makeRulesConfigMaps(ruleFiles map[string]string) ([]corev1.ConfigMap, error) {
	//check if none of the rule files is too large for a single ConfigMap
	for filename, file := range ruleFiles {
		if len(file) > maxConfigMapDataSize {
			return nil, errors.Errorf(
				"rule file '%v' is too large for a single Kubernetes ConfigMap",
				filename,
			)
		}
	}

	buckets := []map[string]string{
		{},
	}
	currBucketIndex := 0

	// To make bin packing algorithm deterministic, sort ruleFiles filenames and
	// iterate over filenames instead of ruleFiles map (not deterministic).
	fileNames := []string{}
	for n := range ruleFiles {
		fileNames = append(fileNames, n)
	}
	sort.Strings(fileNames)

	for _, filename := range fileNames {
		// If rule file doesn't fit into current bucket, create new bucket.
		if bucketSize(buckets[currBucketIndex])+len(ruleFiles[filename]) > maxConfigMapDataSize {
			buckets = append(buckets, map[string]string{})
			currBucketIndex++
		}
		buckets[currBucketIndex][filename] = ruleFiles[filename]
	}

	ruleFileConfigMaps := []corev1.ConfigMap{}
	for i, bucket := range buckets {
		ruleFileConfigMaps = append(ruleFileConfigMaps, corev1.ConfigMap{
			ObjectMeta: r.meta(r.name("rulefiles", strconv.Itoa(i))),
			Data:       bucket,
		})
	}

	return ruleFileConfigMaps, nil
}

func bucketSize(bucket map[string]string) int {
	totalSize := 0
	for _, v := range bucket {
		totalSize += len(v)
	}

	return totalSize
}

// GenerateContent takes a PrometheusRuleSpec and generates the rule content
func GenerateContent(promRule promv1.PrometheusRuleSpec, log logr.Logger) (string, error) {
	content, err := yaml.Marshal(promRule)
	if err != nil {

		return "", errors.Wrap(err, "failed to marshal content")
	}
	errs := prometheus.ValidateRule(promRule)
	if len(errs) != 0 {
		const m = "Invalid rule"
		log.V(9).WithValues("msg", m, "content", content).Info("")
		for _, err := range errs {
			log.WithValues("msg", m).Error(err, "")
		}
		return "", errors.New(m)
	}
	return string(content), nil
}
