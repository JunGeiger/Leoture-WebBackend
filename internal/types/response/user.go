package response

import (
	"LeotureWeb/internal/model"
	"time"
)

// UserResp 用户响应数据结构
type UserResp struct {
	ID        string    `json:"userId"`    // 用户ID
	Username  string    `json:"userName"`  // 用户名称
	Nickname  string    `json:"nickName"`  // 用户昵称
	Email     string    `json:"email"`     // 邮箱地址
	Phone     string    `json:"phone"`     // 手机号码
	Status    int16     `json:"status"`    // 数据状态：1:正常, 2:禁用
	CreatedAt time.Time `json:"createdAt"` // 创建时间
	UpdatedAt time.Time `json:"updatedAt"` // 更新时间
}

// ToUserResp 将数据库模型转换为响应结构体
func ToUserResp(m *model.User) *UserResp {
	return &UserResp{
		ID:        m.ID.String(),
		Username:  m.Username,
		Nickname:  m.Nickname,
		Email:     m.Email,
		Phone:     m.Phone,
		Status:    int16(m.Status),
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}
