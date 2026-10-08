package role

import (
	"LeotureWeb/internal/repository"
	"LeotureWeb/internal/types/request"
	"LeotureWeb/internal/types/response"
	"LeotureWeb/internal/utils"
	"context"

	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, data *request.RoleCreate) error
	Update(ctx context.Context, data *request.RoleUpdate) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*response.RoleResp, error)
	GetByCode(ctx context.Context, code string) (*response.RoleResp, error)
	List(ctx context.Context, data *request.RoleQuery) ([]*response.RoleResp, error)
	ListByPage(ctx context.Context, data *request.RoleQuery) (*response.Paginated[response.RoleResp], error)
}

type roleService struct {
	transactor repository.Transaction
	repo       repository.Role
}

func NewService(transactor repository.Transaction, repo repository.Role) Service {
	return &roleService{
		transactor: transactor,
		repo:       repo,
	}
}

func (s *roleService) Create(ctx context.Context, data *request.RoleCreate) error {
	m := data.ToCreateModel()
	return s.repo.Create(ctx, m, utils.GenUUIDv7())
}

func (s *roleService) Update(ctx context.Context, data *request.RoleUpdate) error {
	m := data.ToUpdateModel()
	id, err := uuid.Parse(data.ID)
	if err != nil {
		return err
	}
	return s.repo.Update(ctx, m, id)
}

func (s *roleService) Delete(ctx context.Context, id string) error {
	idn, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(ctx, idn)
}

func (s *roleService) GetByID(ctx context.Context, id string) (*response.RoleResp, error) {
	idn, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	role, err := s.repo.GetByID(ctx, idn)
	if err != nil {
		return nil, err
	}
	return response.ToRoleResp(role), nil
}

func (s *roleService) GetByCode(ctx context.Context, code string) (*response.RoleResp, error) {
	role, err := s.repo.GetByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	return response.ToRoleResp(role), nil
}

func (s *roleService) List(ctx context.Context, data *request.RoleQuery) ([]*response.RoleResp, error) {
	result, err := s.repo.List(ctx, data)
	if err != nil {
		return nil, err
	}
	resp := make([]*response.RoleResp, 0, len(result))
	for _, role := range result {
		resp = append(resp, response.ToRoleResp(role))
	}
	return resp, nil
}

func (s *roleService) ListByPage(ctx context.Context, data *request.RoleQuery) (*response.Paginated[response.RoleResp], error) {
	count, err := s.repo.Count(ctx, data)
	if err != nil {
		return nil, err
	}
	data.InitPagination()
	var result = response.PagingData[response.RoleResp](nil, count, data.CurrentPage, data.PageSize)
	if count == 0 {
		return result, nil
	}
	list, err := s.repo.List(ctx, data)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return result, nil
	}
	roles := make([]*response.RoleResp, 0, len(list))
	for _, role := range list {
		roles = append(roles, response.ToRoleResp(role))
	}
	result.List = roles
	return result, nil
}
