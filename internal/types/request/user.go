package request

import (
	"LeotureWeb/internal/model"
)

// UserCreate 创建用户请求数据结构
type UserCreate struct {
	Username string `json:"username" binding:"required"`         // 用户名称
	Password string `json:"password" binding:"required"`         // 用户密码
	Nickname string `json:"nickname" binding:"required"`         // 用户昵称
	Email    string `json:"email" binding:"required,email"`      // 用户邮箱
	Phone    string `json:"phone" binding:"required"`            // 用户手机号码
	Status   int16  `json:"status" binding:"required,oneof=1 2"` // 数据状态：1:正常, 2:禁用
}

// ToModelWithPassword 转换为数据库映射模型，并补充密码字段
func (r *UserCreate) ToModelWithPassword(password string) *model.User {
	m := &model.User{
		Username: r.Username,
		Password: password,
		Nickname: r.Nickname,
		Email:    r.Email,
		Phone:    r.Phone,
	}
	m.Status.FromInt(int(r.Status))
	return m
}

// UserUpdate 更新用户请求数据结构
type UserUpdate struct {
	Username string `json:"username" binding:"-"`         // 用户名称
	Nickname string `json:"nickname" binding:"-"`         // 用户昵称
	Email    string `json:"email" binding:"-,email"`      // 用户邮箱
	Phone    string `json:"phone" binding:"-"`            // 用户手机号码
	Status   int16  `json:"status" binding:"-,oneof=1 2"` // 数据状态：1:正常, 2:禁用
	DataID
}

// ToModel 转换为数据库映射模型
func (r *UserUpdate) ToModel() *model.User {
	m := &model.User{
		Username: r.Username,
		Nickname: r.Nickname,
		Email:    r.Email,
		Phone:    r.Phone,
	}
	m.Status.FromInt(int(r.Status))
	return m
}

// UserQuery 用户查询请求数据结构
type UserQuery struct {
	Username string `json:"username" form:"username" binding:"-"`               // 用户名称
	Nickname string `json:"nickname" form:"nickname" binding:"-"`               // 用户昵称
	Email    string `json:"email" form:"email" binding:"-"`                     // 邮箱地址
	Phone    string `json:"phone" form:"phone" binding:"-"`                     // 手机号码
	Status   int16  `json:"status" form:"status" binding:"omitempty,oneof=1 2"` // 数据状态
	PageParams
	DateRange
}
