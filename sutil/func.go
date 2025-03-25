package sutil

import (
	"reflect"
	"runtime"
	"strings"
)

func FuncName(f any) string {
	name := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	return strings.Split(name, "-")[0]
}
