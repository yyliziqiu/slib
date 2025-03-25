package sreflect

import (
	"fmt"
	"reflect"
)

func FieldsOf(s any) []string {
	mt := reflect.TypeOf(s)
	var fields []string
	for i := 0; i < mt.NumField(); i++ {
		fields = append(fields, mt.Field(i).Name)
	}
	return fields
}

func ValuesOf(s any) []string {
	mv := reflect.ValueOf(s)
	var values []string
	for i := 0; i < mv.NumField(); i++ {
		values = append(values, fmt.Sprintf("%v", mv.Field(i).Interface()))
	}
	return values
}

func ValueOf(s any, fieldName string) (any, bool) {
	val := reflect.ValueOf(s)
	if val.Kind() == reflect.Pointer {
		val = val.Elem()
	}
	field := val.FieldByName(fieldName)
	if !field.IsValid() {
		return nil, false
	}
	return field.Interface(), true
}
