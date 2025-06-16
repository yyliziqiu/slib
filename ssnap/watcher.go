package ssnap

import (
	"context"
	"reflect"
	"time"

	"github.com/yyliziqiu/slib/slog"
	"github.com/yyliziqiu/slib/stime"
)

type Watcher interface {
	Load() error
	Save(exit bool) error
}

type Config struct {
	Name string
	Poll time.Duration
}

type WatcherConfig interface {
	Config() Config
}

func Watches(ctx context.Context, watchers []Watcher) error {
	// load
	err := watcherLoad(watchers)
	if err != nil {
		return err
	}

	// save
	for _, watcher := range watchers {
		go runWatcherSave(ctx, watcher)
	}

	return nil
}

func watcherLoad(watchers []Watcher) error {
	timer := stime.NewTimer()

	for _, watcher := range watchers {
		err := watcher.Load()
		if err != nil {
			slog.Errorf("Load snap failed, name: %s, error: %v.", watcherName(watcher), err)
			return err
		}
		slog.Infof("Load snap succeed, name: %s, cost: %s.", watcherName(watcher), timer.Pauses())
	}

	slog.Infof("Load all snaps compeleted, cost: %s.", timer.Stops())

	return nil
}

func watcherName(watcher Watcher) string {
	conf := watcherConfig(watcher)
	if conf.Name != "" {
		return conf.Name
	}

	typ := reflect.TypeOf(watcher)
	if typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}

	return typ.Name()
}

func watcherConfig(watcher any) Config {
	c, ok := watcher.(WatcherConfig)
	if ok {
		return c.Config()
	}
	return Config{}
}

func runWatcherSave(ctx context.Context, watcher Watcher) {
	conf := watcherConfig(watcher)

	if conf.Poll <= 0 {
		<-ctx.Done()
		_ = watcherSave(watcher, true)
		return
	}

	ticker := time.NewTicker(conf.Poll)
	for {
		select {
		case <-ticker.C:
			_ = watcherSave(watcher, false)
		case <-ctx.Done():
			_ = watcherSave(watcher, true)
			return
		}
	}
}

func watcherSave(watcher Watcher, exit bool) error {
	timer := stime.NewTimer()

	err := watcher.Save(exit)
	if err != nil {
		slog.Errorf("Save snap failed, name: %s, error: %v.", watcherName(watcher), err)
	} else {
		slog.Infof("Save snap succeed, name: %s, cost: %s.", watcherName(watcher), timer.Stops())
	}

	return err
}
