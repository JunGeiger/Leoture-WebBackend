package request

import (
	"LeotureWeb/internal/model"
)

// RoleCreate 创建角色请求数据结构
type RoleCreate struct {
	Name   string `json:"name" binding:"required"`             // 角色名称
	Code   string `json:"code" binding:"required"`             // 角色编码
	Remark string `json:"remark" binding:"-"`                  // 角色说明
	Status int16  `json:"status" binding:"required,oneof=1 2"` // 数据状态：1:正常, 2:禁用
}

// ToModel 转换为数据库映射模型
func (r *RoleCreate) ToCreateModel() *model.Role {
	m := &model.Role{
		Name:   r.Name,
		Code:   r.Code,
		Remark: r.Remark,
	}
	m.Status.FromInt(int(r.Status))
	return m
}

// RoleUpdate 更新角色请求数据结构
type RoleUpdate struct {
	Name   string `json:"name" binding:"-"`             // 角色名称
	Code   string `json:"code" binding:"-"`             // 角色编码
	Remark string `json:"remark" binding:"-"`           // 角色说明
	Status int16  `json:"status" binding:"-,oneof=1 2"` // 数据状态：1:正常, 2:禁用
	DataID
}

// ToModel 转换为数据库映射模型
func (r *RoleUpdate) ToUpdateModel() *model.Role {
	m := &model.Role{
		Name:   r.Name,
		Code:   r.Code,
		Remark: r.Remark,
	}
	m.Status.FromInt(int(r.Status))
	return m
}

// RoleQuery 角色查询请求数据结构
type RoleQuery struct {
	Name   string `form:"name" binding:"-"`             // 角色名称
	Code   string `form:"code" binding:"-"`             // 角色编码
	Remark string `form:"remark" binding:"-"`           // 角色说明
	Status int16  `form:"status" binding:"-,oneof=1 2"` // 数据状态：1:正常, 2:禁用
	DateRange
	PageParams
}
