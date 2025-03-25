package sboot

import (
	"context"

	"github.com/yyliziqiu/slib/slog"
	"github.com/yyliziqiu/slib/sreflect"
)

type InitFunc func() error

type InitFuncs []InitFunc

func (list InitFuncs) Init() error {
	for _, fun := range list {
		slog.Infof("Init moudle: %s", sreflect.FuncName(fun))
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
		slog.Infof("Boot moudle: %s", sreflect.FuncName(fun))
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
