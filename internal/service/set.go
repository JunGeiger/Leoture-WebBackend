package service

import (
	"LeotureWeb/internal/authorization"
	"LeotureWeb/internal/repository"
	"LeotureWeb/internal/service/auth"
	"LeotureWeb/internal/service/dict"
	"LeotureWeb/internal/service/menu"
	"LeotureWeb/internal/service/role"
	"LeotureWeb/internal/service/role_menu"
	"LeotureWeb/internal/service/user"
	"LeotureWeb/internal/service/user_role"
)

type Set struct {
	UserService     user.Service
	RoleService     role.Service
	MenuService     menu.Service
	UserRoleService user_role.Service
	RoleMenuService role_menu.Service
	AuthService     auth.Service
	DictService     dict.Service
}

func NewSet(repoSet *repository.Set, jwtSvc *authorization.JWTService, casbinSvc *authorization.CasbinService) *Set {
	return &Set{
		UserService:     user.NewService(repoSet.Transactor, repoSet.UserRepo),
		RoleService:     role.NewService(repoSet.Transactor, repoSet.RoleRepo),
		MenuService:     menu.NewService(repoSet.Transactor, repoSet.MenuRepo),
		UserRoleService: user_role.NewService(repoSet.Transactor, casbinSvc, repoSet.UserRoleRepo, repoSet.UserRepo, repoSet.RoleRepo),
		RoleMenuService: role_menu.NewService(repoSet.Transactor, casbinSvc, repoSet.RoleMenuRepo, repoSet.RoleRepo, repoSet.MenuRepo),
		AuthService:     auth.NewAuthService(repoSet.Transactor, jwtSvc, repoSet.UserRepo, repoSet.UserRoleRepo, repoSet.RoleMenuRepo, repoSet.LoginLogRepo, repoSet.UserSessionRepo),
		DictService:     dict.NewService(repoSet.Transactor, repoSet.DictRepo),
	}
}
