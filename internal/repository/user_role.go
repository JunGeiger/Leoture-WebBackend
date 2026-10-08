package repository

import (
	"LeotureWeb/internal/model"
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRole interface {
	Create(ctx context.Context, data []*model.UserRole) error
	Delete(ctx context.Context, userID uuid.UUID) error
	GetRolesByUserId(ctx context.Context, userID uuid.UUID) ([]*model.Role, error)
}

type userRoleRepo struct {
	db *pgxpool.Pool
}

func NewUserRoleRepo(db *pgxpool.Pool) UserRole {
	return &userRoleRepo{db: db}
}

func (r userRoleRepo) Create(ctx context.Context, data []*model.UserRole) error {
	batch := &pgx.Batch{}
	for _, userRole := range data {
		batch.Queue(
			`INSERT INTO user_role (user_id, role_id) VALUES ($1, $2)`,
			userRole.UserID, userRole.RoleID,
		)
	}

	br := r.db.SendBatch(ctx, batch)
	defer func() {
		if err := br.Close(); err != nil {
			slog.ErrorContext(ctx, "create user role correlations, failed to close database batch",
				slog.String("error", err.Error()),
			)
		}
	}()

	for i := 0; i < len(data); i++ {
		if _, err := br.Exec(); err != nil {
			return fmt.Errorf("create user role correlations, userId: %s, roleId: %s, exec batch error: %w",
				data[i].UserID.String(), data[i].RoleID.String(), err)
		}
	}
	return nil
}

func (r userRoleRepo) Delete(ctx context.Context, userID uuid.UUID) error {
	sql := `DELETE FROM user_role WHERE user_id = $1`
	tag, err := r.db.Exec(ctx, sql, userID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (r userRoleRepo) GetRolesByUserId(ctx context.Context, userID uuid.UUID) ([]*model.Role, error) {
	sql := `SELECT r.id, r.name, r.code, r.remark, r.status, r.created_at, r.updated_at FROM user_role ur
    	LEFT JOIN roles r ON ur.role_id = r.id WHERE ur.user_id = $1`
	rows, err := r.db.Query(ctx, sql, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*model.Role
	for rows.Next() {
		var m model.Role
		err = rows.Scan(&m.ID, &m.Name, &m.Code, &m.Remark, &m.Status, &m.CreatedAt, &m.UpdatedAt)
		if err != nil {
			return nil, err
		}
		list = append(list, &m)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return list, nil
}
