package sreflect

import (
	"reflect"
	"runtime"
	"strings"
)

var FuncNamePrefixes []string

func FuncName(f any) string {
	func0 := runtime.FuncForPC(reflect.ValueOf(f).Pointer())

	name := strings.Split(func0.Name(), "-")[0]

	for _, prefix := range FuncNamePrefixes {
		if strings.HasPrefix(name, prefix) {
			return strings.TrimPrefix(name, prefix)
		}
	}

	return name
}
