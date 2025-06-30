package ssnap

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"golang.org/x/net/context"

	"github.com/yyliziqiu/slib/slog"
)

func TestWatchers(t *testing.T) {
	slog.Default, _ = slog.New(slog.Config{Console: true})

	data1 := map[string]string{
		"hello": "world1",
	}
	data2 := map[string]string{
		"hello": "world2",
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	watchers := DefaultWatchers([]DefaultWatcherConfig{
		{
			Path: "/private/ws/self/slib/data/data1.json",
			Data: &data1,
			Mu:   &sync.RWMutex{},
		},
		{
			Path: "/private/ws/self/slib/data/data2.json",
			Data: &data2,
			Poll: 5 * time.Second,
			Mu:   &sync.Mutex{},
		},
	}...)

	err := Watches(ctx, watchers)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(data1)
	fmt.Println(data2)

	time.Sleep(60 * time.Second)

	cancel()

	time.Sleep(time.Second)
}
