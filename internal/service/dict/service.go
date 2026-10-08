package dict

import (
	"LeotureWeb/internal/repository"
	"LeotureWeb/internal/types/request"
	"LeotureWeb/internal/types/response"
	"LeotureWeb/internal/utils"
	"context"

	"github.com/google/uuid"
)

// Service 字典业务接口
type Service interface {
	Create(ctx context.Context, data *request.DictCreate) error
	Update(ctx context.Context, data *request.DictUpdate) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*response.DictResp, error)
	GetByKey(ctx context.Context, category, key string) (*response.DictResp, error)
	List(ctx context.Context, data *request.DictQuery) ([]*response.DictResp, error)
	ListByPage(ctx context.Context, data *request.DictQuery) (*response.Paginated[response.DictResp], error)
}

type dictService struct {
	transactor repository.Transaction
	repo       repository.Dict
}

func NewService(transactor repository.Transaction, repo repository.Dict) Service {
	return &dictService{
		transactor: transactor,
		repo:       repo,
	}
}

func (s *dictService) Create(ctx context.Context, data *request.DictCreate) error {
	m := data.ToModel()
	return s.repo.Create(ctx, m, utils.GenUUIDv7())
}

func (s *dictService) Update(ctx context.Context, data *request.DictUpdate) error {
	m := data.ToModel()
	id, err := uuid.Parse(data.ID)
	if err != nil {
		return err
	}
	return s.repo.Update(ctx, m, id)
}

func (s *dictService) Delete(ctx context.Context, id string) error {
	idn, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(ctx, idn)
}

func (s *dictService) GetByID(ctx context.Context, id string) (*response.DictResp, error) {
	idn, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	dict, err := s.repo.GetByID(ctx, idn)
	if err != nil {
		return nil, err
	}
	return response.ToDictResp(dict), nil
}

func (s *dictService) GetByKey(ctx context.Context, category, key string) (*response.DictResp, error) {
	dict, err := s.repo.GetByKey(ctx, category, key)
	if err != nil {
		return nil, err
	}
	return response.ToDictResp(dict), nil
}

func (s *dictService) List(ctx context.Context, data *request.DictQuery) ([]*response.DictResp, error) {
	result, err := s.repo.List(ctx, data)
	if err != nil {
		return nil, err
	}
	resp := make([]*response.DictResp, 0, len(result))
	for _, dict := range result {
		resp = append(resp, response.ToDictResp(dict))
	}
	return resp, nil
}

func (s *dictService) ListByPage(ctx context.Context, data *request.DictQuery) (*response.Paginated[response.DictResp], error) {
	count, err := s.repo.Count(ctx, data)
	if err != nil {
		return nil, err
	}
	data.InitPagination()
	result := response.PagingData[response.DictResp](nil, count, data.CurrentPage, data.PageSize)
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
	dicts := make([]*response.DictResp, 0, len(list))
	for _, dict := range list {
		dicts = append(dicts, response.ToDictResp(dict))
	}
	result.List = dicts
	return result, nil
}
