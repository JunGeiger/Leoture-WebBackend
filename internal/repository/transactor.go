package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
)

type Transaction interface {
	WithDeferrable() TransactionOption
	WithReadWrite() TransactionOption
	WithReadOnly() TransactionOption
	WithIsoLevel(level pgx.TxIsoLevel) TransactionOption
	WithTransaction(ctx context.Context, fn func(tx pgx.Tx) error, opts ...TransactionOption) error
}

type transactor struct {
	db *pgxpool.Pool
}

func NewTransactor(db *pgxpool.Pool) Transaction {
	return &transactor{db: db}
}

// transactionConfig 内部配置结构体，封装 pgx.TxOptions 并为未来扩展预留字段
type transactionConfig struct {
	pgx.TxOptions
	// 未来可在此处扩展非 pgx 原生的配置，例如：
	// MetricsLabel string
	// SlowThreshold time.Duration
}

// TransactionOption 函数选项类型
type TransactionOption func(*transactionConfig)

// WithIsoLevel 设置事务隔离级别
func (t transactor) WithIsoLevel(level pgx.TxIsoLevel) TransactionOption {
	return func(c *transactionConfig) {
		c.IsoLevel = level
	}
}

// WithReadOnly 设置只读事务模式
func (t transactor) WithReadOnly() TransactionOption {
	return func(c *transactionConfig) {
		c.AccessMode = pgx.ReadOnly
	}
}

// WithReadWrite 显式设置读写事务模式（默认即为读写，通常不需要调用）
func (t transactor) WithReadWrite() TransactionOption {
	return func(c *transactionConfig) {
		c.AccessMode = pgx.ReadWrite
	}
}

// WithDeferrable 设置 DEFERRABLE 模式
// 注意：仅在 IsoLevel=Serializable 且 AccessMode=ReadOnly 时生效
func (t transactor) WithDeferrable() TransactionOption {
	return func(c *transactionConfig) {
		c.DeferrableMode = pgx.Deferrable
	}
}

// WithTransaction 在事务中执行 fn
// opts 为可选的事务配置，不传时使用数据库默认配置（ReadCommitted + ReadWrite）
func (t transactor) WithTransaction(ctx context.Context, fn func(tx pgx.Tx) error, opts ...TransactionOption) error {
	// 1. 初始化默认配置（零值即 pgx 默认安全行为）
	cfg := &transactionConfig{}

	// 2. 应用所有传入的选项
	for _, opt := range opts {
		opt(cfg)
	}

	// 3. 开启事务
	tx, err := t.db.BeginTx(ctx, cfg.TxOptions)
	if err != nil {
		return fmt.Errorf("database: failed beginning pgx transaction: %w", err)
	}

	// 4. panic 保护：确保异常时连接被正确释放
	defer func() {
		if r := recover(); r != nil {
			// 使用不受原 ctx 取消影响的上下文执行回滚
			cleanupCtx := context.WithoutCancel(ctx)
			_ = tx.Rollback(cleanupCtx)
			// 重新抛出异常，方式无法跟踪原始错误
			panic(r)
		}
	}()

	// 5. 执行业务逻辑
	if err := fn(tx); err != nil {
		// 业务错误时回滚，使用独立上下文避免 ctx 已取消导致回滚失败
		cleanupCtx := context.WithoutCancel(ctx)
		if rbErr := tx.Rollback(cleanupCtx); rbErr != nil && !errors.Is(rbErr, pgx.ErrTxClosed) {
			slog.ErrorContext(cleanupCtx, "database: failed rolling back pgx transaction",
				slog.String("fn_err", err.Error()),
				slog.String("rb_err", rbErr.Error()),
			)
		}
		return err
	}

	// 6、提交事务；Commit 失败后无需再 Rollback（pgx 内部已处理）
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("database: failed committing pgx transaction: %w", err)
	}

	return nil
}
