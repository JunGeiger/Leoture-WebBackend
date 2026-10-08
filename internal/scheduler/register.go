package scheduler

import (
	"LeotureWeb/internal/config"
	"github.com/robfig/cron/v3"
	"log/slog"
	"runtime/debug"
)

// safeJob wraps a cron.Job to recover from panics
// and prevent the cron scheduler from crashing.
type safeJob struct {
	jobName string
	job     cron.Job
}

func (s safeJob) Run() {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("scheduler: Cronjob panic",
				slog.String("job", s.jobName),
				slog.String("stack", string(debug.Stack())),
			)
		}
	}()
	s.job.Run()
}

// registerJobs 注册所有定时任务
func registerJobs(jobs map[string]config.CronJob, c *cron.Cron, registry map[string]func() cron.Job) {
	if len(jobs) == 0 {
		return
	}

	for name, cfg := range jobs {
		if !cfg.Enabled {
			slog.Info("scheduler: Skipping disabled cronjob", slog.String("job", name))
			continue
		}

		factory, ok := registry[name]
		if !ok {
			slog.Error("scheduler: Unknown cronjob", slog.String("job", name))
			continue
		}

		addJob(c, factory(), name, cfg)
	}
}

func addJob(c *cron.Cron, j cron.Job, jobName string, jobCfg config.CronJob) {
	// 防止 panic 导致整个 cron 崩溃
	wrapped := safeJob{jobName: jobName, job: j}

	// jobCfg.Spec 为秒级 cron 表达式
	id, err := c.AddJob(jobCfg.Spec, wrapped)

	if err != nil {
		slog.Error("scheduler: cronjob register failed",
			slog.Any("error", err),
			slog.String("job", jobName),
			slog.String("spec", jobCfg.Spec),
			slog.String("desc", jobCfg.Desc))
		return
	}
	slog.Info("scheduler: Cronjob registered",
		slog.String("job", jobName),
		slog.String("spec", jobCfg.Spec),
		slog.String("desc", jobCfg.Desc),
		slog.Any("entry_id", id))
}
