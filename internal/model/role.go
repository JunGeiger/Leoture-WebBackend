package model

import "github.com/google/uuid"

// Role 角色表
// Table roles
type Role struct {
	ID     uuid.UUID `json:"id" db:"id"`
	Name   string    `json:"name" db:"name"`
	Code   string    `json:"code" db:"code"`
	Remark string    `json:"remark" db:"remark"`
	Status Status    `json:"status" db:"status"`
	BaseModel
}

func (Role) TableName() string { return "roles" }
