package stask

import (
	"context"
	"time"

	"github.com/robfig/cron/v3"

	"github.com/yyliziqiu/slib/slog"
	"github.com/yyliziqiu/slib/sreflect"
)

type CronTask struct {
	Name string
	Spec string
	Func func()
}

func (t CronTask) slug() string {
	if t.Name != "" {
		return t.Name
	}
	return sreflect.FuncName(t.Func)
}

func RunCronTasks(ctx context.Context, tasks []CronTask, loc *time.Location) {
	if loc == nil {
		loc = time.Local
	}

	runner := cron.New(cron.WithSeconds(), cron.WithLocation(loc))

	for _, task := range tasks {
		if task.Spec == "" {
			continue
		}
		_, err := runner.AddFunc(task.Spec, task.Func)
		if err != nil {
			slog.Errorf("Add cron task failed, name: %v, error: %v.", task.slug(), err)
			return
		}
		slog.Infof("Add cron task: %s", task.slug())
	}

	runner.Start()
	slog.Info("Cron task started.")
	<-ctx.Done()

	runner.Stop()
	slog.Info("Cron task exit.")
}

func RunCronTasksWithConfig(ctx context.Context, tasks []CronTask, configs []CronTask, loc *time.Location) {
	index := make(map[string]CronTask)
	for _, config := range configs {
		index[config.slug()] = config
	}

	for i := 0; i < len(tasks); i++ {
		config, ok := index[tasks[i].slug()]
		if ok {
			tasks[i].Spec = config.Spec
		}
	}

	RunCronTasks(ctx, tasks, loc)
}
