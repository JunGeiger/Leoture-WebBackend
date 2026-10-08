package repository

import (
	"LeotureWeb/internal/model"
	"LeotureWeb/internal/types/request"
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Role interface {
	Create(ctx context.Context, m *model.Role, id uuid.UUID) error
	Update(ctx context.Context, m *model.Role, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Role, error)
	GetByIDs(ctx context.Context, dataIDs []uuid.UUID) ([]*model.Role, error)
	GetByCode(ctx context.Context, code string) (*model.Role, error)
	Count(ctx context.Context, data *request.RoleQuery) (int64, error)
	List(ctx context.Context, data *request.RoleQuery) ([]*model.Role, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type roleRepo struct {
	db *pgxpool.Pool
}

func NewRoleRepo(db *pgxpool.Pool) Role {
	return &roleRepo{db: db}
}

func (r *roleRepo) Create(ctx context.Context, m *model.Role, id uuid.UUID) error {
	args := pgx.NamedArgs{
		"id":     id,
		"name":   m.Name,
		"code":   m.Code,
		"desc":   m.Remark,
		"status": m.Status,
	}
	sql := `INSERT INTO roles
    (id, name, code, remark, status) VALUES
    (@id, @name, @code, @remark, @status)`
	_, err := r.db.Exec(ctx, sql, args)
	if err != nil {
		return err
	}
	return nil
}

func (r *roleRepo) Update(ctx context.Context, m *model.Role, id uuid.UUID) error {
	args := pgx.NamedArgs{
		"id":     id,
		"name":   m.Name,
		"code":   m.Code,
		"remark": m.Remark,
		"status": m.Status,
	}
	sql := `UPDATE roles SET
	name=@name, code=@code, remark=@remark, status=@status, updated_at=CURRENT_TIMESTAMP
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

func (r *roleRepo) GetByID(ctx context.Context, id uuid.UUID) (*model.Role, error) {
	sql := `SELECT id, name, code, remark, status, created_at, updated_at FROM roles
    WHERE deleted_at IS NULL AND id=$1`
	var m model.Role
	err := r.db.QueryRow(ctx, sql, id).Scan(
		&m.ID, &m.Name, &m.Code, &m.Remark, &m.Status, &m.CreatedAt, &m.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *roleRepo) GetByIDs(ctx context.Context, dataIDs []uuid.UUID) ([]*model.Role, error) {
	sql := `SELECT id, name, code, remark, status, created_at, updated_at FROM roles
    WHERE deleted_at IS NULL AND id = ANY($1::UUID[])`
	rows, err := r.db.Query(ctx, sql, dataIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	models := make([]*model.Role, 0, len(dataIDs))
	for rows.Next() {
		var m model.Role
		err := rows.Scan(&m.ID, &m.Name, &m.Code, &m.Remark, &m.Status, &m.CreatedAt, &m.UpdatedAt)
		if err != nil {
			return nil, err
		}
		models = append(models, &m)
	}
	return models, nil
}

func (r *roleRepo) GetByCode(ctx context.Context, code string) (*model.Role, error) {
	sql := `SELECT id, name, code, remark, status, created_at, updated_at FROM roles
    WHERE deleted_at IS NULL AND code=$1`
	var m model.Role
	err := r.db.QueryRow(ctx, sql, code).Scan(
		&m.ID, &m.Name, &m.Code, &m.Remark, &m.Status, &m.CreatedAt, &m.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *roleRepo) Count(ctx context.Context, data *request.RoleQuery) (int64, error) {
	var sql strings.Builder
	args := pgx.NamedArgs{}
	sql.WriteString(`SELECT COUNT(*) FROM roles WHERE deleted_at IS NULL`)
	sql.WriteString(buildQuerySQL(data, args))

	var total int64
	err := r.db.QueryRow(ctx, sql.String(), args).Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (r *roleRepo) List(ctx context.Context, data *request.RoleQuery) ([]*model.Role, error) {
	var sql strings.Builder
	args := pgx.NamedArgs{}
	sql.WriteString(`SELECT id, name, code, remark, status, created_at, updated_at
    FROM roles WHERE deleted_at IS NULL`)
	sql.WriteString(buildQuerySQL(data, args))

	sql.WriteString(` ORDER BY created_at DESC`)
	if data.PageSize > 0 && data.CurrentPage > 0 {
		sql.WriteString(` LIMIT @limit OFFSET @offset`)
		args["limit"] = data.PageSize
		args["offset"] = (data.CurrentPage - 1) * data.PageSize
	}

	var list []*model.Role
	rows, err := r.db.Query(ctx, sql.String(), args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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

func (r *roleRepo) Delete(ctx context.Context, id uuid.UUID) error {
	sql := `UPDATE roles SET deleted_at=CURRENT_TIMESTAMP WHERE deleted_at IS NULL AND id=$1`
	tag, err := r.db.Exec(ctx, sql, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func buildQuerySQL(data *request.RoleQuery, args pgx.NamedArgs) string {
	var sql strings.Builder
	if data.Name != "" {
		sql.WriteString(` AND name=@name`)
		args["name"] = data.Name
	}
	if data.Code != "" {
		sql.WriteString(` AND code=@code`)
		args["code"] = data.Code
	}
	if data.Remark != "" {
		sql.WriteString(` AND remark LIKE @remark`)
		args["remark"] = "%" + data.Remark + "%"
	}
	if data.Status != 0 {
		sql.WriteString(` AND status=@status`)
		args["status"] = data.Status
	}
	if data.DateStart != nil && data.DateEnd != nil {
		sql.WriteString(` AND created_at >= @start AND created_at <= @end`)
		args["start"] = data.DateStart
		args["end"] = data.DateEnd
	}
	return sql.String()
}
