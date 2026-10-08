package user

import (
	"LeotureWeb/internal/errors"
	"LeotureWeb/internal/repository"
	"LeotureWeb/internal/types/request"
	"LeotureWeb/internal/types/response"
	"LeotureWeb/internal/utils"
	"context"

	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, data *request.UserCreate) error
	Update(ctx context.Context, data *request.UserUpdate) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*response.UserResp, error)
	GetByUsername(ctx context.Context, username string) (*response.UserResp, error)
	List(ctx context.Context, data *request.UserQuery) ([]*response.UserResp, error)
	ListByPage(ctx context.Context, data *request.UserQuery) (*response.Paginated[response.UserResp], error)
}

type userService struct {
	transactor repository.Transaction
	repo       repository.User
}

func NewService(transactor repository.Transaction, repo repository.User) Service {
	return &userService{
		transactor: transactor,
		repo:       repo,
	}
}

func (s *userService) Create(ctx context.Context, data *request.UserCreate) error {
	m := data.ToModelWithPassword(utils.HashPassword(data.Password))
	if m.Password == "" {
		return errors.ErrInternal
	}
	return s.repo.Create(ctx, m, utils.GenUUIDv7())
}

func (s *userService) Update(ctx context.Context, data *request.UserUpdate) error {
	m := data.ToModel()
	id, err := uuid.Parse(data.ID)
	if err != nil {
		return err
	}
	return s.repo.Update(ctx, m, id)
}

func (s *userService) Delete(ctx context.Context, id string) error {
	idn, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(ctx, idn)
}

func (s *userService) GetByID(ctx context.Context, id string) (*response.UserResp, error) {
	idn, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	user, err := s.repo.GetByID(ctx, idn)
	if err != nil {
		return nil, err
	}
	resp := response.ToUserResp(user)
	return resp, nil
}

func (s *userService) GetByUsername(ctx context.Context, username string) (*response.UserResp, error) {
	user, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	resp := response.ToUserResp(user)
	return resp, nil
}

func (s *userService) List(ctx context.Context, data *request.UserQuery) ([]*response.UserResp, error) {
	result, err := s.repo.List(ctx, data)
	if err != nil {
		return nil, err
	}
	resp := make([]*response.UserResp, 0, len(result))
	for _, user := range result {
		resp = append(resp, response.ToUserResp(user))
	}
	return resp, nil
}

func (s *userService) ListByPage(ctx context.Context, data *request.UserQuery) (*response.Paginated[response.UserResp], error) {
	count, err := s.repo.Count(ctx, data)
	if err != nil {
		return nil, err
	}
	data.InitPagination()
	var result = response.PagingData[response.UserResp](nil, count, data.CurrentPage, data.PageSize)
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
	users := make([]*response.UserResp, 0, len(list))
	for _, user := range list {
		users = append(users, response.ToUserResp(user))
	}
	result.List = users
	return result, nil
}
