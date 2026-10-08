package model

import (
	"time"
)

// BaseModel 包含所有业务表的公共字段
type BaseModel struct {
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"-" db:"deleted_at"`
}

// Status 通用状态枚举
type Status int16

const (
	StatusEnabled  Status = 1 // 正常/显示
	StatusDisabled Status = 2 // 禁用/隐藏
)

// IsZero 实现零值判断，便于校验
func (s *Status) IsZero() bool { return *s == 0 }

func (s *Status) FromInt(i int) {
	switch i {
	case int(StatusEnabled), int(StatusDisabled):
		*s = Status(i)
	default:
		*s = Status(2)
	}
}
