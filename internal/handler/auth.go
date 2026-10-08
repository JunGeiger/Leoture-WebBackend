package handler

import (
	"LeotureWeb/internal/errors"
	"LeotureWeb/internal/service/auth"
	"LeotureWeb/internal/types/request"
	"LeotureWeb/internal/types/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Auth struct {
	service auth.Service
}

func newAuthHandler(svc auth.Service) *Auth {
	return &Auth{service: svc}
}

// LoginByUser godoc
// @Summary user login api
// @Description return accessToken、sessionExpiry、refreshToken
// @Tags login
// @Accept json
// @Produce json
// @Param loginUser body request.LoginUser true "user login data"
// @Success 200 {object} response.SwaggerResponse{data=response.LoginResp}
// @Failure 400,500 {object} response.SwaggerResponse
// @Router /api/v1/auth/login [post]
func (h *Auth) LoginByUser(c *gin.Context) {
	var data request.LoginUser
	err := c.ShouldBindJSON(&data)
	data.IPAddress = c.ClientIP()
	data.UserAgent = c.Request.UserAgent()
	data.OriginUrl = c.GetHeader("Origin")
	data.Location = c.GetHeader("x-client-geo-city")
	if err != nil {
		_ = c.Error(errors.ErrInvalidParam.WithErr(err))
		return
	}

	loginResp, err := h.service.Login(c.Request.Context(), &data)
	if err != nil {
		_ = c.Error(errors.ErrInternal.WithErr(err))
		return
	}

	c.JSON(http.StatusOK, response.Success[*response.LoginResp](c, loginResp))
	return
}

// LoginByRefreshToken godoc
// @Summary 系统Token刷新接口
// @Description 系统Token刷新接口，返回accessToken、refreshToken
// @Tags refresh
// @Accept json
// @Produce json
// @Param refreshToken body request.RefreshToken true "用户RefreshToken"
// @Success 200 {object} response.SwaggerResponse{data=response.LoginResp}
// @Failure 400,500 {object} response.SwaggerResponse
// @Router /api/v1/auth/refreshToken [post]
func (h *Auth) LoginByRefreshToken(c *gin.Context) {
	var data request.RefreshToken
	err := c.ShouldBindJSON(&data)
	if err != nil {
		_ = c.Error(errors.ErrInvalidParam.WithErr(err))
		return
	}
	loginResp, err := h.service.RefreshToken(c.Request.Context(), &data)
	if err != nil {
		_ = c.Error(errors.ErrInternal.WithErr(err))
		return
	}

	c.JSON(http.StatusOK, response.Success[*response.LoginRefreshResp](c, loginResp))
	return
}

// GetUserInfo godoc
// @Summary get user info
// @Description 根据UserID获取用户信息，自动从Token解析UserID
// @Tags login
// @Accept json
// @Produce json
// @Param userID uri request.DataID true "用户ID"
// @Success 200 {object} response.SwaggerResponse
// @Failure 400,500 {object} response.SwaggerResponse
// @Router /api/v1/auth/getUserInfo [get]
func (h *Auth) GetUserInfo(c *gin.Context) {
	var data request.DataID
	claims := request.GetClaims(c)
	if claims == nil {
		_ = c.Error(errors.ErrInvalidParam)
		return
	}
	data.ID = claims.Subject
	userResp, err := h.service.GetUserInfo(c.Request.Context(), data.ID)
	if err != nil {
		_ = c.Error(errors.ErrInternal.WithErr(err))
		return
	}

	c.JSON(http.StatusOK, response.Success[*response.UserInfoResp](c, userResp))
	return
}

// GetConstantRoutes godoc
// @Summary get constant routes
// @Description 根据UserID获取用户所属角色下的静态路由，后端暂时不保存静态路由数据，返回空数组
// @Tags ConstantRoutes
// @Accept json
// @Produce json
// @Param userID query request.DataID false "用户ID"
// @Success 200 {object} response.SwaggerResponse{data=[]response.ConstantRouteResp}
// @Failure 400,500 {object} response.SwaggerResponse
// @Router /api/v1/route/getConstantRoutes [get]
func (h *Auth) GetConstantRoutes(c *gin.Context) {
	constantRouteResp := make([]*response.UserRoutesResp, 0)
	c.JSON(http.StatusOK, response.Success[[]*response.UserRoutesResp](c, constantRouteResp))
	return
}

// GetUserRoutes godoc
// @Summary get user routes
// @Description 根据UserID获取用户所属角色下的菜单路由数据，自动从Token解析UserID
// @Tags UserRoutes
// @Accept json
// @Produce json
// @Param userID query request.DataID false "用户ID"
// @Success 200 {object} response.SwaggerResponse{data=[]response.RouteResp}
// @Failure 400,500 {object} response.SwaggerResponse
// @Router /api/v1/route/getUserRoutes [get]
func (h *Auth) GetUserRoutes(c *gin.Context) {
	var data request.DataID
	claims := request.GetClaims(c)
	if claims == nil {
		_ = c.Error(errors.ErrInvalidParam)
		return
	}
	data.ID = claims.Subject
	routeResp, err := h.service.GetUserRoutes(c.Request.Context(), data.ID)
	if err != nil {
		_ = c.Error(errors.ErrInternal.WithErr(err))
		return
	}
	userRoutes := &response.UserRoutesResp{
		Routes: routeResp,
		Home:   "home",
	}
	c.JSON(http.StatusOK, response.Success[*response.UserRoutesResp](c, userRoutes))
	return
}

// GetUserInfo godoc
// @Summary 系统登录接口
// @Description 系统登录接口，返回accessToken、refreshToken、用户权限信息
// @Tags login
// @Accept json
// @Produce json
// @Param userID uri request.DataID true "用户ID"
// @Success 200 {object} response.SwaggerResponse
// @Failure 400,500 {object} response.SwaggerResponse
// @Router /api/v1/route/isRouteExist [get]
func (h *Auth) IsRouteExist(c *gin.Context) {
	var data request.DataID
	claims := request.GetClaims(c)
	if claims == nil {
		_ = c.Error(errors.ErrInvalidParam)
		return
	}
	data.ID = claims.Subject
	userResp, err := h.service.GetUserInfo(c.Request.Context(), data.ID)
	if err != nil {
		_ = c.Error(errors.ErrInternal.WithErr(err))
		return
	}
	c.JSON(http.StatusOK, response.Success[*response.UserInfoResp](c, userResp))
	return
}

// Logout godoc
// @Summary 系统登录接口
// @Description 系统登录接口，返回accessToken、refreshToken、用户权限信息
// @Tags login
// @Accept json
// @Produce json
// @Param userID uri request.DataID true "用户ID"
// @Success 200 {object} response.SwaggerResponse
// @Failure 400,500 {object} response.SwaggerResponse
// @Router /api/v1/auth/logout/:id [get]
func (h *Auth) Logout(c *gin.Context) {
	var data request.DataID
	if err := c.ShouldBindUri(&data); err != nil {
		_ = c.Error(errors.ErrInvalidParam.WithErr(err))
		return
	}
	if err := h.service.Logout(c.Request.Context(), data.ID); err != nil {
		_ = c.Error(errors.ErrInternal.WithErr(err))
		return
	}
	return
}
