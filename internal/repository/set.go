package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Set struct {
	Transactor      Transaction
	UserRepo        User
	RoleRepo        Role
	MenuRepo        Menu
	UserRoleRepo    UserRole
	RoleMenuRepo    RoleMenu
	LoginLogRepo    LoginLog
	UserSessionRepo UserSession
	DictRepo        Dict
}

func NewRepoSet(db *pgxpool.Pool) *Set {
	return &Set{
		Transactor:      NewTransactor(db),
		UserRepo:        NewUserRepo(db),
		RoleRepo:        NewRoleRepo(db),
		MenuRepo:        NewMenuRepo(db),
		UserRoleRepo:    NewUserRoleRepo(db),
		RoleMenuRepo:    NewRoleMenuRepo(db),
		LoginLogRepo:    NewLoginLogRepo(db),
		UserSessionRepo: NewUserSessionRepo(db),
		DictRepo:        NewDictRepo(db),
	}
}
