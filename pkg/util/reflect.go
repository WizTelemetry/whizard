package util

import "reflect"

func Override(left interface{}, right interface{}) {
	if reflect.ValueOf(left).IsNil() || reflect.ValueOf(right).IsNil() {
		return
	}

	if reflect.ValueOf(left).Type().Kind() != reflect.Ptr ||
		reflect.ValueOf(right).Type().Kind() != reflect.Ptr ||
		reflect.ValueOf(left).Kind() != reflect.ValueOf(right).Kind() {
		return
	}

	oldVal := reflect.ValueOf(left).Elem()
	newVal := reflect.ValueOf(right).Elem()

	for i := 0; i < oldVal.NumField(); i++ {
		val := newVal.Field(i).Interface()
		if !reflect.DeepEqual(val, reflect.Zero(reflect.TypeOf(val)).Interface()) {
			oldVal.Field(i).Set(reflect.ValueOf(val))
		}
	}
}
