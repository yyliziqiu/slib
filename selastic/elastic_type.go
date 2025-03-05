package selastic

import (
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"

	"github.com/yyliziqiu/slib/slog"
)

const DefaultId = "default"

type Config struct {
	Id           string   // optional
	Hosts        []string // must
	Username     string   // must
	Password     string   // must
	EnableLogger bool     // optional

	Logger *logrus.Logger `json:"-"` // optional
	Client elastic.Doer   `json:"-"` // optional
}

func (c Config) Default() Config {
	if c.Id == "" {
		c.Id = DefaultId
	}
	if c.EnableLogger && c.Logger == nil {
		c.Logger = slog.Default
	}
	return c
}
