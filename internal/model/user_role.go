package model

import (
	"github.com/google/uuid"
	"time"
)

// UserRole 用户角色表
// Table user_roles
type UserRole struct {
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	RoleID    uuid.UUID `json:"role_id" db:"role_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

func (UserRole) TableName() string { return "user_role" }
