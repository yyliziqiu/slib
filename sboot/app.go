package sboot

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yyliziqiu/slib/sconfig"
	"github.com/yyliziqiu/slib/slog"
	"github.com/yyliziqiu/slib/sreflect"
)

type App struct {
	// app 名称
	Name string

	// app 版本
	Version string

	// 默认配置文件路径，可被命令行参数覆盖
	ConfigPath string

	// 全局配置
	ConfigRoot any

	// 模块
	InitFuncs func() InitFuncs
	BootFuncs func() BootFuncs

	hasInitConfig bool
	hasInitModule bool
}

// Init app
func (app *App) Init() (err error) {
	err = app.InitConfig()
	if err != nil {
		return err
	}

	return app.InitModule()
}

func (app *App) InitConfig() (err error) {
	if app.hasInitConfig {
		return nil
	}
	app.hasInitConfig = true

	log.Println("Init config.")

	// 加载配置文件
	err = sconfig.Init(app.ConfigPath, app.ConfigRoot)
	if err != nil {
		return fmt.Errorf("init config failed [%v]", err)
	}

	// 检查配置是否正确
	check, ok := app.ConfigRoot.(Check)
	if ok {
		err = check.Check()
		if err != nil {
			return err
		}
	}

	// 为配置设置默认值
	default0, ok := app.ConfigRoot.(Default)
	if ok {
		default0.Default()
	}

	log.Println("Init log.")

	// 初始化日志
	logConfig := slog.Config{Console: true}
	if logConfig1, ok1 := sreflect.ValueOf(app.ConfigRoot, "Log"); ok1 {
		logConfig2, ok2 := logConfig1.(slog.Config)
		if ok2 {
			logConfig = logConfig2
		}
	}
	err = slog.Init(logConfig)
	if err != nil {
		return fmt.Errorf("init log failed [%v]", err)
	}

	return nil
}

func (app *App) InitModule() (err error) {
	if app.hasInitModule {
		return nil
	}
	app.hasInitModule = true

	err = app.InitFuncs().Init()
	if err != nil {
		slog.Errorf("Init modules failed, error: %v", err)
		return err
	}

	return nil
}

// Run app
func (app *App) Run() (err error) {
	err = app.Init()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())

	err = app.BootFuncs().Boot(ctx)
	if err != nil {
		slog.Errorf("Boot modules failed, error: %v", err)
		cancel()
		return err
	}

	slog.Info("App boot successfully.")

	exitCh := make(chan os.Signal)
	signal.Notify(exitCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-exitCh

	cancel()

	if exitWait, ok := sreflect.ValueOf(app.ConfigRoot, "ExitWait"); ok {
		exitWait2, ok2 := exitWait.(time.Duration)
		if ok2 && exitWait2 > 0 {
			time.Sleep(exitWait2)
		}
	}

	slog.Info("App exit.")

	return nil
}
