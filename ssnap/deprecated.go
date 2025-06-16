package ssnap

import (
	"context"
	"reflect"
	"time"

	"github.com/yyliziqiu/slib/slog"
	"github.com/yyliziqiu/slib/stime"
)

// Handler deprecated
type Handler interface {
	Load() error
	Save() error
}

// HandlerName deprecated
type HandlerName interface {
	Name() string
}

// HandlerInterval deprecated
type HandlerInterval interface {
	Interval() time.Duration
}

// Watch deprecated
func Watch(ctx context.Context, handlers []Handler) error {
	err := watchLoad(handlers)
	if err != nil {
		return err
	}

	for _, handler := range handlers {
		go runWatchSave(ctx, handler)
	}

	return nil
}

func watchLoad(handlers []Handler) error {
	timer := stime.NewTimer()

	for _, handler := range handlers {
		err := handler.Load()
		if err != nil {
			slog.Errorf("Load snap failed, name: %s, error: %v.", watchName(handler), err)
			return err
		}
		slog.Infof("Load snap succeed, name: %s, cost: %s.", watchName(handler), timer.Pauses())
	}

	slog.Infof("Load snaps compeleted, cost: %s.", timer.Stops())

	return nil
}

func watchName(handler Handler) string {
	i, ok := handler.(HandlerName)
	if ok {
		return i.Name()
	}

	typ := reflect.TypeOf(handler)
	if typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}

	return typ.Name()
}

func runWatchSave(ctx context.Context, handler Handler) {
	interval := watchInterval(handler)

	if interval <= 0 {
		<-ctx.Done()
		_ = watchSave(handler)
		return
	}

	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ticker.C:
			_ = watchSave(handler)
		case <-ctx.Done():
			_ = watchSave(handler)
			return
		}
	}
}

func watchInterval(handler Handler) time.Duration {
	i, ok := handler.(HandlerInterval)
	if ok {
		return i.Interval()
	}
	return 0
}

func watchSave(handler Handler) error {
	timer := stime.NewTimer()

	err := handler.Save()
	if err != nil {
		slog.Errorf("Save snap failed, name: %s, error: %v.", watchName(handler), err)
	} else {
		slog.Infof("Save snap succeed, name: %s, cost: %s.", watchName(handler), timer.Stops())
	}

	return err
}
