package role_menu

import (
	"LeotureWeb/internal/authorization"
	"LeotureWeb/internal/model"
	"LeotureWeb/internal/repository"
	"LeotureWeb/internal/types/request"
	"LeotureWeb/internal/types/response"
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Service interface {
	Update(ctx context.Context, data *request.RoleMenus) error
	GetMenusByRoleID(ctx context.Context, roleID string) ([]*response.MenuTreeResp, error)
	UpdateCasbinPolicies(code string, paths []string) error
}

type roleMenuService struct {
	transactor repository.Transaction
	casbinSvc  *authorization.CasbinService
	repo       repository.RoleMenu
	roleRepo   repository.Role
	menuRepo   repository.Menu
}

func NewService(transactor repository.Transaction, casbinSvc *authorization.CasbinService,
	repo repository.RoleMenu, roleRepo repository.Role, menuRepo repository.Menu) Service {
	return &roleMenuService{
		transactor: transactor,
		casbinSvc:  casbinSvc,
		repo:       repo,
		roleRepo:   roleRepo,
		menuRepo:   menuRepo,
	}
}

func (s roleMenuService) Update(ctx context.Context, data *request.RoleMenus) error {
	err := s.transactor.WithTransaction(ctx, func(tx pgx.Tx) error {
		roleID, err := uuid.Parse(data.RoleID)
		if err != nil {
			return err
		}
		role, err := s.roleRepo.GetByID(ctx, roleID)
		if err != nil {
			return err
		}
		if role == nil {
			return nil
		}
		if err := s.repo.Delete(ctx, roleID); err != nil {
			return err
		}
		if len(data.MenuIDs) == 0 {
			return nil
		}
		menuIDs := make([]uuid.UUID, len(data.MenuIDs))
		for _, menuID := range data.MenuIDs {
			menuUUID, err := uuid.Parse(menuID)
			if err != nil {
				return err
			}
			menuIDs = append(menuIDs, menuUUID)
		}
		menus, err := s.menuRepo.GetByIDs(ctx, menuIDs)
		if err != nil {
			return err
		}
		if len(menus) == 0 {
			return nil
		}
		models := make([]*model.RoleMenu, 0, len(menus))
		menuPaths := make([]string, 0, len(menus))
		for _, menu := range menus {
			var roleMenu *model.RoleMenu
			roleMenu.RoleID = roleID
			roleMenu.MenuID = menu.ID
			models = append(models, roleMenu)
			menuPaths = append(menuPaths, menu.Path)
		}
		if err := s.repo.Create(ctx, models); err != nil {
			return err
		}
		if err := s.UpdateCasbinPolicies(role.Code, menuPaths); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (s roleMenuService) GetMenusByRoleID(ctx context.Context, roleID string) ([]*response.MenuTreeResp, error) {
	id, err := uuid.Parse(roleID)
	if err != nil {
		return nil, err
	}
	menus, err := s.repo.GetMenusBYRoleID(ctx, id)
	if err != nil {
		return nil, err
	}
	if len(menus) == 0 {
		return nil, nil
	}
	menuMap := response.ToMenuTreeResp(menus)
	return menuMap, nil
}

func (s roleMenuService) UpdateCasbinPolicies(roleCode string, menuPaths []string) error {
	if _, err := s.casbinSvc.Enforcer.RemoveFilteredNamedPolicy("p", 0, roleCode); err != nil {
		return err
	}
	if len(menuPaths) == 0 {
		return nil
	}
	rules := make([][]string, 0, len(menuPaths))
	for _, menuPath := range menuPaths {
		var rule []string
		rule = append(rule, roleCode)
		rule = append(rule, menuPath)
		rule = append(rule, ".*")
		rule = append(rule, authorization.EftAllow)
		rules = append(rules, rule)
	}

	if _, err := s.casbinSvc.Enforcer.AddPolicies(rules); err != nil {
		return err
	}
	return nil
}
