package scsv

import (
	"fmt"
	"reflect"
)

func structFields(model any) []string {
	mt := reflect.TypeOf(model)
	var fields []string
	for i := 0; i < mt.NumField(); i++ {
		fields = append(fields, mt.Field(i).Name)
	}
	return fields
}

func structValues(model any) []string {
	mv := reflect.ValueOf(model)
	var values []string
	for i := 0; i < mv.NumField(); i++ {
		values = append(values, fmt.Sprintf("%v", mv.Field(i).Interface()))
	}
	return values
}

func structValue(s any, name string) (any, bool) {
	val := reflect.ValueOf(s)
	if val.Kind() == reflect.Pointer {
		val = val.Elem()
	}
	field := val.FieldByName(name)
	if !field.IsValid() {
		return nil, false
	}
	return field.Interface(), true
}
