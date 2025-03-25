package stask

import (
	"context"

	"github.com/yyliziqiu/slib/slog"
	"github.com/yyliziqiu/slib/sreflect"
)

type OnceTask struct {
	Name string                    // 名称
	Cons int                       // 并发数
	Func func(ctx context.Context) // 方法
}

func (t OnceTask) slug() string {
	if t.Name != "" {
		return t.Name
	}
	return sreflect.FuncName(t.Func)
}

func StartOnceTasks(ctx context.Context, tasks []OnceTask) {
	for _, task := range tasks {
		if task.Cons <= 0 {
			continue
		}
		for i := 0; i < task.Cons; i++ {
			go task.Func(ctx)
		}
		slog.Infof("Add once task: %s (%d)", task.slug(), task.Cons)
	}
}

func StartOnceTasksWithConfig(ctx context.Context, tasks []OnceTask, configs []OnceTask) {
	index := make(map[string]OnceTask)
	for _, config := range configs {
		index[config.slug()] = config
	}

	for i := 0; i < len(tasks); i++ {
		config, ok := index[tasks[i].slug()]
		if ok {
			tasks[i].Cons = config.Cons
		}
	}

	StartOnceTasks(ctx, tasks)
}
