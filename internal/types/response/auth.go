package response

import (
	"LeotureWeb/internal/model"

	"github.com/google/uuid"
)

// LoginResp 登录成功响应
type LoginResp struct {
	// 访问令牌，建议同时通过 Authorization Header 返回
	AccessToken string `json:"token"`
	// 刷新令牌，用于无感续签，仅在 Body 中返回
	RefreshToken string `json:"refreshToken"`
}

// LoginRefreshResp 刷新Token成功响应
type LoginRefreshResp struct {
	// 访问令牌，建议同时通过 Authorization Header 返回
	AccessToken string `json:"access_token"`
	// 刷新令牌，用于无感续签，仅在 Body 中返回
	RefreshToken string `json:"refreshToken"`
}

type UserInfoResp struct {
	UserResp
	Buttons []string `json:"buttons"`
}

func ToUserInfoResp(m *model.User, buttons []string) *UserInfoResp {
	return &UserInfoResp{
		UserResp: *ToUserResp(m),
		Buttons:  buttons,
	}
}

type UserRoutesResp struct {
	Routes []*RouteResp `json:"routes"`
	Home   string       `json:"home"`
}

type RouteMeta struct {
	Title       string `json:"title"`
	Icon        string `json:"icon"`
	KeepAlive   bool   `json:"keepAlive"`
	OrderNo     int    `json:"order"`
	ExternalUrl string `json:"href"`
	IsHidden    bool   `json:"hideInMenu"`
}

type RouteResp struct {
	ID        uuid.UUID    `json:"-"`
	Name      string       `json:"name"`
	Path      string       `json:"path"`
	Component string       `json:"component"`
	Meta      RouteMeta    `json:"meta"`
	Children  []*RouteResp `json:"children"`
}

func ToRouteResp(ms []*model.Menu) []*RouteResp {
	menuMap := make(map[uuid.UUID]*RouteResp)
	var roots []*RouteResp

	// 第一遍：创建所有节点
	for _, m := range ms {
		node := &RouteResp{
			ID:        m.ID,
			Name:      m.Name,
			Path:      m.Path,
			Component: m.Component,
			Meta: RouteMeta{
				Title:       m.Title,
				Icon:        m.Icon,
				KeepAlive:   m.KeepAlive,
				OrderNo:     m.OrderNo,
				ExternalUrl: m.ExternalUrl,
				IsHidden:    m.IsHidden,
			},
			Children: make([]*RouteResp, 0),
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
