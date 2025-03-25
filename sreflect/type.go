package sreflect

import (
	"reflect"
)

func FieldValue(s any, name string) (any, bool) {
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
