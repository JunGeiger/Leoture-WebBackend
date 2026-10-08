package model

import (
	"github.com/google/uuid"
	"time"
)

// Dict 系统字典表
// Table dictionaries
type Dict struct {
	ID        uuid.UUID `json:"id"  db:"id"`
	Name      string    `json:"name"  db:"name"`
	Category  string    `json:"category"  db:"category"` // 字典类别
	Key       string    `json:"key" db:"key"`            // 字典key
	Val       string    `json:"val"  db:"val"`           // 字典value
	Sort      int       `json:"sort"  db:"sort"`
	IsDefault bool      `json:"is_default"  db:"is_default"` // 是否选项默认值
	Status    Status    `json:"status" db:"status"`          // 状态，1:正常, 2:禁用
	Remark    string    `json:"remark"  db:"remark"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// TableName 返回表名，方便后续 ORM/QueryBuilder 使用
func (Dict) TableName() string { return "dictionaries" }
