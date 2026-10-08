package repository

import (
	"LeotureWeb/internal/model"
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserSession interface {
	Create(ctx context.Context, m *model.UserSession) error
	UpdateStatusByID(ctx context.Context, status int16, id uuid.UUID) error
	UpdateExpiredByUserID(ctx context.Context, status int16, UserID uuid.UUID) error
	GetByIDAndUserID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*model.UserSession, error)
}

type userSessionRepo struct {
	db *pgxpool.Pool
}

func NewUserSessionRepo(db *pgxpool.Pool) UserSession {
	return &userSessionRepo{db: db}
}

func (r userSessionRepo) Create(ctx context.Context, m *model.UserSession) error {
	args := pgx.NamedArgs{
		"id":           m.ID,
		"user_id":      m.UserID,
		"refresh_jti":  m.RefreshJTI,
		"login_log_id": m.LoginLogID,
		"status":       m.Status,
	}
	sql := `INSERT INTO user_session 
    (id, user_id, refresh_jti, login_log_id, status) VALUES 
    (@id, @user_id, @refresh_jti, @login_log_id, @status)`
	_, err := r.db.Exec(ctx, sql, args)
	if err != nil {
		return err
	}
	return nil
}

func (r userSessionRepo) UpdateStatusByID(ctx context.Context, status int16, id uuid.UUID) error {
	sql := `UPDATE user_session SET status = $1, updated_at = CURRENT_TIMESTAMP WHERE id=$2`
	_, err := r.db.Exec(ctx, sql, status, id)
	if err != nil {
		return err
	}
	return nil
}

func (r userSessionRepo) UpdateExpiredByUserID(ctx context.Context, status int16, UserID uuid.UUID) error {
	sql := `UPDATE user_session SET status = $1, updated_at = CURRENT_TIMESTAMP WHERE status='1' AND user_id=$2`
	_, err := r.db.Exec(ctx, sql, status, UserID)
	if err != nil {
		return err
	}
	return nil
}

func (r userSessionRepo) GetByIDAndUserID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*model.UserSession, error) {
	sql := `SELECT id, user_id, refresh_jti, login_log_id, status, created_at, updated_at FROM user_session 
    WHERE id=$1 AND user_id=$2`
	var m model.UserSession
	err := r.db.QueryRow(ctx, sql, id, userID).Scan(
		&m.ID, &m.UserID, &m.RefreshJTI, &m.LoginLogID, &m.Status, &m.CreatedAt, &m.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &m, nil
}
