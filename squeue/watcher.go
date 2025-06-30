package squeue

import (
	"path/filepath"
	"time"

	"github.com/yyliziqiu/slib/ssnap"
)

type Watcher struct {
	Queue *Queue
	Item  any
	Conf  ssnap.WatchConfig
}

func (w *Watcher) Save(exit bool) error {
	if exit {
		return w.Queue.Save()
	}

	d := w.Conf.Poll
	if d == 0 {
		return nil
	}

	return w.Queue.Duplicate(d*3 + 10*time.Second)
}

func (w *Watcher) Load() error {
	return w.Queue.Load(w.Item)
}

func (w *Watcher) WatchConfig() ssnap.WatchConfig {
	return w.Conf
}

type WatcherConfig struct {
	Queue *Queue
	Item  any
	Path  string
	Name  string
	Poll  time.Duration
}

func Watchers(configs ...WatcherConfig) []ssnap.Watcher {
	watchers := make([]ssnap.Watcher, 0, len(configs))
	for _, config := range configs {
		if config.Path != "" {
			config.Queue.path = config.Path
		}
		if config.Queue.path == "" {
			continue
		}
		if config.Name == "" {
			config.Name = filepath.Base(config.Queue.path)
		}
		watchers = append(watchers, &Watcher{
			Queue: config.Queue,
			Item:  config.Item,
			Conf: ssnap.WatchConfig{
				Name: config.Name,
				Poll: config.Poll,
			},
		})
	}
	return watchers
}
