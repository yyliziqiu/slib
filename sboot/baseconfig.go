package sboot

import (
	"path/filepath"
	"time"

	"github.com/yyliziqiu/slib/sdb"
	"github.com/yyliziqiu/slib/senv"
	"github.com/yyliziqiu/slib/ses"
	"github.com/yyliziqiu/slib/skafka"
	"github.com/yyliziqiu/slib/skvs"
	"github.com/yyliziqiu/slib/slog"
	"github.com/yyliziqiu/slib/sredis"
	"github.com/yyliziqiu/slib/sserver"
	"github.com/yyliziqiu/slib/stask"
)

type Config struct {
	Env      string
	AppId    string
	InsId    string
	BasePath string
	DataPath string
	ExitWait time.Duration

	Server sserver.Config
	Log    slog.Config

	Db        sdb.Config
	Es        ses.Config
	Redis     sredis.Config
	Kafka     skafka.Config
	DbList    []sdb.Config
	EsList    []ses.Config
	RedisList []sredis.Config
	KafkaList []skafka.Config

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
