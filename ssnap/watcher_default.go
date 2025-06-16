package ssnap

import (
	"path/filepath"
	"time"
)

type DefaultWatcher struct {
	snap *Snap
	conf Config
}

func (w *DefaultWatcher) Save(exit bool) error {
	if exit {
		return w.snap.Save()
	}

	d := w.conf.Poll
	if d == 0 {
		return nil
	}

	if d < 30*time.Minute {
		d = 30 * time.Minute // 防止保存时间太短导致副本丢失
	}

	return w.snap.Duplicate(d*3 + 10) // 至少保存 3 分副本
}

func (w *DefaultWatcher) Load() error {
	return w.snap.Load()
}

func (w *DefaultWatcher) Config() Config {
	return w.conf
}

type Setting struct {
	Path string
	Data any
	Name string
	Poll time.Duration
}

func DefaultWatchers(settings ...Setting) []Watcher {
	watchers := make([]Watcher, 0, len(settings))
	for _, setting := range settings {
		if setting.Name == "" {
			setting.Name = filepath.Base(setting.Path)
		}
		watchers = append(watchers, &DefaultWatcher{
			snap: New(setting.Path, setting.Data),
			conf: Config{
				Name: setting.Name,
				Poll: setting.Poll,
			},
		})
	}
	return watchers
}
