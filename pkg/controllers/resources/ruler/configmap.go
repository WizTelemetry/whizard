package ruler

import (
	"crypto/md5"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/prometheus-operator/prometheus-operator/pkg/operator"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"github.com/WhizardTelemetry/whizard/pkg/controllers/resources"
	"github.com/WhizardTelemetry/whizard/pkg/util"
)

var maxConfigMapDataSize = int(float64(corev1.MaxSecretSize) * 0.5)

var errResourcesFunc = func(err error) []resources.Resource {
	return []resources.Resource{
		func() (runtime.Object, resources.Operation, error) {
			return nil, resources.OperationCreateOrUpdate, err
		},
	}
}

func (r *Ruler) ruleConfigMaps() (retResources []resources.Resource) {

	prometheusRules, err := r.selectPrometheusRules()
	if err != nil {
		return errResourcesFunc(err)
	}

	// Although the shards attribute has a default value in the updated crd,
	// it may be still nil in previous ruler instances. So check it.
	var shards uint64 = 1
	if r.ruler.Spec.Shards != nil {
		shards = uint64(*r.ruler.Spec.Shards)
	}

	// generate rule files for each shard
	var shardsRuleFiles = make([]map[string]string, shards)
	for shardSn := range shardsRuleFiles {
		shardsRuleFiles[shardSn] = make(map[string]string)
	}
	if shards > 1 {
		for file, spec := range prometheusRules {
			if len(spec.Groups) == 0 {
				continue
			}
			var shardSpecs = make([]promv1.PrometheusRuleSpec, shards)
			for i := range spec.Groups {
				// hashmod to generate shard sequence number by file and group name
				name := fmt.Sprintf("%s/%s", file, spec.Groups[i].Name)
				shardSn := sum64(md5.Sum([]byte(name))) % shards
				shardSpecs[shardSn].Groups = append(shardSpecs[shardSn].Groups, spec.Groups[i])
			}
			for shardSn, shardSpec := range shardSpecs {
				if len(shardSpec.Groups) == 0 {
					continue
				}
				content, err := GenerateContent(shardSpec, r.Log)
				if err != nil {
					return errResourcesFunc(err)
				}
				shardsRuleFiles[shardSn][file] = content
			}
		}
	} else {
		for file, spec := range prometheusRules {
			if len(spec.Groups) == 0 {
				continue
			}
			content, err := GenerateContent(*spec, r.Log)
			if err != nil {
				return errResourcesFunc(err)
			}
			shardsRuleFiles[0][file] = content
		}
	}

	// generate configmaps
	var targets = make(map[string]*corev1.ConfigMap)
	r.shardsRuleConfigMapNames = make([]map[string]struct{}, shards)
	for shardSn := range shardsRuleFiles {
		ruleFiles := shardsRuleFiles[shardSn]
		cms, err := r.makeRulesConfigMaps(ruleFiles, shardSn)
		if err != nil {
			return errResourcesFunc(err)
		}
		for j := range cms {
			targets[cms[j].Name] = &cms[j]

			if r.shardsRuleConfigMapNames[shardSn] == nil {
				r.shardsRuleConfigMapNames[shardSn] = make(map[string]struct{})
			}
			r.shardsRuleConfigMapNames[shardSn][cms[j].Name] = struct{}{}
		}
	}

	var cmList corev1.ConfigMapList
	err = r.Client.List(r.Context, &cmList, client.InNamespace(r.ruler.Namespace))
	if err != nil {
		return errResourcesFunc(err)
	}
	// check configmaps to be deleted.
	// the configmaps owned by the ruler have a name
	// which concatenates a same name prefix, a shard sequence number and a configmap sequence number.
	currents := make(map[string]corev1.ConfigMap)
	namePrefix := r.name("rulefiles") + "-"
	for i := range cmList.Items {
		cm := cmList.Items[i]
		if !strings.HasPrefix(cm.Name, namePrefix) {
			continue
		}
		suffix := strings.TrimPrefix(cm.Name, namePrefix)
		sns := strings.Split(suffix, "-")
		if len(sns) != 2 {
			continue
		}
		shardSn, cmSn := sns[0], sns[1]
		if sequenceNumberRegexp.MatchString(shardSn) && sequenceNumberRegexp.MatchString(cmSn) {
			currents[cm.Name] = cmList.Items[i]
			if _, ok := targets[cm.Name]; !ok {
				retResources = append(retResources, func() (runtime.Object, resources.Operation, error) {
					return &cm, resources.OperationDelete, nil
				})
			}
		}
	}

	// create or update the targets if needed
	for name := range targets {
		target := targets[name]
		if current, ok := currents[name]; !ok || !reflect.DeepEqual(target.Data, current.Data) {
			retResources = append(retResources, func() (runtime.Object, resources.Operation, error) {
				return target, resources.OperationCreateOrUpdate, nil
			})
		}
	}

	return
}

// sum64 sums the md5 hash to an uint64.
func sum64(hash [md5.Size]byte) uint64 {
	var s uint64

	for i, b := range hash {
		shift := uint64((md5.Size - i - 1) * 8)

		s |= uint64(b) << shift
	}
	return s
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

func (r *Ruler) selectPrometheusRules() (map[string]*promv1.PrometheusRuleSpec, error) {
	rules := make(map[string]*promv1.PrometheusRuleSpec)

	namespaces, err := r.selectNamespaces(r.ruler.Spec.RuleNamespaceSelector)
	if err != nil {
		return nil, err
	}
	var ruleSelectors []labels.Selector
	for _, s := range r.ruler.Spec.RuleSelectors {
		ruleSelector, err := metav1.LabelSelectorAsSelector(s)
		if err != nil {
			return nil, err
		}
		ruleSelectors = append(ruleSelectors, ruleSelector)
	}

	for _, ns := range namespaces {
		var prometheusRules []*promv1.PrometheusRule
		for _, s := range ruleSelectors {
			var prometheusRuleList promv1.PrometheusRuleList
			err = r.Client.List(r.Context, &prometheusRuleList,
				client.MatchingLabelsSelector{Selector: s}, client.InNamespace(ns))
			if err != nil {
				return nil, err
			}
			prometheusRules = append(prometheusRules, prometheusRuleList.Items...)
		}
		for _, promRule := range prometheusRules {
			file := fmt.Sprintf("%v-%v.yaml", promRule.Namespace, promRule.Name)
			if _, ok := rules[file]; !ok {
				rules[file] = promRule.Spec.DeepCopy()
			}
		}
	}

	return rules, nil
}

// makeRulesConfigMaps refers to prometheus-operator and
// adds a shard sequence number argument to support ruler sharding.
//
// makeRulesConfigMaps takes a Ruler configuration and rule files and
// returns a list of Kubernetes ConfigMaps to be later on mounted into the
// Ruler instance.
// If the total size of rule files exceeds the Kubernetes ConfigMap limit,
// they are split up via the simple first-fit [1] bin packing algorithm. In the
// future this can be replaced by a more sophisticated algorithm, but for now
// simplicity should be sufficient.
// [1] https://en.wikipedia.org/wiki/Bin_packing_problem#First-fit_algorithm
func (r *Ruler) makeRulesConfigMaps(ruleFiles map[string]string, shardSn int) ([]corev1.ConfigMap, error) {
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
			ObjectMeta: r.meta(fmt.Sprintf("%s-%d-%d", r.name("rulefiles"), shardSn, i)),
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

	errs := operator.ValidateRule(promRule)
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

func (r *Ruler) envoyConfigMap(data map[string]string) error {
	var cm = &corev1.ConfigMap{ObjectMeta: r.meta(r.name("envoy-config"))}

	var buff strings.Builder
	tmpl := util.EnvoyStaticConfigTemplate
	if err := tmpl.Execute(&buff, data); err != nil {
		return err
	}

	cm.Data = map[string]string{
		envoyConfigFile: buff.String(),
	}

	if err := ctrl.SetControllerReference(r.ruler, cm, r.Scheme); err != nil {
		return err
	}
	_, err := controllerutil.CreateOrPatch(r.Context, r.Client, cm, configmapDataMutate(cm, cm.Data))
	return err
}

func configmapDataMutate(cm *corev1.ConfigMap, data map[string]string) controllerutil.MutateFn {
	return func() error {
		cm.Data = data
		return nil
	}
}
