package menu

import (
	"LeotureWeb/internal/repository"
	"LeotureWeb/internal/types/request"
	"LeotureWeb/internal/types/response"
	"LeotureWeb/internal/utils"
	"context"

	"github.com/google/uuid"
)

// Service 菜单业务接口
type Service interface {
	Create(ctx context.Context, data *request.MenuCreate) error
	Update(ctx context.Context, data *request.MenuUpdate) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*response.MenuResp, error)
	List(ctx context.Context, data *request.MenuQuery) ([]*response.MenuResp, error)
	ListByPage(ctx context.Context, data *request.MenuQuery) (*response.Paginated[response.MenuResp], error)
	Tree(ctx context.Context) ([]*response.MenuTreeResp, error)
}

type menuService struct {
	transactor repository.Transaction
	repo       repository.Menu
}

func NewService(transactor repository.Transaction, repo repository.Menu) Service {
	return &menuService{
		transactor: transactor,
		repo:       repo,
	}
}

func (s *menuService) Create(ctx context.Context, data *request.MenuCreate) error {
	m := data.ToModel()
	return s.repo.Create(ctx, m, utils.GenUUIDv7())
}

func (s *menuService) Update(ctx context.Context, data *request.MenuUpdate) error {
	m := data.ToModel()
	id, err := uuid.Parse(data.ID)
	if err != nil {
		return err
	}
	return s.repo.Update(ctx, m, id)
}

func (s *menuService) Delete(ctx context.Context, id string) error {
	idn, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(ctx, idn)
}

func (s *menuService) GetByID(ctx context.Context, id string) (*response.MenuResp, error) {
	idn, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	menu, err := s.repo.GetByID(ctx, idn)
	if err != nil {
		return nil, err
	}
	return response.ToMenuResp(menu), nil
}

func (s *menuService) List(ctx context.Context, data *request.MenuQuery) ([]*response.MenuResp, error) {
	result, err := s.repo.List(ctx, data)
	if err != nil {
		return nil, err
	}
	resp := make([]*response.MenuResp, 0, len(result))
	for i := range result {
		resp = append(resp, response.ToMenuResp(result[i]))
	}
	return resp, nil
}

func (s *menuService) ListByPage(ctx context.Context, data *request.MenuQuery) (*response.Paginated[response.MenuResp], error) {
	count, err := s.repo.Count(ctx, data)
	if err != nil {
		return nil, err
	}
	data.InitPagination()
	result := response.PagingData[response.MenuResp](nil, count, data.CurrentPage, data.PageSize)
	if count == 0 {
		return result, nil
	}
	list, err := s.repo.List(ctx, data)
	if err != nil {
		return nil, err
	}
	menus := make([]*response.MenuResp, 0, len(list))
	for i := range list {
		menus = append(menus, response.ToMenuResp(list[i]))
	}
	result.List = menus
	return result, nil
}

func (s *menuService) Tree(ctx context.Context) ([]*response.MenuTreeResp, error) {
	allMenus, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	if len(allMenus) == 0 {
		return nil, nil
	}
	// 构建树形结构
	roots := response.ToMenuTreeResp(allMenus)
	return roots, nil
}
