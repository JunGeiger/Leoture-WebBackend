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

type RoleMenu interface {
	Create(ctx context.Context, data []*model.RoleMenu) error
	Delete(ctx context.Context, roleID uuid.UUID) error
	GetMenusBYRoleID(ctx context.Context, roleID uuid.UUID) ([]*model.Menu, error)
	GetMenusBYRoleIDs(ctx context.Context, menuTypes []int16, roleIDs []uuid.UUID) ([]*model.Menu, error)
}

type roleMenuRepo struct {
	db *pgxpool.Pool
}

func NewRoleMenuRepo(db *pgxpool.Pool) RoleMenu {
	return &roleMenuRepo{db: db}
}

func (r roleMenuRepo) Create(ctx context.Context, data []*model.RoleMenu) error {
	batch := &pgx.Batch{}
	for _, roleMenu := range data {
		batch.Queue(
			`INSERT INTO role_menu (role_id, menu_id) VALUES ($1, $2)`,
			roleMenu.RoleID, roleMenu.MenuID,
		)
	}
	br := r.db.SendBatch(ctx, batch)
	defer func() {
		if err := br.Close(); err != nil {
			slog.ErrorContext(ctx, "create role menu correlations, failed to close database batch",
				slog.String("error", err.Error()),
			)
		}
	}()
	for i := 0; i < len(data); i++ {
		if _, err := br.Exec(); err != nil {
			return fmt.Errorf("create role menu correlations, userId: %s, roleId: %s, exec batch error: %w",
				data[i].RoleID.String(), data[i].MenuID.String(), err)
		}
	}
	return nil
}

func (r roleMenuRepo) Delete(ctx context.Context, roleID uuid.UUID) error {
	sql := `DELETE FROM role_menu WHERE role_id = $1`
	tag, err := r.db.Exec(ctx, sql, roleID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (r roleMenuRepo) GetMenusBYRoleID(ctx context.Context, roleID uuid.UUID) ([]*model.Menu, error) {
	sql := `SELECT m.id, m.name, m.icon, m.type, m.parent_id, m.perm_code, m.path, m.component, m.order_no,
    m.status, m.created_at, m.updated_at FROM role_menu rm LEFT JOIN menus m ON rm.menu_id=m.id WHERE rm.role_id=$1`

	rows, err := r.db.Query(ctx, sql, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*model.Menu
	for rows.Next() {
		var m model.Menu
		err = rows.Scan(
			&m.ID, &m.Name, &m.Icon, &m.Type, &m.ParentID, &m.PermCode, &m.Path, &m.Component, &m.OrderNo,
			&m.Status, &m.CreatedAt, &m.UpdatedAt)
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

func (r roleMenuRepo) GetMenusBYRoleIDs(ctx context.Context, menuTypes []int16, roleIDs []uuid.UUID) ([]*model.Menu, error) {
	sql := `SELECT DISTINCT m.id, m.name, m.icon, m.type, m.parent_id, m.perm_code, m.path, m.component, m.order_no,
    m.status, m.created_at, m.updated_at FROM role_menu rm INNER JOIN menus m ON rm.menu_id = m.id WHERE m.type = ANY($1) AND rm.role_id = ANY($2)`

	rows, err := r.db.Query(ctx, sql, menuTypes, roleIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[model.Menu])
	if err != nil {
		return nil, err
	}

	return list, nil
}
