package repository

import (
	"LeotureWeb/internal/model"
	"LeotureWeb/internal/types/request"
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"strings"
)

// Menu 菜单数据库操作接口
type Menu interface {
	Create(ctx context.Context, m *model.Menu, id uuid.UUID) error
	Update(ctx context.Context, m *model.Menu, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Menu, error)
	GetByIDs(ctx context.Context, dataIDs []uuid.UUID) ([]*model.Menu, error)
	GetAll(ctx context.Context) ([]*model.Menu, error)
	Count(ctx context.Context, data *request.MenuQuery) (int64, error)
	List(ctx context.Context, data *request.MenuQuery) ([]*model.Menu, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type menuRepo struct {
	db *pgxpool.Pool
}

func NewMenuRepo(db *pgxpool.Pool) Menu {
	return &menuRepo{db: db}
}

func (r *menuRepo) Create(ctx context.Context, m *model.Menu, id uuid.UUID) error {
	args := pgx.NamedArgs{
		"id":           m.ID,
		"parent_id":    m.ParentID,
		"type":         m.Type,
		"name":         m.Name,
		"title":        m.Title,
		"keep_alive":   m.KeepAlive,
		"icon":         m.Icon,
		"order_no":     m.OrderNo,
		"external_url": m.ExternalUrl,
		"is_hidden":    m.IsHidden,
		"perm_code":    m.PermCode,
		"api_path":     m.ApiPath,
		"status":       m.Status,
	}
	sql := `INSERT INTO menus
    (id, parent_id, type, name, title, keep_alive, icon, order_no, external_url, is_hidden, perm_code, api_path, status) VALUES
    (@id, @parent_id, @type, @name, @title, @keep_alive, @icon, @order_no, @external_url, @is_hidden, @perm_code, @api_path, @status)`

	_, err := r.db.Exec(ctx, sql, args)
	return err
}

func (r *menuRepo) Update(ctx context.Context, m *model.Menu, id uuid.UUID) error {
	args := pgx.NamedArgs{
		"id":           m.ID,
		"parent_id":    m.ParentID,
		"type":         m.Type,
		"name":         m.Name,
		"title":        m.Title,
		"keep_alive":   m.KeepAlive,
		"icon":         m.Icon,
		"order_no":     m.OrderNo,
		"external_url": m.ExternalUrl,
		"is_hidden":    m.IsHidden,
		"perm_code":    m.PermCode,
		"api_path":     m.ApiPath,
		"status":       m.Status,
	}
	sql := `UPDATE menus SET
	id=@id, parent_id=@parent_id, type=@type, name=@name, title=@title, keep_alive=@keep_alive, icon=@icon,
	order_no=@order_no, external_url=@external_url, is_hidden=@is_hidden, perm_code=@perm_code, api_path=@api_path,
	status=@status, updated_at=CURRENT_TIMESTAMP
    WHERE deleted_at IS NULL AND id=@id`

	tag, err := r.db.Exec(ctx, sql, args)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (r *menuRepo) GetByID(ctx context.Context, id uuid.UUID) (*model.Menu, error) {
	sql := `SELECT id, parent_id, type, name, title, keep_alive, icon, order_no, external_url, is_hidden, perm_code,
       api_path, status, created_at, updated_at FROM menus WHERE deleted_at IS NULL AND id=$1 ORDER BY sort ASC`

	var m model.Menu
	err := r.db.QueryRow(ctx, sql, id).Scan(
		&m.ID, &m.ParentID, &m.Type, &m.Name, &m.Title, &m.KeepAlive, &m.Icon, &m.OrderNo, &m.ExternalUrl, &m.IsHidden, &m.PermCode,
		&m.ApiPath, &m.Status, &m.CreatedAt, &m.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *menuRepo) GetByIDs(ctx context.Context, ids []uuid.UUID) ([]*model.Menu, error) {
	sql := `SELECT id, parent_id, type, name, title, keep_alive, icon, order_no, external_url, is_hidden, perm_code,
       api_path, status, created_at, updated_at FROM menus WHERE deleted_at IS NULL AND id=ANY($1::UUID[]) ORDER BY sort ASC`

	rows, err := r.db.Query(ctx, sql, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	models := make([]*model.Menu, 0, len(ids))
	for rows.Next() {
		var m model.Menu
		err := rows.Scan(
			&m.ID, &m.ParentID, &m.Type, &m.Name, &m.Title, &m.KeepAlive, &m.Icon, &m.OrderNo, &m.ExternalUrl, &m.IsHidden, &m.PermCode,
			&m.ApiPath, &m.Status, &m.CreatedAt, &m.UpdatedAt)
		if err != nil {
			return nil, err
		}
		models = append(models, &m)
	}
	return models, nil
}

func (r *menuRepo) GetAll(ctx context.Context) ([]*model.Menu, error) {
	sql := `SELECT id, parent_id, type, name, title, keep_alive, icon, order_no, external_url, is_hidden, perm_code,
       api_path, status, created_at, updated_at FROM menus WHERE deleted_at IS NULL ORDER BY sort ASC`

	rows, err := r.db.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*model.Menu
	for rows.Next() {
		var m model.Menu
		err = rows.Scan(
			&m.ID, &m.ParentID, &m.Type, &m.Name, &m.Title, &m.KeepAlive, &m.Icon, &m.OrderNo, &m.ExternalUrl, &m.IsHidden, &m.PermCode,
			&m.ApiPath, &m.Status, &m.CreatedAt, &m.UpdatedAt)
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

func (r *menuRepo) Count(ctx context.Context, req *request.MenuQuery) (int64, error) {
	var sql strings.Builder
	args := pgx.NamedArgs{}
	sql.WriteString(`SELECT COUNT(*) FROM menus WHERE deleted_at IS NULL`)
	sql.WriteString(buildMenuQuerySQL(req, args))

	var total int64
	err := r.db.QueryRow(ctx, sql.String(), args).Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (r *menuRepo) List(ctx context.Context, req *request.MenuQuery) ([]*model.Menu, error) {
	var sql strings.Builder
	args := pgx.NamedArgs{}
	sql.WriteString(`SELECT id, parent_id, type, name, title, keep_alive, icon, order_no, external_url, is_hidden, perm_code,
       api_path, status, created_at, updated_at FROM menus WHERE deleted_at IS NULL`)

	sql.WriteString(buildMenuQuerySQL(req, args))
	sql.WriteString(` ORDER BY sort ASC`)

	if req.PageSize > 0 && req.CurrentPage > 0 {
		sql.WriteString(` LIMIT @limit OFFSET @offset`)
		args["limit"] = req.PageSize
		args["offset"] = (req.CurrentPage - 1) * req.PageSize
	}

	rows, err := r.db.Query(ctx, sql.String(), args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*model.Menu
	for rows.Next() {
		var m model.Menu
		err = rows.Scan(
			&m.ID, &m.ParentID, &m.Type, &m.Name, &m.Title, &m.KeepAlive, &m.Icon, &m.OrderNo, &m.ExternalUrl, &m.IsHidden, &m.PermCode,
			&m.ApiPath, &m.Status, &m.CreatedAt, &m.UpdatedAt)
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

func (r *menuRepo) Delete(ctx context.Context, id uuid.UUID) error {
	sql := `UPDATE menus SET deleted_at=CURRENT_TIMESTAMP WHERE deleted_at IS NULL AND id=$1`
	tag, err := r.db.Exec(ctx, sql, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func buildMenuQuerySQL(data *request.MenuQuery, args pgx.NamedArgs) string {
	var sql strings.Builder
	if data.Name != "" {
		sql.WriteString(` AND name LIKE @name`)
		args["name"] = "%" + data.Name + "%"
	}
	if data.Type == 0 {
		sql.WriteString(` AND type=@type`)
		args["type"] = data.Type
	}
	if data.ParentID != uuid.Nil {
		sql.WriteString(` AND parent_id=@parent_id`)
		args["parent_id"] = data.ParentID
	}
	if data.Status != 0 {
		sql.WriteString(` AND status=@status`)
		args["status"] = data.Status
	}
	if data.DateRange.DateStart != nil && data.DateRange.DateEnd != nil {
		sql.WriteString(` AND created_at >= @start AND created_at <= @end`)
		args["start"] = data.DateRange.DateStart
		args["end"] = data.DateRange.DateEnd
	}
	return sql.String()
}
