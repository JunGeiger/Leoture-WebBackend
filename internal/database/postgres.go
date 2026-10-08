package database

import (
	"LeotureWeb/internal/config"
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// DB 全局私有变量，存储连接池实例
var (
	pg   *pgxpool.Pool
	once sync.Once
)

// SetupPostgres 初始化 PostgreSQL 连接池 (线程安全单例模式)
// 参数:
//
//	cfg: 数据库配置结构体
//	logger: 日志实例，用于记录连接状态
func SetupPostgres(cfg config.Database) (*pgxpool.Pool, error) {
	//	ctx: 上下文，用于控制初始化超时
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var setupErr error
	// sync.Once 确保闭包内的代码仅执行一次
	once.Do(func() {
		setupErr = setup(ctx, cfg)
	})

	return pg, setupErr
}

func setup(ctx context.Context, cfg config.Database) error {
	// 构建 DSN (Data Source Name)
	// 格式: postgres://username:password@host:port/database?sslmode=disable
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
		cfg.SSLMode,
	)

	// 解析配置
	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return fmt.Errorf("postgres: configuration resolution failed: %w", err)
	}

	// 配置连接池参数 (关键性能调优点)
	poolConfig.MaxConns = cfg.MaxConns
	poolConfig.MinConns = cfg.MinConns
	poolConfig.MaxConnLifetime = cfg.MaxConnLifetime
	poolConfig.MaxConnIdleTime = cfg.MaxConnIdleTime
	poolConfig.HealthCheckPeriod = cfg.HealthCheckPeriod
	poolConfig.ConnConfig.ConnectTimeout = cfg.ConnectTimeout

	// 运行时参数
	poolConfig.ConnConfig.RuntimeParams = map[string]string{
		"statement_timeout":                   strconv.Itoa(int(cfg.StatementTimeout.Milliseconds())),
		"lock_timeout":                        strconv.Itoa(int(cfg.LockTimeout.Milliseconds())),
		"idle_in_transaction_session_timeout": strconv.Itoa(int(cfg.IdleInTransactionSessionTimeout.Milliseconds())),
	}

	// 创建连接池
	pg, err = pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return fmt.Errorf("postgres: failed to create connection pool: %w", err)
	}

	// Ping 数据库验证连接有效性
	if err := pg.Ping(ctx); err != nil {
		pg.Close() // 如果Ping不通，关闭已创建的连接池
		return fmt.Errorf("postgres: connection verification failed: %w", err)
	}

	slog.Info("postgres: connection succeeded",
		slog.String("host", cfg.Host),
		slog.Int("port", cfg.Port),
		slog.String("database", cfg.DBName),
	)
	return nil
}

// Close 关闭数据库连接池，建议在程序退出时调用
func Close() {
	if pg != nil {
		pg.Close()
		pg = nil
	}
	slog.Info("postgres: connection pool closed")
}
