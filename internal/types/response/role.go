package response

import "LeotureWeb/internal/model"

// RoleResp 角色响应数据结构
type RoleResp struct {
	ID        string `json:"id"`         // 角色ID
	Name      string `json:"name"`       // 角色名称
	Code      string `json:"code"`       // 角色编码
	Remark    string `json:"remark"`     // 角色说明
	Status    int16  `json:"status"`     // 数据状态：1:正常, 2:禁用
	CreatedAt string `json:"created_at"` // 创建时间
	UpdatedAt string `json:"updated_at"` // 更新时间
}

func ToRoleResp(m *model.Role) *RoleResp {
	return &RoleResp{
		ID:        m.ID.String(),
		Name:      m.Name,
		Code:      m.Code,
		Remark:    m.Remark,
		Status:    int16(m.Status),
		CreatedAt: m.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: m.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
