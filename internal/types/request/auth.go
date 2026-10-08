package request

import (
	"LeotureWeb/internal/model"
)

// LoginUser 用户登录请求数据结构
type LoginUser struct {
	Username string `json:"username" binding:"required"` // 用户名称
	Password string `json:"password" binding:"required"` // 用户密码
	Captcha  string `json:"captcha" binding:"omitempty"`
	LoginLog
}

// LoginLog 用户登录请求数据结构
type LoginLog struct {
	IPAddress string `json:"ipAddress" binding:"omitempty"`
	UserAgent string `json:"userAgent" binding:"omitempty"`
	OriginUrl string `json:"originUrl" binding:"omitempty"`
	Location  string `json:"location" binding:"omitempty"`
}

func (r LoginUser) ToLoginLogModel(result string) *model.LoginLog {
	return &model.LoginLog{
		Username:  r.Username,
		IPAddress: r.LoginLog.IPAddress,
		UserAgent: r.LoginLog.UserAgent,
		OriginUrl: r.LoginLog.OriginUrl,
		Location:  r.LoginLog.Location,
		Result:    result,
	}
}

type LogoutUser struct {
	UserID string `json:"userId" binding:"required"` // 用户ID
}

type RefreshToken struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}
