package main

import (
	"LeotureWeb/docs/swagger"
	"LeotureWeb/internal/authorization"
	"LeotureWeb/internal/cache"
	"LeotureWeb/internal/config"
	"LeotureWeb/internal/database"
	"LeotureWeb/internal/handler"
	"LeotureWeb/internal/logger"
	"LeotureWeb/internal/repository"
	"LeotureWeb/internal/router"
	"LeotureWeb/internal/scheduler"
	"LeotureWeb/internal/service"
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// @title           LeotureWeb API
// @version         1.0
// @description     这是一个基于 Gin 框架的 Web 项目 API 文档
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey  BearerAuth
// @in                          header
// @name                        Authorization
// @description                 Bearer Token for authentication
func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	// 初始化配置
	cfg, err := config.SetupViper()
	if err != nil {
		log.Fatalf("Setup viper error: %v", err)
	}

	// 初始化日志
	if err := logger.SetupSlog(*cfg.Log); err != nil {
		log.Fatalf("Setup slog error: %v", err)
	}

	// 初始化本地单机定时任务
	if err := scheduler.SetupScheduler(*cfg.Scheduler); err != nil {
		log.Fatalf("Setup scheduler error: %v", err)
	}

	// 服务启动日志
	appCfg := cfg.App
	slog.Info("Application is starting",
		slog.String("name", appCfg.Name),
		slog.String("version", appCfg.Version),
		slog.String("env", appCfg.Env))

	// 初始化 pgsql 客户端
	pgxPool, err := database.SetupPostgres(*cfg.Database)
	if err != nil {
		log.Fatalf("Setup postgres error: %v", err)
	}
	repoSet := repository.NewRepoSet(pgxPool)

	// 初始化 Reids 客户端
	redisCli, err := cache.SetupRedisCli(*cfg.Redis)
	if err != nil {
		log.Fatalf("Setup redis error: %v", err)
	}

	// 初始化权限管理JWT、Casbin
	jwtService, err := authorization.SetupJWT(*cfg.JWT)
	if err != nil {
		log.Fatalf("Setup JWT error: %v", err)
	}
	casbinService, err := authorization.SetupCasbin(*cfg.Casbin, pgxPool, redisCli)
	if err != nil {
		log.Fatalf("Setup Casbin error: %v", err)
	}

	// 初始化Gin路由，同步初始化repo、service、handler依赖
	serviceSet := service.NewSet(repoSet, jwtService, casbinService)
	handlerSet := handler.NewSet(serviceSet)
	r := router.SetupGinRouter(*cfg, handlerSet, jwtService, casbinService)

	// 初始化Swagger
	swagger.SetupSwagger(r, *cfg)

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler: r,
	}
	go func() {
		slog.Info("Server started", slog.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server error: %v", err)
		}
	}()

	<-ctx.Done()
	slog.Info("Shutdown signal received")

	// 停止定时任组件
	scheduler.Stop(10 * time.Second)

	// 关闭资源
	database.Close()
	cache.Close()

	slog.Info("Application exited gracefully")
}
