package repository

import (
	"LeotureWeb/internal/model"
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// LoginLog 登录日志数据库操作接口
type LoginLog interface {
	Create(ctx context.Context, m *model.LoginLog) error
}

type loginLogRepo struct {
	db *pgxpool.Pool
}

func NewLoginLogRepo(db *pgxpool.Pool) LoginLog {
	return &loginLogRepo{db: db}
}

func (r *loginLogRepo) Create(ctx context.Context, m *model.LoginLog) error {
	args := pgx.NamedArgs{
		"id":         m.ID,
		"username":   m.Username,
		"ip_address": m.IPAddress,
		"user_agent": m.UserAgent,
		"origin_url": m.OriginUrl,
		"location":   m.Location,
		"result":     m.Result,
	}
	sql := `INSERT INTO login_log 
    (id, username, ip_address, user_agent, origin_url, location, result) VALUES 
    (@id, @username, @ip_address, @user_agent, @origin_url, @location, @result)`
	_, err := r.db.Exec(ctx, sql, args)
	if err != nil {
		return err
	}
	return nil
}
