package ssnap

import (
	"path/filepath"
	"sync"
	"time"
)

type DefaultWatcher struct {
	Snap *Snap
	Conf WatchConfig
	Mu   sync.Locker
}

func (w *DefaultWatcher) WatchConfig() WatchConfig {
	return w.Conf
}

func (w *DefaultWatcher) Save(exit bool) error {
	w.Mu.Lock()
	defer w.Mu.Unlock()

	if exit {
		return w.Snap.Save()
	}

	return w.Snap.Duplicate()
}

func (w *DefaultWatcher) Load() error {
	w.Mu.Lock()
	defer w.Mu.Unlock()

	return w.Snap.Load()
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
			Snap: New3(config.Path, config.Data, config.Poll, 3),
			Conf: WatchConfig{
				Name: config.Name,
				Poll: config.Poll,
			},
			Mu: config.Mu,
		})
	}
	return watchers
}
