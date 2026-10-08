package model

import (
	"github.com/google/uuid"
)

// Menu 菜单表
// Table menus
type Menu struct {
	ID          uuid.UUID `json:"id" db:"id"`
	ParentID    uuid.UUID `json:"parent_id" db:"parent_id"` // 父菜单ID，为0则说明是根节点
	Type        int16     `json:"type" db:"type"`           // 类型：1-目录, 2-菜单, 3-按钮, 4-路由
	Name        string    `json:"name" db:"name"`
	Title       string    `json:"title" db:"title"`
	Path        string    `json:"path" db:"path"`
	Component   string    `json:"component" db:"component"`
	Icon        string    `json:"icon" db:"icon"`
	KeepAlive   bool      `json:"keepAlive" db:"keep_alive"`      // 是否开启页面缓存
	OrderNo     int       `json:"order_no" db:"order_no"`         // 菜单排序
	ExternalUrl string    `json:"external_url" db:"external_url"` // 外部链接
	IsHidden    bool      `json:"is_hidden" db:"is_hidden"`       // 是否隐藏菜单
	PermCode    string    `json:"perm_code" db:"perm_code"`       // 权限标识
	ApiMethod   string    `json:"api_method" db:"api_method"`     // 后端Api HTTP方法
	ApiPath     string    `json:"api_path" db:"api_path"`         // 后端Api地址
	Status      Status    `json:"status" db:"status"`
	BaseModel
}

func (Menu) TableName() string { return "menus" }
