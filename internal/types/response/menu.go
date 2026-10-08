package response

import (
	"LeotureWeb/internal/model"
	"time"

	"github.com/google/uuid"
)

// MenuResp 菜单响应数据结构
type MenuResp struct {
	ID          uuid.UUID `json:"id"`       // 菜单ID
	ParentID    uuid.UUID `json:"parentId"` // 父菜单ID，为NULL则说明是根节点`
	Type        int16     `json:"type"`     // 类型：1-目录, 2-菜单, 3-按钮
	Name        string    `json:"name"`
	Title       string    `json:"title"`
	Path        string    `json:"path"`
	Component   string    `json:"component"`
	Icon        string    `json:"localIcon"`
	KeepAlive   bool      `json:"keepAlive"`   // 是否开启页面缓存
	OrderNo     int       `json:"orderNo"`     // 菜单排序
	ExternalUrl string    `json:"externalUrl"` // 外部链接
	IsHidden    bool      `json:"hideInMenu"`  // 是否隐藏菜单
	PermCode    string    `json:"permCode"`    // 权限标识
	ApiMethod   string    `json:"apiMethod"`   // 后端Api HTTP方法
	ApiPath     string    `json:"apiPath"`     // 后端Api地址
	Status      int16     `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// ToMenuResp 转换数据库映射模型为响应数据结构
func ToMenuResp(m *model.Menu) *MenuResp {
	return &MenuResp{
		ID:          m.ID,
		ParentID:    m.ParentID,
		Type:        m.Type,
		Name:        m.Name,
		Title:       m.Title,
		Path:        m.Path,
		Component:   m.Component,
		Icon:        m.Icon,
		KeepAlive:   m.KeepAlive,
		OrderNo:     m.OrderNo,
		ExternalUrl: m.ExternalUrl,
		IsHidden:    m.IsHidden,
		PermCode:    m.PermCode,
		Status:      int16(m.Status),
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

// MenuTreeResp 菜单树形响应数据结构
type MenuTreeResp struct {
	MenuResp
	Children []*MenuTreeResp `json:"children,omitempty"`
}

// ToMenuTreeResp 模型转树形响应
func ToMenuTreeResp(ms []*model.Menu) []*MenuTreeResp {
	menuMap := make(map[uuid.UUID]*MenuTreeResp)
	var roots []*MenuTreeResp

	// 第一遍：创建所有节点
	for _, m := range ms {
		node := &MenuTreeResp{
			MenuResp: *ToMenuResp(m),
			Children: make([]*MenuTreeResp, 0),
		}
		menuMap[node.ID] = node
	}
	// 第二遍：建立父子关系
	for i := range ms {
		node := menuMap[ms[i].ID]
		parentID := ms[i].ParentID
		if parentID == uuid.Nil {
			// 根节点
			roots = append(roots, node)
		} else if parent, exists := menuMap[parentID]; exists {
			// 有父节点，添加到children
			parent.Children = append(parent.Children, node)
		} else {
			// 父节点不存在（可能被删除），作为根节点处理
			roots = append(roots, node)
		}
	}

	return roots
}
