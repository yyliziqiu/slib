package slog

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	rotate "github.com/lestrrat/go-file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"

	"github.com/yyliziqiu/slib/sfile"
)

var (
	_config Config

	Default *logrus.Logger
	Console *logrus.Logger
)

func Init(config Config) (err error) {
	_config = config.Default()

	Default, err = New(_config)
	if err != nil {
		return err
	}

	Console, err = newConsoleLogger(_config)
	if err != nil {
		return err
	}

	return nil
}

func New(config Config) (*logrus.Logger, error) {
	if config.Console {
		return newConsoleLogger(config)
	}
	return newFileLogger(config)
}

func New2(name string) (*logrus.Logger, error) {
	config := _config
	config.Name = name
	return New(config)
}

func New3(name string) *logrus.Logger {
	logger, err := New2(name)
	if err != nil {
		return Default
	}
	return logger
}

func newConsoleLogger(config Config) (*logrus.Logger, error) {
	// 创建日志
	logger := logrus.New()

	// 设置日志等级
	logger.SetLevel(parseLevel(config.Level))

	// 设置日志格式
	logger.SetFormatter(newFormatter(config))

	// 禁止输出方法名
	logger.SetReportCaller(config.ShowCaller)

	return logger, nil
}

func parseLevel(name string) logrus.Level {
	level, err := logrus.ParseLevel(name)
	if err != nil {
		return logrus.DebugLevel
	}
	return level
}

func newFormatter(config Config) logrus.Formatter {
	var (
		dataFormat = config.DataFormat
		dateFormat = config.DateFormat
	)

	if dateFormat == "" {
		dateFormat = time.DateTime
	}

	switch dataFormat {
	case JsonFormat:
		return &logrus.JSONFormatter{
			TimestampFormat: dateFormat,
		}
	default:
		return &logrus.TextFormatter{
			DisableQuote:    true,
			TimestampFormat: dateFormat,
		}
	}
}

func newFileLogger(config Config) (*logrus.Logger, error) {
	// 创建日志
	logger := logrus.New()

	// 禁止控制台输出
	logger.SetOutput(io.Discard)

	// 设置日志等级
	logger.SetLevel(parseLevel(config.Level))

	// 禁止输出方法名
	logger.SetReportCaller(config.ShowCaller)

	// 日志按天分割
	hook, err := newHook(config)
	if err != nil {
		return nil, err
	}
	logger.AddHook(hook)

	return logger, nil
}

func newHook(config Config) (*lfshook.LfsHook, error) {
	var (
		name = config.Name
		path = config.Path
		age  = config.MaxAge
		rtt  = config.RotationTime
		rtl  = config.RotationLevel
	)

	// 确保日志目录存在
	err := sfile.MakeDir(path)
	if err != nil {
		return nil, fmt.Errorf("make log dir failed [%v]", err)
	}

	// 美化日志文件名
	if !strings.HasSuffix(name, "-") {
		name = name + "-"
	}

	// 生成 output
	var output any
	if rtl == 1 {
		output, err = newRotation(path, name, age, rtt)
	} else {
		output, err = newOutput(newDispatch(name, rtl), path, age, rtt)
	}
	if err != nil {
		return nil, err
	}

	return lfshook.NewHook(output, newFormatter(config)), nil
}

func newRotation(name string, dir string, age time.Duration, rtt time.Duration) (*rotate.RotateLogs, error) {
	rt, err := rotate.New(filepath.Join(dir, name+"-%Y%m%d.log"), rotate.WithMaxAge(age), rotate.WithRotationTime(rtt))
	if err != nil {
		return nil, fmt.Errorf("new rotation failed [%v]", err)
	}
	return rt, nil
}

func newOutput(dispatch Dispatch, dir string, age time.Duration, rtt time.Duration) (lfshook.WriterMap, error) {
	output := lfshook.WriterMap{}
	for name, levels := range dispatch {
		rt, err := newRotation(name, dir, age, rtt)
		if err != nil {
			return nil, err
		}
		for _, level := range levels {
			output[level] = rt
		}
	}
	return output, nil
}

func newDispatch(name string, rtl int) Dispatch {
	var (
		dispatch         Dispatch
		d, i, w, e, f, p logrus.Level = 5, 4, 3, 2, 1, 0
	)

	switch rtl {
	case 3:
		dispatch = Dispatch{
			name:           {d, i},
			name + "warn":  {w},
			name + "error": {e, f, p},
		}
	case 4:
		dispatch = Dispatch{
			name + "debug": {d},
			name + "info":  {i},
			name + "warn":  {w},
			name + "error": {e, f, p},
		}
	case 5:
		dispatch = Dispatch{
			name:           {d, i},
			name + "error": {w, e, f, p},
		}
	default:
		dispatch = Dispatch{
			name:           {d, i, w},
			name + "error": {e, f, p},
		}
	}

	return dispatch
}
