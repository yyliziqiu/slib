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

func Init(conf Config) (err error) {
	_config = conf.Default()

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

func New(conf Config) (*logrus.Logger, error) {
	if conf.Console {
		return newConsoleLogger(conf)
	}
	return newFileLogger(conf)
}

func New2(name string) (*logrus.Logger, error) {
	conf := _config
	conf.Name = name
	return New(conf)
}

func New3(name string) *logrus.Logger {
	logger, err := New2(name)
	if err != nil {
		return Default
	}
	return logger
}

func newConsoleLogger(conf Config) (*logrus.Logger, error) {
	// 创建日志
	logger := logrus.New()

	// 设置日志等级
	logger.SetLevel(parseLevel(conf.Level))

	// 设置日志格式
	logger.SetFormatter(newFormatter(conf))

	// 禁止输出方法名
	logger.SetReportCaller(conf.ShowCaller)

	return logger, nil
}

func parseLevel(name string) logrus.Level {
	level, err := logrus.ParseLevel(name)
	if err != nil {
		return logrus.DebugLevel
	}
	return level
}

func newFormatter(conf Config) logrus.Formatter {
	var (
		dataFormat = conf.DataFormat
		dateFormat = conf.DateFormat
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

func newFileLogger(conf Config) (*logrus.Logger, error) {
	// 创建日志
	logger := logrus.New()

	// 禁止控制台输出
	logger.SetOutput(io.Discard)

	// 设置日志等级
	logger.SetLevel(parseLevel(conf.Level))

	// 禁止输出方法名
	logger.SetReportCaller(conf.ShowCaller)

	// 日志按天分割
	hook, err := newHook(conf)
	if err != nil {
		return nil, err
	}
	logger.AddHook(hook)

	return logger, nil
}

func newHook(conf Config) (*lfshook.LfsHook, error) {
	var (
		name = conf.Name
		path = conf.Path
		age  = conf.MaxAge
		rtt  = conf.RotateTime
		rtl  = conf.RotateLevel
		rtz  = conf.RotateTimezone
	)

	err := sfile.MakeDir(path)
	if err != nil {
		return nil, fmt.Errorf("make log dir failed [%v]", err)
	}

	if strings.HasSuffix(name, "-") {
		name = strings.TrimSuffix(name, "-")
	}

	var output any
	if rtl == 1 {
		output, err = newRotateLogs(path, name, age, rtt, rtz)
	} else {
		output, err = newWriterMap(path, name, age, rtt, rtl, rtz)
	}
	if err != nil {
		return nil, err
	}

	return lfshook.NewHook(output, newFormatter(conf)), nil
}

func newRotateLogs(dir string, name string, age time.Duration, rtt time.Duration, rtz string) (*rotate.RotateLogs, error) {
	path := filepath.Join(dir, name+"-%Y%m%d.log")

	loc, err := time.LoadLocation(rtz)
	if err != nil {
		loc = time.UTC
	}

	rls, err := rotate.New(path, rotate.WithMaxAge(age), rotate.WithRotationTime(rtt), rotate.WithLocation(loc))
	if err != nil {
		return nil, fmt.Errorf("new rotation failed [%v]", err)
	}

	return rls, nil
}

func newWriterMap(dir string, name string, age time.Duration, rtt time.Duration, rtl int, rtz string) (lfshook.WriterMap, error) {
	wm := make(lfshook.WriterMap)
	for prefix, levels := range levelDispatch(rtl, name) {
		rls, err := newRotateLogs(dir, prefix, age, rtt, rtz)
		if err != nil {
			return nil, err
		}
		for _, level := range levels {
			wm[level] = rls
		}
	}
	return wm, nil
}

func levelDispatch(rtl int, name string) (dispatch LevelDispatch) {
	var d, i, w, e, f, p logrus.Level = 5, 4, 3, 2, 1, 0

	switch rtl {
	case 3:
		dispatch = LevelDispatch{
			name:            {d, i},
			name + "-warn":  {w},
			name + "-error": {e, f, p},
		}
	case 4:
		dispatch = LevelDispatch{
			name + "-debug": {d},
			name + "-info":  {i},
			name + "-warn":  {w},
			name + "-error": {e, f, p},
		}
	case 5:
		dispatch = LevelDispatch{
			name:            {d, i},
			name + "-error": {w, e, f, p},
		}
	default:
		dispatch = LevelDispatch{
			name:            {d, i, w},
			name + "-error": {e, f, p},
		}
	}

	return dispatch
}
