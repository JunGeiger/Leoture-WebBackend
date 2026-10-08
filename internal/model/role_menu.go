package model

import (
	"github.com/google/uuid"
	"time"
)

// RoleMenu 角色菜单表
// Table role_menus
type RoleMenu struct {
	RoleID    uuid.UUID `json:"role_id" db:"role_id"`
	MenuID    uuid.UUID `json:"menu_id" db:"menu_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

func (RoleMenu) TableName() string { return "role_menu" }
