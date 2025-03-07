package sboot

import (
	"context"
	"reflect"
)

type InitFunc func() error

type InitFuncs []InitFunc

func (list InitFuncs) Init() error {
	for _, fun := range list {
		err := fun()
		if err != nil {
			return err
		}
	}
	return nil
}

type BootFunc func(context.Context) error

type BootFuncs []BootFunc

func (list BootFuncs) Boot(ctx context.Context) error {
	for _, fun := range list {
		err := fun(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

// Check 检查配置是否正确
type Check interface {
	Check() error
}

// Default 为配置设置默认值
type Default interface {
	Default()
}

func fieldValue(s any, name string) (any, bool) {
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
