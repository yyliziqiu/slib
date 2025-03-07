package sboot

import (
	"context"
	"time"

	"github.com/yyliziqiu/slib/slog"
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

// LogConfig 获取日志配置
type LogConfig interface {
	LogConfig() slog.Config
}

// ExitWait 获取应用退出时的等待时长
type ExitWait interface {
	ExitWait() time.Duration
}
