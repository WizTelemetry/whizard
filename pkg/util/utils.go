package util

import (
	"reflect"
	"strings"

	"gopkg.in/yaml.v3"
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
