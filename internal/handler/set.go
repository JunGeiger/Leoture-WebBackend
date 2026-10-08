package handler

import (
	"LeotureWeb/internal/service"
)

type Set struct {
	UserHandler     *User
	RoleHandler     *Role
	MenuHandler     *Menu
	UserRoleHandler *UserRole
	RoleMenuHandler *RoleMenu
	AuthHandler     *Auth
	DictHandler     *Dict
}

func NewSet(set *service.Set) *Set {
	return &Set{
		UserHandler:     newUserHandler(set.UserService),
		RoleHandler:     newRoleHandler(set.RoleService),
		MenuHandler:     newMenuHandler(set.MenuService),
		UserRoleHandler: newUserRoleHandler(set.UserRoleService),
		RoleMenuHandler: newRoleMenuHandler(set.RoleMenuService),
		AuthHandler:     newAuthHandler(set.AuthService),
		DictHandler:     newDictHandler(set.DictService),
	}
}
