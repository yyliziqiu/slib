package squeue

import (
	"path/filepath"
	"time"

	"github.com/yyliziqiu/slib/ssnap"
)

type SnapWatcher struct {
	Queue *Queue
	Item  any
	Conf  ssnap.WatchConfig
}

func (w *SnapWatcher) Save(exit bool) error {
	if exit {
		return w.Queue.Save()
	}

	d := w.Conf.Poll
	if d == 0 {
		return nil
	}

	return w.Queue.SnapDuplicate(d*3 + time.Second)
}

func (w *SnapWatcher) Load() error {
	return w.Queue.Load(w.Item)
}

func (w *SnapWatcher) WatchConfig() ssnap.WatchConfig {
	return w.Conf
}

type SnapWatcherConfig struct {
	Queue *Queue
	Item  any
	Path  string
	Name  string
	Poll  time.Duration
}

func SnapWatchers(configs ...SnapWatcherConfig) []ssnap.Watcher {
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
		watchers = append(watchers, &SnapWatcher{
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
