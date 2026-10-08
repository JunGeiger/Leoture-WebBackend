package user_role

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
	Update(ctx context.Context, data *request.UserRoles) error
	GetRolesByUserID(ctx context.Context, userID string) ([]*response.RoleResp, error)
}

type userRoleService struct {
	transactor repository.Transaction
	casbinSvc  *authorization.CasbinService
	repo       repository.UserRole
	userRepo   repository.User
	roleRepo   repository.Role
}

func NewService(transactor repository.Transaction, casbinSvc *authorization.CasbinService,
	repo repository.UserRole, userRepo repository.User, roleRepo repository.Role) Service {
	return &userRoleService{
		transactor: transactor,
		casbinSvc:  casbinSvc,
		repo:       repo,
		userRepo:   userRepo,
		roleRepo:   roleRepo,
	}
}

func (s userRoleService) Update(ctx context.Context, data *request.UserRoles) error {
	err := s.transactor.WithTransaction(ctx, func(tx pgx.Tx) error {
		userID, err := uuid.Parse(data.UserID)
		if err != nil {
			return err
		}
		user, err := s.userRepo.GetByID(ctx, userID)
		if err != nil {
			return err
		}
		if user == nil {
			return nil
		}
		if err := s.repo.Delete(ctx, userID); err != nil {
			return err
		}
		if len(data.RoleIDs) == 0 {
			return nil
		}
		roleIDs := make([]uuid.UUID, len(data.RoleIDs))
		for _, roleID := range data.RoleIDs {
			roleUUID, err := uuid.Parse(roleID)
			if err != nil {
				return err
			}
			roleIDs = append(roleIDs, roleUUID)
		}
		roles, err := s.roleRepo.GetByIDs(ctx, roleIDs)
		if err != nil {
			return err
		}
		if len(roles) == 0 {
			return nil
		}
		models := make([]*model.UserRole, 0, len(roles))
		roleCodes := make([]string, 0, len(roles))
		for _, role := range roles {
			var userRole *model.UserRole
			userRole.UserID = userID
			userRole.RoleID = role.ID
			models = append(models, userRole)
			roleCodes = append(roleCodes, role.Code)
		}
		if err := s.repo.Create(ctx, models); err != nil {
			return err
		}

		if err := s.UpdateCasbinGroupingPolicy(user.Username, roleCodes); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (s userRoleService) GetRolesByUserID(ctx context.Context, userID string) ([]*response.RoleResp, error) {
	id, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}
	list, err := s.repo.GetRolesByUserId(ctx, id)
	if err != nil {
		return nil, err
	}
	roles := make([]*response.RoleResp, 0, len(list))
	for _, role := range list {
		roles = append(roles, response.ToRoleResp(role))
	}
	return roles, nil
}

func (s userRoleService) UpdateCasbinGroupingPolicy(username string, roleCodes []string) error {
	if _, err := s.casbinSvc.Enforcer.RemoveFilteredNamedPolicy("g", 0, username); err != nil {
		return err
	}
	if len(roleCodes) == 0 {
		return nil
	}
	rules := make([][]string, 0, len(roleCodes))
	for _, roleCode := range roleCodes {
		var rule []string
		rule = append(rule, username)
		rule = append(rule, roleCode)
		rules = append(rules, rule)
	}

	if _, err := s.casbinSvc.Enforcer.AddGroupingPolicies(rules); err != nil {
		return err
	}
	return nil
}
