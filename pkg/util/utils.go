package util

import (
	"reflect"
	"strings"

	"gopkg.in/yaml.v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/kubesphere/whizard/pkg/constants"
)

func AppendLabel(m1, m2 map[string]string) {
	if m1 == nil {
		m1 = make(map[string]string)
	}

	for k, v := range m2 {
		if _, ok := m1[k]; !ok {
			m1[k] = v
		}
	}
}

func Contains(list []string, key string) bool {
	for _, k := range list {
		if k == key {
			return true
		}
	}

	return false
}

func YamlMarshal(val interface{}) (string, error) {

	out, err := yaml.Marshal(val)
	if err != nil {
		return "", err
	}

	return string(out), nil
}

func Join(sep string, elem ...string) string {
	return strings.Join(elem, sep)
}

func MegerMap(src, dest map[string]string) {

	for k, v := range src {
		dest[k] = v
	}
}

func ReplaceInSlice(s interface{}, fn func(v interface{}) bool, new interface{}) bool {
	v := reflect.ValueOf(s)
	if v.Kind() != reflect.Slice || v.IsNil() {
		return false
	}

	for i := 0; i < v.Len(); i++ {
		elem := v.Index(i)
		if fn(elem.Interface()) {
			elem.Set(reflect.ValueOf(new))
			return true
		}
	}

	return false
}

func GetArgName(s string) string {

	keys := strings.Split(s, "=")
	return keys[0]
}

func ManagedLabelByService(service metav1.Object) map[string]string {
	return map[string]string{
		constants.ServiceLabelKey: service.GetNamespace() + "." + service.GetName(),
	}
}

func ManagedLabelBySameService(o metav1.Object) map[string]string {
	return map[string]string{
		constants.ServiceLabelKey: o.GetLabels()[constants.ServiceLabelKey],
	}
}

func ServiceNamespacedName(managedByService metav1.Object) *types.NamespacedName {
	ls := managedByService.GetLabels()
	if len(ls) == 0 {
		return nil
	}

	namespacedName := ls[constants.ServiceLabelKey]
	arr := strings.Split(namespacedName, ".")
	if len(arr) != 2 {
		return nil
	}

	return &types.NamespacedName{
		Namespace: arr[0],
		Name:      arr[1],
	}
}

func ManagedLabelByStorage(storage metav1.Object) map[string]string {
	return map[string]string{
		constants.StorageLabelKey: storage.GetNamespace() + "." + storage.GetName(),
	}
}

func StorageNamespacedName(managedByStorage metav1.Object) *types.NamespacedName {
	ls := managedByStorage.GetLabels()
	if len(ls) == 0 {
		return nil
	}

	namespacedName := ls[constants.StorageLabelKey]
	arr := strings.Split(namespacedName, ".")
	if len(arr) != 2 {
		return nil
	}

	return &types.NamespacedName{
		Namespace: arr[0],
		Name:      arr[1],
	}
}
