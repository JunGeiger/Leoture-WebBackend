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

// Dict 字典数据库操作接口
type Dict interface {
	Create(ctx context.Context, m *model.Dict, id uuid.UUID) error
	Update(ctx context.Context, m *model.Dict, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Dict, error)
	GetByKey(ctx context.Context, category, key string) (*model.Dict, error)
	Count(ctx context.Context, data *request.DictQuery) (int64, error)
	List(ctx context.Context, data *request.DictQuery) ([]*model.Dict, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type dictRepo struct {
	db *pgxpool.Pool
}

func NewDictRepo(db *pgxpool.Pool) Dict {
	return &dictRepo{db: db}
}

func (r *dictRepo) Create(ctx context.Context, m *model.Dict, id uuid.UUID) error {
	args := pgx.NamedArgs{
		"id":         id,
		"name":       m.Name,
		"category":   m.Category,
		"key":        m.Key,
		"val":        m.Val,
		"sort":       m.Sort,
		"is_default": m.IsDefault,
		"status":     m.Status,
		"remark":     m.Remark,
	}
	sql := `INSERT INTO dictionaries
    (id, name, category, key, val, sort, is_default, status, remark) VALUES
    (@id, @name, @category, @key, @val, @sort, @is_default, @status, @remark)`
	_, err := r.db.Exec(ctx, sql, args)
	if err != nil {
		return err
	}
	return nil
}

func (r *dictRepo) Update(ctx context.Context, m *model.Dict, id uuid.UUID) error {
	args := pgx.NamedArgs{
		"id":         id,
		"name":       m.Name,
		"category":   m.Category,
		"key":        m.Key,
		"val":        m.Val,
		"sort":       m.Sort,
		"is_default": m.IsDefault,
		"status":     m.Status,
		"remark":     m.Remark,
	}
	sql := `UPDATE dictionaries SET
	name=@name, category=@category, key=@key, val=@val,
	sort=@sort, is_default=@is_default, status=@status, remark=@remark, updated_at=CURRENT_TIMESTAMP
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

func (r *dictRepo) GetByID(ctx context.Context, id uuid.UUID) (*model.Dict, error) {
	sql := `SELECT id, name, category, key, val, sort, is_default, status, remark, created_at, updated_at
    FROM dictionaries WHERE deleted_at IS NULL AND id=$1`
	var m model.Dict
	err := r.db.QueryRow(ctx, sql, id).Scan(
		&m.ID, &m.Name, &m.Category, &m.Key, &m.Val, &m.Sort,
		&m.IsDefault, &m.Status, &m.Remark, &m.CreatedAt, &m.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *dictRepo) GetByKey(ctx context.Context, category, key string) (*model.Dict, error) {
	sql := `SELECT id, name, category, key, val, sort, is_default, status, remark, created_at, updated_at
    FROM dictionaries WHERE deleted_at IS NULL AND category=$1 AND key=$2`
	var m model.Dict
	err := r.db.QueryRow(ctx, sql, category, key).Scan(
		&m.ID, &m.Name, &m.Category, &m.Key, &m.Val, &m.Sort,
		&m.IsDefault, &m.Status, &m.Remark, &m.CreatedAt, &m.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *dictRepo) Count(ctx context.Context, data *request.DictQuery) (int64, error) {
	var sql strings.Builder
	args := pgx.NamedArgs{}
	sql.WriteString(`SELECT COUNT(*) FROM dictionaries WHERE deleted_at IS NULL`)
	sql.WriteString(buildDictQuerySQL(data, args))

	var total int64
	err := r.db.QueryRow(ctx, sql.String(), args).Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (r *dictRepo) List(ctx context.Context, data *request.DictQuery) ([]*model.Dict, error) {
	var sql strings.Builder
	args := pgx.NamedArgs{}
	sql.WriteString(`SELECT id, name, category, key, val, sort, is_default, status, remark, created_at, updated_at
    FROM dictionaries WHERE deleted_at IS NULL`)
	sql.WriteString(buildDictQuerySQL(data, args))
	sql.WriteString(` ORDER BY sort ASC, created_at DESC`)
	if data.PageSize > 0 && data.CurrentPage > 0 {
		sql.WriteString(` LIMIT @limit OFFSET @offset`)
		args["limit"] = data.PageSize
		args["offset"] = (data.CurrentPage - 1) * data.PageSize
	}

	var list []*model.Dict
	rows, err := r.db.Query(ctx, sql.String(), args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var m model.Dict
		err = rows.Scan(&m.ID, &m.Name, &m.Category, &m.Key, &m.Val, &m.Sort,
			&m.IsDefault, &m.Status, &m.Remark, &m.CreatedAt, &m.UpdatedAt)
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

func (r *dictRepo) Delete(ctx context.Context, id uuid.UUID) error {
	sql := `UPDATE dictionaries SET deleted_at=CURRENT_TIMESTAMP WHERE deleted_at IS NULL AND id=$1`
	tag, err := r.db.Exec(ctx, sql, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func buildDictQuerySQL(req *request.DictQuery, args pgx.NamedArgs) string {
	var sql strings.Builder
	if req.Name != "" {
		sql.WriteString(` AND name LIKE @name`)
		args["name"] = "%" + req.Name + "%"
	}
	if req.Category != "" {
		sql.WriteString(` AND category=@category`)
		args["category"] = req.Category
	}
	if req.Key != "" {
		sql.WriteString(` AND key LIKE @key`)
		args["key"] = "%" + req.Key + "%"
	}
	if req.Status != 0 {
		sql.WriteString(` AND status=@status`)
		args["status"] = req.Status
	}
	if req.DateStart != nil && req.DateEnd != nil {
		sql.WriteString(` AND created_at >= @start AND created_at <= @end`)
		args["start"] = req.DateStart
		args["end"] = req.DateEnd
	}
	return sql.String()
}
