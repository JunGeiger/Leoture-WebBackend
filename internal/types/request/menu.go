package request

import (
	"LeotureWeb/internal/model"
	"github.com/google/uuid"
)

// MenuCreate 创建菜单请求数据结构
type MenuCreate struct {
	ParentID    uuid.UUID `json:"parent_id" binding:"omitempty"`       // 父菜单ID，为NULL则说明是根节点
	Type        int16     `json:"type" binding:"required,oneof=1 2 3"` // 类型：1-目录, 2-菜单, 3-按钮
	Name        string    `json:"name" binding:"omitempty"`
	Title       string    `json:"title" binding:"omitempty"`
	KeepAlive   bool      `json:"keepAlive" binding:"omitempty"` // 是否开启页面缓存
	Icon        string    `json:"icon" binding:"omitempty"`
	OrderNo     int       `json:"order_no" binding:"required"`      // 菜单排序
	ExternalUrl string    `json:"external_url" binding:"omitempty"` // 外部链接
	IsHidden    bool      `json:"is_hidden" binding:"required"`     // 是否隐藏菜单
	PermCode    string    `json:"perm_code" binding:"required"`     // 权限标识
	ApiPath     string    `json:"api_path" binding:"required"`      // 后端Api地址
	Status      int16     `json:"status" binding:"required"`
}

// ToModel 转换为数据库映射模型
func (r *MenuCreate) ToModel() *model.Menu {
	m := &model.Menu{
		ParentID:    r.ParentID,
		Type:        r.Type,
		Name:        r.Name,
		Title:       r.Title,
		KeepAlive:   r.KeepAlive,
		Icon:        r.Icon,
		OrderNo:     r.OrderNo,
		ExternalUrl: r.ExternalUrl,
		IsHidden:    r.IsHidden,
		PermCode:    r.PermCode,
		ApiPath:     r.ApiPath,
	}
	m.Status.FromInt(int(r.Status))
	return m
}

// MenuUpdate 更新菜单请求数据结构
type MenuUpdate struct {
	DataID
	ParentID    uuid.UUID `json:"parent_id" binding:"omitempty"`       // 父菜单ID，为NULL则说明是根节点
	Type        int16     `json:"type" binding:"required,oneof=1 2 3"` // 类型：1-目录, 2-菜单, 3-按钮
	Name        string    `json:"name" binding:"omitempty"`
	Title       string    `json:"title" binding:"omitempty"`
	KeepAlive   bool      `json:"keepAlive" binding:"omitempty"` // 是否开启页面缓存
	Icon        string    `json:"icon" binding:"omitempty"`
	OrderNo     int       `json:"order_no" binding:"required"`      // 菜单排序
	ExternalUrl string    `json:"external_url" binding:"omitempty"` // 外部链接
	IsHidden    bool      `json:"is_hidden" binding:"required"`     // 是否隐藏菜单
	PermCode    string    `json:"perm_code" binding:"required"`     // 权限标识
	ApiPath     string    `json:"api_path" binding:"required"`      // 后端Api地址
	Status      int16     `json:"status" binding:"required"`
}

// ToModel 转换为数据库映射模型
func (r *MenuUpdate) ToModel() *model.Menu {
	m := &model.Menu{
		ParentID:  r.ParentID,
		Type:      r.Type,
		Name:      r.Name,
		Title:     r.Title,
		KeepAlive: r.KeepAlive,
	}
	m.Status.FromInt(int(r.Status))
	return m
}

// MenuQuery 菜单查询请求数据结构
type MenuQuery struct {
	Type     int16     `form:"type" binding:"omitempty,oneof=1 2 3 "` // 菜单类型
	Name     string    `form:"name" binding:"omitempty"`              // 菜单名称
	Title    string    `form:"title" binding:"omitempty"`
	Status   int16     `form:"status" binding:"omitempty,oneof=1 2"` // 数据状态
	ParentID uuid.UUID `form:"parentId" binding:"omitempty"`         // 父菜单ID
	PageParams
	DateRange
}
