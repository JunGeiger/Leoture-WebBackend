package auth

import (
	"LeotureWeb/internal/authorization"
	"LeotureWeb/internal/errors"
	"LeotureWeb/internal/model"
	"LeotureWeb/internal/repository"
	"LeotureWeb/internal/types/request"
	"LeotureWeb/internal/types/response"
	"LeotureWeb/internal/utils"
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Service interface {
	Login(ctx context.Context, data *request.LoginUser) (*response.LoginResp, error)
	RefreshToken(ctx context.Context, r *request.RefreshToken) (*response.LoginRefreshResp, error)
	GetUserInfo(ctx context.Context, id string) (*response.UserInfoResp, error)
	GetUserRoutes(ctx context.Context, id string) ([]*response.RouteResp, error)
	Logout(ctx context.Context, userID string) error
}

type authService struct {
	jwtSvc       *authorization.JWTService
	transactor   repository.Transaction
	userRepo     repository.User
	userRoleRepo repository.UserRole
	roleMenuRepo repository.RoleMenu
	logRepo      repository.LoginLog
	sessionRepo  repository.UserSession
}

func NewAuthService(transactor repository.Transaction, jwtSvc *authorization.JWTService, userRepo repository.User,
	userRoleRepo repository.UserRole, roleMenuRepo repository.RoleMenu, logRepo repository.LoginLog, sessionRepo repository.UserSession) Service {
	return authService{
		jwtSvc:       jwtSvc,
		transactor:   transactor,
		userRepo:     userRepo,
		userRoleRepo: userRoleRepo,
		roleMenuRepo: roleMenuRepo,
		logRepo:      logRepo,
		sessionRepo:  sessionRepo,
	}
}

func (s authService) Login(ctx context.Context, data *request.LoginUser) (*response.LoginResp, error) {
	// 登录日志ID
	loginLogID := utils.GenUUIDv7()
	resp, err := s.loginByUser(ctx, data, loginLogID)

	// 创建登录日志
	msg := "成功"
	if err != nil {
		msg = "失败"
	}
	s.createLoginLog(ctx, data, loginLogID, msg)

	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s authService) loginByUser(ctx context.Context, data *request.LoginUser, logID uuid.UUID) (*response.LoginResp, error) {
	// 获取用户数据
	user, err := s.userRepo.GetByUsername(ctx, data.Username)
	if err != nil {
		return nil, err
	}
	if user == nil || !utils.CheckPassword(data.Password, user.Password) {
		return nil, errors.New(errors.CodeForbidden, "用户名或密码错误")
	}
	// 判断用户状态是否正常
	if user.Status != model.StatusEnabled {
		return nil, errors.ErrForbidden.WithMsg("用户状态异常")
	}

	// 登录返回数据
	resp := &response.LoginResp{}

	// 用户SessionID
	sessionID := utils.GenUUIDv7()
	// 生成AccessToken
	resp.AccessToken, err = s.jwtSvc.GenerateAccessToken(user.ID.String(), sessionID.String(), user.Username)
	if err != nil {
		return nil, err
	}
	// RefreshTokenID
	var refreshJTI uuid.UUID
	// 生成RefreshToken
	resp.RefreshToken, refreshJTI, err = s.jwtSvc.GenerateRefreshToken(user.ID.String(), sessionID.String())
	if err != nil {
		return nil, err
	}

	err = s.transactor.WithTransaction(ctx, func(tx pgx.Tx) error {
		// 注销该用户名下其他有效RefreshToken状态
		if err := s.sessionRepo.UpdateExpiredByUserID(ctx, authorization.StatusKicked, user.ID); err != nil {
			return err
		}
		// 保存新的RefreshToken
		err := s.sessionRepo.Create(ctx, &model.UserSession{
			ID:         sessionID,
			UserID:     user.ID,
			RefreshJTI: refreshJTI,
			LoginLogID: logID,
			Status:     authorization.StatusValid})
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s authService) createLoginLog(ctx context.Context, data *request.LoginUser, logID uuid.UUID, msg string) {
	err := s.logRepo.Create(ctx, &model.LoginLog{
		ID:        logID,
		Username:  data.Username,
		IPAddress: data.IPAddress,
		UserAgent: data.UserAgent,
		OriginUrl: data.OriginUrl,
		Location:  data.Location,
		Result:    msg,
	})
	if err != nil {
		slog.Error("Login日志保存失败", slog.String("err", err.Error()))
	}
}

func (s authService) RefreshToken(ctx context.Context, r *request.RefreshToken) (*response.LoginRefreshResp, error) {
	// 解析Token数据
	claims, err := s.jwtSvc.ParseToken(r.RefreshToken)
	if err != nil {
		return nil, err
	}
	// 非RefreshToken直接返回错误
	if claims.Type != authorization.RefreshToken {
		slog.Error("RefreshToken, Token类型错误:  + claims.Type")
		return nil, errors.ErrInternal.WithMsg("Token类型错误: " + claims.Type)
	}
	// 解析UserSessionID
	sessionID, err := uuid.Parse(claims.SessionID)
	if err != nil {
		return nil, err
	}
	// 解析用户ID
	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return nil, err
	}
	// 获取用户保存的Session数据
	userSession, err := s.sessionRepo.GetByIDAndUserID(ctx, sessionID, userID)
	if err != nil {
		return nil, err
	}
	if userSession == nil || userSession.Status != authorization.StatusValid {
		return nil, errors.ErrUnauthenticated
	}
	// 获取用户数据
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user.Status != model.StatusEnabled {
		return nil, errors.New(errors.CodeForbidden, "用户状态异常")
	}

	resp := &response.LoginRefreshResp{
		RefreshToken: r.RefreshToken,
	}
	// 生成AccessToken
	resp.AccessToken, err = s.jwtSvc.GenerateAccessToken(user.ID.String(), sessionID.String(), user.Username)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s authService) GetUserInfo(ctx context.Context, id string) (*response.UserInfoResp, error) {
	idn, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	user, err := s.userRepo.GetByID(ctx, idn)
	if err != nil {
		return nil, err
	}

	// 获取用户角色
	roles, err := s.userRoleRepo.GetRolesByUserId(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	if len(roles) == 0 {
		return response.ToUserInfoResp(user, []string{}), nil
	}

	roleIDs := make([]uuid.UUID, len(roles))
	for i, role := range roles {
		roleIDs[i] = role.ID
	}
	// 获取用户角色关联的所有菜单
	menuTypes := []int16{3}
	menus, err := s.roleMenuRepo.GetMenusBYRoleIDs(ctx, menuTypes, roleIDs)
	if err != nil {
		return nil, err
	}
	if len(menus) == 0 {
		return response.ToUserInfoResp(user, make([]string, 0)), nil
	}
	buttons := make([]string, len(menus))
	for _, v := range menus {
		buttons = append(buttons, v.PermCode)
	}
	return response.ToUserInfoResp(user, buttons), nil
}

// GetUserRoutes 获取用户角色包含的权限信息
func (s authService) GetUserRoutes(ctx context.Context, id string) ([]*response.RouteResp, error) {
	idn, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	user, err := s.userRepo.GetByID(ctx, idn)
	if err != nil {
		return nil, err
	}

	// 获取用户角色
	roles, err := s.userRoleRepo.GetRolesByUserId(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	if len(roles) == 0 {
		return make([]*response.RouteResp, 0), nil
	}

	roleIDs := make([]uuid.UUID, len(roles))
	for i, role := range roles {
		roleIDs[i] = role.ID
	}
	// 获取用户角色关联的所有菜单
	menuTypes := []int16{1, 2}
	menus, err := s.roleMenuRepo.GetMenusBYRoleIDs(ctx, menuTypes, roleIDs)
	if err != nil {
		return nil, err
	}
	if len(menus) == 0 {
		return make([]*response.RouteResp, 0), nil
	}

	return response.ToRouteResp(menus), nil
}

func (s authService) Logout(ctx context.Context, userID string) error {
	id, err := uuid.Parse(userID)
	if err != nil {
		return err
	}
	if err := s.sessionRepo.UpdateExpiredByUserID(ctx, authorization.StatusRevoked, id); err != nil {
		return err
	}
	return nil
}
