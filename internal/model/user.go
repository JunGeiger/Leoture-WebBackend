package model

import (
	"github.com/google/uuid"
)

// User 用户表
// Table users
type User struct {
	ID       uuid.UUID `json:"id"  db:"id"`
	Username string    `json:"username"  db:"username"`
	Password string    `json:"-"  db:"password"`
	Nickname string    `json:"nickname" db:"nickname"`
	Email    string    `json:"email"  db:"email"`
	Phone    string    `json:"phone"  db:"phone"`
	Status   Status    `json:"status" db:"status"`
	BaseModel
}

// TableName 返回表名，方便后续 ORM/QueryBuilder 使用
func (User) TableName() string { return "users" }
