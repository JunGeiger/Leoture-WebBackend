package model

import (
	"github.com/google/uuid"
	"time"
)

// LoginLog 用户登录日志表
// Table login_log
type LoginLog struct {
	ID        uuid.UUID `json:"id" db:"id"` // 自动递增ID字段
	Username  string    `json:"username" db:"username"`
	IPAddress string    `json:"ip_address" db:"ip_address"`
	UserAgent string    `json:"user_agent" db:"user_agent"`
	OriginUrl string    `json:"origin_url" db:"origin_url"` // 登录请求来源页面URL
	Location  string    `json:"location" db:"location"`     // 用户自报地点信息(可选)
	Result    string    `json:"result" db:"result"`         // 登录返回结果，或者内部错误详情
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// TableName 返回表名，方便后续 ORM/QueryBuilder 使用
func (LoginLog) TableName() string {
	return "login_log"
}
