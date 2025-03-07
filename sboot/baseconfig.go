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

type Config struct {
	Env      string
	AppId    string
	InsId    string
	BasePath string
	DataPath string
	ExitWait time.Duration

	Log slog.Config
	Web sweb.Config

	Db       sdb.Config
	Redis    sredis.Config
	Kafka    skafka.Config
	Elastic  selastic.Config
	Dbs      []sdb.Config
	Redises  []sredis.Config
	Kafkas   []skafka.Config
	Elastics []selastic.Config

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
