package ssnap

import (
	"path/filepath"
	"time"
)

type DefaultWatcher struct {
	Snap *Snap
	Conf Config
}

func (w *DefaultWatcher) Save(exit bool) error {
	if exit {
		return w.Snap.Save()
	}

	d := w.Conf.Poll
	if d == 0 {
		return nil
	}

	if d < 30*time.Minute {
		d = 30 * time.Minute // 防止保存时间太短导致副本丢失
	}

	return w.Snap.Duplicate(d*3 + 10) // 至少保存 3 分副本
}

func (w *DefaultWatcher) Load() error {
	return w.Snap.Load()
}

func (w *DefaultWatcher) Config() Config {
	return w.Conf
}

type DefaultWatcherConfig struct {
	Path string
	Data any
	Name string
	Poll time.Duration
}

func DefaultWatchers(configs ...DefaultWatcherConfig) []Watcher {
	watchers := make([]Watcher, 0, len(configs))
	for _, config := range configs {
		if config.Name == "" {
			config.Name = filepath.Base(config.Path)
		}
		watchers = append(watchers, &DefaultWatcher{
			Snap: New(config.Path, config.Data),
			Conf: Config{
				Name: config.Name,
				Poll: config.Poll,
			},
		})
	}
	return watchers
}
