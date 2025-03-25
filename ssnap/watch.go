package ssnap

import (
	"context"
	"reflect"
	"time"

	"github.com/yyliziqiu/slib/slog"
	"github.com/yyliziqiu/slib/stime"
)

type Handler interface {
	Load() error
	Save() error
}

type HandlerName interface {
	Name() string
}

func handlerName(handler Handler) string {
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

type HandlerInterval interface {
	Interval() time.Duration
}

func handlerInterval(handler Handler) time.Duration {
	i, ok := handler.(HandlerInterval)
	if ok {
		return i.Interval()
	}
	return 0
}

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
			slog.Errorf("Load snap failed, name: %s, error: %v.", handlerName(handler), err)
			return err
		}
		slog.Infof("Load snap succeed, name: %s, cost: %s.", handlerName(handler), timer.Pauses())
	}

	slog.Infof("Load snaps compeleted, cost: %s.", timer.Stops())

	return nil
}

func runWatchSave(ctx context.Context, handler Handler) {
	interval := handlerInterval(handler)

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

func watchSave(handler Handler) error {
	timer := stime.NewTimer()

	err := handler.Save()
	if err != nil {
		slog.Errorf("Save snap failed, name: %s, error: %v.", handlerName(handler), err)
	} else {
		slog.Infof("Save snap succeed, name: %s, cost: %s.", handlerName(handler), timer.Stops())
	}

	return err
}
