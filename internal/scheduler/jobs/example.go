package jobs

import (
	"log/slog"
	"time"
)

// ExampleJob 示例定时任务
type ExampleJob struct{}

// Run 必须实现 cron.Job 接口
func (j *ExampleJob) Run() {
	slog.Info("ExampleJob 执行时间:",
		slog.String("Run time", time.Now().Format("2006-01-02 15:04:05")))
	slog.Info("定时任务执行成功", "job", "ExampleJob")
}
