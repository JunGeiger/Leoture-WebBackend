package model

import (
	"github.com/google/uuid"
	"time"
)

// UserSession 用户会话Session表
// Table user_session
type UserSession struct {
	ID         uuid.UUID `json:"id" db:"id"`
	UserID     uuid.UUID `json:"UserID" db:"user_id"`
	RefreshJTI uuid.UUID `json:"refreshJTI" db:"refresh_jti"`  // Refresh Token唯一标识
	LoginLogID uuid.UUID `json:"loginLogID" db:"login_log_id"` // 关联登录日志ID
	Status     int16     `json:"status" db:"status"`           // 1: valid 正常可用, 2: revoked 主动注销, 3: kicked 被新登录顶掉, 4: locked 风控 / 封禁
	CreatedAt  time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt  time.Time `json:"updatedAt" db:"updated_at"`
}

// TableName 返回表名，方便后续 ORM/QueryBuilder 使用
func (UserSession) TableName() string {
	return "user_session"
}
