package ssnap

import (
	"path/filepath"
	"sync"
	"time"
)

type DefaultWatcher struct {
	Snap *Snap
	Conf Config
	Mu   sync.Locker
}

func (w *DefaultWatcher) Save(exit bool) error {
	w.Mu.Lock()
	defer w.Mu.Unlock()

	if exit {
		return w.Snap.Save()
	}

	d := w.Conf.Poll
	if d == 0 {
		return nil
	}

	if d < 10*time.Minute {
		d = 10 * time.Minute // 防止保存时间太短导致副本丢失
	}

	return w.Snap.Duplicate(d*3 + 10*time.Second) // 至少保存 3 分副本
}

func (w *DefaultWatcher) Load() error {
	w.Mu.Lock()
	defer w.Mu.Unlock()

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
	Mu   sync.Locker
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
			Mu: config.Mu,
		})
	}
	return watchers
}
