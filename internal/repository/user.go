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

type User interface {
	Create(ctx context.Context, m *model.User, id uuid.UUID) error
	Update(ctx context.Context, m *model.User, id uuid.UUID) error
	UpdatePassword(ctx context.Context, id uuid.UUID, old string, new string) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	Count(ctx context.Context, data *request.UserQuery) (int64, error)
	List(ctx context.Context, data *request.UserQuery) ([]*model.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type userRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) User {
	return &userRepo{db: db}
}

func (r *userRepo) Create(ctx context.Context, m *model.User, id uuid.UUID) error {
	args := pgx.NamedArgs{
		"id":       id,
		"username": m.Username,
		"password": m.Password,
		"nickname": m.Nickname,
		"email":    m.Email,
		"phone":    m.Phone,
		"status":   m.Status,
	}
	sql := `INSERT INTO users
    (id, username, password, nickname, email, phone, status) VALUES
    (@id, @username, @password, @nickname, @email, @phone, @status)`
	_, err := r.db.Exec(ctx, sql, args)
	if err != nil {
		return err
	}
	return nil
}

func (r *userRepo) Update(ctx context.Context, m *model.User, id uuid.UUID) error {
	args := pgx.NamedArgs{
		"id":       id,
		"username": m.Username,
		"nickname": m.Nickname,
		"email":    m.Email,
		"phone":    m.Phone,
		"status":   m.Status,
	}
	sql := `UPDATE users SET
	username=@username, nickname=@nickname, email=@email, phone=@phone, status=@status, updated_at=CURRENT_TIMESTAMP
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

func (r *userRepo) UpdatePassword(ctx context.Context, id uuid.UUID, old string, new string) error {
	args := pgx.NamedArgs{
		"id":  id,
		"old": old,
		"new": new,
	}
	sql := `UPDATE users SET password=@new, updated_at=CURRENT_TIMESTAMP
    WHERE deleted_at IS NULL AND id=@id AND password=@old`
	tag, err := r.db.Exec(ctx, sql, args)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (r *userRepo) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	sql := `SELECT id, username, password, nickname, email, phone, status, created_at, updated_at FROM users
    WHERE deleted_at IS NULL AND id=$1`
	var m model.User
	err := r.db.QueryRow(ctx, sql, id).Scan(
		&m.ID, &m.Username, &m.Password, &m.Nickname, &m.Email, &m.Phone, &m.Status, &m.CreatedAt, &m.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *userRepo) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	sql := `SELECT id, username, password, nickname, email, phone, status, created_at, updated_at FROM users
    WHERE deleted_at IS NULL AND username=$1`
	var m model.User
	err := r.db.QueryRow(ctx, sql, username).Scan(
		&m.ID, &m.Username, &m.Password, &m.Nickname, &m.Email, &m.Phone, &m.Status, &m.CreatedAt, &m.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *userRepo) Count(ctx context.Context, data *request.UserQuery) (int64, error) {
	var sql strings.Builder
	args := pgx.NamedArgs{}
	sql.WriteString(`SELECT COUNT(*) FROM users WHERE deleted_at IS NULL`)
	sql.WriteString(buildUserQuerySQL(data, args))

	var total int64
	err := r.db.QueryRow(ctx, sql.String(), args).Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (r *userRepo) List(ctx context.Context, data *request.UserQuery) ([]*model.User, error) {
	var sql strings.Builder
	args := pgx.NamedArgs{}
	sql.WriteString(`SELECT id, username, password, nickname, email, phone, status, created_at, updated_at
    FROM users WHERE deleted_at IS NULL`)
	sql.WriteString(buildUserQuerySQL(data, args))
	sql.WriteString(` ORDER BY created_at DESC`)
	if data.PageSize > 0 && data.CurrentPage > 0 {
		sql.WriteString(` LIMIT @limit OFFSET @offset`)
		args["limit"] = data.PageSize
		args["offset"] = (data.CurrentPage - 1) * data.PageSize
	}

	var list []*model.User
	rows, err := r.db.Query(ctx, sql.String(), args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var m model.User
		err = rows.Scan(&m.ID, &m.Username, &m.Password, &m.Nickname, &m.Email, &m.Phone, &m.Status, &m.CreatedAt, &m.UpdatedAt)
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

func (r *userRepo) Delete(ctx context.Context, id uuid.UUID) error {
	sql := `UPDATE users SET deleted_at=CURRENT_TIMESTAMP WHERE deleted_at IS NULL AND id=$1`
	tag, err := r.db.Exec(ctx, sql, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func buildUserQuerySQL(req *request.UserQuery, args pgx.NamedArgs) string {
	var sql strings.Builder
	if req.Username != "" {
		sql.WriteString(` AND username=@username`)
		args["username"] = req.Username
	}
	if req.Nickname != "" {
		sql.WriteString(` AND nickname LIKE '%@nickname%'`)
		args["nickname"] = req.Nickname
	}
	if req.Email != "" {
		sql.WriteString(` AND email=@email`)
		args["email"] = req.Email
	}
	if req.Phone != "" {
		sql.WriteString(` AND phone=@phone`)
		args["phone"] = req.Phone
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
