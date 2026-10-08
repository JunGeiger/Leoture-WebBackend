package scheduler

import (
	"LeotureWeb/internal/config"
	"LeotureWeb/internal/scheduler/jobs"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
)

var (
	scheduler *CronScheduler
	once      sync.Once
)

// CronScheduler cron 调度器封装
type CronScheduler struct {
	cron *cron.Cron
}

// SetupScheduler 获取单例调度器
func SetupScheduler(cfg config.Scheduler) error {
	if !cfg.Enabled {
		return nil
	}
	var initErr error
	once.Do(func() {
		// 解析时区
		tz, err := time.LoadLocation(cfg.Timezone)
		if err != nil {
			initErr = fmt.Errorf("scheduler: failed to load timezone '%s': %w", cfg.Timezone, err)
			return
		}
		scheduler = &CronScheduler{
			// 启用秒级精度
			cron: cron.New(cron.WithSeconds(), cron.WithLocation(tz)),
		}

		// 注册任务
		registerJobs(cfg.Jobs, scheduler.cron, buildRegistry())
		scheduler.cron.Start()
	})
	if initErr != nil {
		slog.Info("scheduler started", slog.String("timezone", cfg.Timezone))
	}
	return initErr
}

// 构建任务，需要与配置文件中的名称一致
func buildRegistry() map[string]func() cron.Job {
	return map[string]func() cron.Job{
		"example_job": func() cron.Job { return &jobs.ExampleJob{} },
		// 新增任务只需在此处添加
	}
}

// Stop 停止调度器（优雅关闭）
func Stop(timeout time.Duration) {
	if scheduler == nil {
		return
	}
	ctx := scheduler.cron.Stop()
	select {
	case <-ctx.Done():
		slog.Info("scheduler stopped")
	case <-time.After(timeout):
		slog.Warn("scheduler stop timeout, some jobs may still be running")
	}
	scheduler = nil
}

// ListEntries 获取所有任务信息
func ListEntries() []cron.Entry {
	if scheduler == nil {
		return []cron.Entry{}
	}
	return scheduler.cron.Entries()
}
