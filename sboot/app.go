package sboot

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yyliziqiu/slib/sconfig"
	"github.com/yyliziqiu/slib/slog"
)

type App struct {
	// app 名称
	Name string

	// app 版本
	Version string

	// 配置文件路径
	ConfigPath string

	// 全局配置
	ConfigRoot any

	// 模块
	InitFuncs   InitFuncs
	BootFuncs   BootFuncs
	InitFuncsCb func() InitFuncs
	BootFuncsCb func() BootFuncs

	hasCallInitFuncs bool
}

// Init app
func (app *App) Init() (err error) {
	err = app.InitConfig()
	if err != nil {
		return err
	}
	return app.CallInitFuncs()
}

func (app *App) InitConfig() (err error) {
	// 加载配置文件
	err = sconfig.Init(app.ConfigPath, app.ConfigRoot)
	if err != nil {
		return fmt.Errorf("init config error [%v]", err)
	}

	// 检查配置是否正确
	icheck, ok := app.ConfigRoot.(ICheck)
	if ok {
		err = icheck.Check()
		if err != nil {
			return err
		}
	}

	// 为配置项设置默认值
	idefault, ok := app.ConfigRoot.(IDefault)
	if ok {
		idefault.Default()
	}

	// 初始化日志
	logc := slog.Config{Console: true}
	ilog, ok := app.ConfigRoot.(IGetLog)
	if ok {
		logc = ilog.GetLog()
	}
	err = slog.Init(logc)
	if err != nil {
		return fmt.Errorf("init log error [%v]", err)
	}

	return nil
}

func (app *App) CallInitFuncs() (err error) {
	if app.hasCallInitFuncs {
		return nil
	}
	app.hasCallInitFuncs = true

	initFuncs := app.InitFuncs
	if app.InitFuncsCb != nil {
		initFuncs = app.InitFuncsCb()
	}

	slog.Info("Prepare init funcs.")
	err = initFuncs.Init()
	if err != nil {
		slog.Errorf("Init funcs failed, error: %v", err)
		return err
	}
	slog.Info("Init funcs succeed.")

	return nil
}

// Start app
func (app *App) Start() (err error, f context.CancelFunc) {
	err = app.InitConfig()
	if err != nil {
		return err, nil
	}
	return app.CallBootFuncs()
}

func (app *App) CallBootFuncs() (error, context.CancelFunc) {
	err := app.CallInitFuncs()
	if err != nil {
		return err, nil
	}

	ctx, cancel := context.WithCancel(context.Background())

	bootFuncs := app.BootFuncs
	if app.BootFuncsCb != nil {
		bootFuncs = app.BootFuncsCb()
	}

	slog.Info("Prepare boot funcs.")
	err = bootFuncs.Boot(ctx)
	if err != nil {
		slog.Errorf("Boot funcs failed, error: %v", err)
		cancel()
		return err, nil
	}
	slog.Info("Boot funcs successfully.")

	return nil, cancel
}

// Run app
func (app *App) Run() (err error) {
	err = app.InitConfig()
	if err != nil {
		return err
	}
	return app.CallBootFuncsBlocked()
}

func (app *App) CallBootFuncsBlocked() (err error) {
	err, cancel := app.CallBootFuncs()
	if err != nil {
		return err
	}

	slog.Info("App run successfully.")

	exitCh := make(chan os.Signal)
	signal.Notify(exitCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-exitCh

	slog.Info("App prepare exit.")

	cancel()

	iwaittime, ok := app.ConfigRoot.(IGetWaitTime)
	if ok {
		time.Sleep(iwaittime.GetWaitTime())
	}

	slog.Info("App exit.")

	return nil
}
