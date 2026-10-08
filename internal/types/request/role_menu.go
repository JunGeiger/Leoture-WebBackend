package request

// RoleMenus 更新角色与菜单关联权限请求数据结构
type RoleMenus struct {
	RoleID  string   `json:"roleId" form:"roleId" binding:"required"`   // 角色ID
	MenuIDs []string `json:"menuIDs" form:"menuIDs" binding:"required"` // 多个菜单ID
}
