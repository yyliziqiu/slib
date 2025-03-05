package sboot

import (
	"path/filepath"
	"time"

	"github.com/yyliziqiu/slib/sdb"
	"github.com/yyliziqiu/slib/selastic"
	"github.com/yyliziqiu/slib/senv"
	"github.com/yyliziqiu/slib/skafka"
	"github.com/yyliziqiu/slib/skvs"
	"github.com/yyliziqiu/slib/slog"
	"github.com/yyliziqiu/slib/sredis"
	"github.com/yyliziqiu/slib/stask"
	"github.com/yyliziqiu/slib/sweb"
)

// ICheck 检查配置是否正确
type ICheck interface {
	Check() error
}

// IDefault 为配置项设置默认值
type IDefault interface {
	Default()
}

// IGetLog 获取日志配置
type IGetLog interface {
	GetLog() slog.Config
}

// IGetWaitTime 获取应用退出时等待时长配置
type IGetWaitTime interface {
	GetWaitTime() time.Duration
}

type Config struct {
	Env      string
	AppId    string
	InsId    string
	BasePath string
	DataPath string
	WaitTime time.Duration

	Log slog.Config
	Web sweb.Config

	Db      []sdb.Config
	Redis   []sredis.Config
	Elastic []selastic.Config
	Kafka   []skafka.Config

	CronTask []stask.CronTask
	OnceTask []stask.OnceTask

	Values skvs.Kvs
}

func (c *Config) Default() {
	if c.Env == "" {
		c.Env = senv.Prod
	}
	if c.AppId == "" {
		c.AppId = "app"
	}
	if c.InsId == "" {
		c.InsId = "1"
	}
	if c.BasePath == "" {
		c.BasePath = "."
	}
	if c.DataPath == "" {
		c.DataPath = filepath.Join(c.BasePath, "data")
	}
}

func (c *Config) GetLog() slog.Config {
	return c.Log
}

func (c *Config) GetWaitTime() time.Duration {
	return c.WaitTime
}
