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

	Console, err = newConsole(_config)
	if err != nil {
		return err
	}

	return nil
}

func New(config Config) (*logrus.Logger, error) {
	if config.Console {
		return newConsole(config)
	}
	return newFile(config)
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

func newConsole(config Config) (*logrus.Logger, error) {
	logger := logrus.New()

	// 设置日志等级
	logger.SetLevel(level(config.Level))

	// 设置日志格式
	logger.SetFormatter(formatter(config))

	// 禁止输出方法名
	logger.SetReportCaller(config.ShowCaller)

	return logger, nil
}

func level(name string) logrus.Level {
	lvl, err := logrus.ParseLevel(name)
	if err != nil {
		return logrus.DebugLevel
	}
	return lvl
}

func formatter(config Config) logrus.Formatter {
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

func newFile(config Config) (*logrus.Logger, error) {
	logger := logrus.New()

	// 禁止控制台输出
	logger.SetOutput(io.Discard)

	// 设置日志等级
	logger.SetLevel(level(config.Level))

	// 禁止输出方法名
	logger.SetReportCaller(config.ShowCaller)

	// 日志按天分割
	var err error
	var hook *lfshook.LfsHook
	switch config.RotationLevel {
	case 0:
		hook, err = rotatesHook0(config)
	case 1:
		hook, err = rotatesHook1(config)
	default:
		hook, err = rotatesHook2(config)
	}
	if err != nil {
		return nil, fmt.Errorf("create hook failed [%v]", err)
	}
	logger.AddHook(hook)

	return logger, nil
}

func rotatesHook0(config Config) (*lfshook.LfsHook, error) {
	var (
		name         = config.Name
		path         = config.Path
		maxAge       = config.MaxAge
		rotationTime = config.RotationTime
	)

	// 确保日志目录存在
	err := sfile.MakeDir(config.Path)
	if err != nil {
		return nil, fmt.Errorf("create logs dir failed [%v]", err)
	}

	// 美化日志文件名
	if !strings.HasSuffix(name, "-") {
		name = name + "-"
	}

	// 创建分割器
	rotation, err := rotates(path, name+"%Y%m%d.log", maxAge, rotationTime)
	if err != nil {
		return nil, fmt.Errorf("create rotate failed [%v]", err)
	}
	errorRotation, err := rotates(path, name+"error-%Y%m%d.log", maxAge, rotationTime)
	if err != nil {
		return nil, fmt.Errorf("create rotate failed [%v]", err)
	}

	return lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: rotation,
		logrus.InfoLevel:  rotation,
		logrus.WarnLevel:  rotation,
		logrus.ErrorLevel: errorRotation,
		logrus.FatalLevel: errorRotation,
		logrus.PanicLevel: errorRotation,
	}, formatter(config)), nil
}

func rotatesHook1(config Config) (*lfshook.LfsHook, error) {
	var (
		name         = config.Name
		path         = config.Path
		maxAge       = config.MaxAge
		rotationTime = config.RotationTime
	)

	// 确保日志目录存在
	err := sfile.MakeDir(path)
	if err != nil {
		return nil, fmt.Errorf("create logs dir failed [%v]", err)
	}

	// 美化日志文件名
	if !strings.HasSuffix(name, "-") {
		name = name + "-"
	}

	// 创建分割器
	rotation, err := rotates(path, name+"%Y%m%d.log", maxAge, rotationTime)
	if err != nil {
		return nil, fmt.Errorf("create rotate failed [%v]", err)
	}

	return lfshook.NewHook(rotation, formatter(config)), nil
}

func rotatesHook2(config Config) (*lfshook.LfsHook, error) {
	var (
		name         = config.Name
		path         = config.Path
		maxAge       = config.MaxAge
		rotationTime = config.RotationTime
	)

	// 确保日志目录存在
	err := sfile.MakeDir(config.Path)
	if err != nil {
		return nil, fmt.Errorf("create logs dir failed [%v]", err)
	}

	// 美化日志文件名
	if !strings.HasSuffix(name, "-") {
		name = name + "-"
	}

	// 创建分割器
	debugRotation, err := rotates(path, name+"debug-%Y%m%d.log", maxAge, rotationTime)
	if err != nil {
		return nil, fmt.Errorf("create rotate failed [%v]", err)
	}
	infoRotation, err := rotates(path, name+"info-%Y%m%d.log", maxAge, rotationTime)
	if err != nil {
		return nil, fmt.Errorf("create rotate failed [%v]", err)
	}
	warnRotation, err := rotates(path, name+"warn-%Y%m%d.log", maxAge, rotationTime)
	if err != nil {
		return nil, fmt.Errorf("create rotate failed [%v]", err)
	}
	errorRotation, err := rotates(path, name+"error-%Y%m%d.log", maxAge, rotationTime)
	if err != nil {
		return nil, fmt.Errorf("create rotate failed [%v]", err)
	}

	return lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: debugRotation,
		logrus.InfoLevel:  infoRotation,
		logrus.WarnLevel:  warnRotation,
		logrus.ErrorLevel: errorRotation,
		logrus.FatalLevel: errorRotation,
		logrus.PanicLevel: errorRotation,
	}, formatter(config)), nil
}

func rotates(dirname string, filename string, maxAge time.Duration, RotationTime time.Duration) (*rotate.RotateLogs, error) {
	return rotate.New(filepath.Join(dirname, filename), rotate.WithMaxAge(maxAge), rotate.WithRotationTime(RotationTime))
}
