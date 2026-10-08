package handler

import (
	"LeotureWeb/internal/errors"
	"LeotureWeb/internal/service/user"
	"LeotureWeb/internal/types/request"
	"LeotureWeb/internal/types/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

type User struct {
	service user.Service
}

func newUserHandler(svc user.Service) *User {
	return &User{service: svc}
}

// Create godoc
// @Summary create user
// @Description create user
// @Tags user
// @Accept json
// @Produce json
// @Param user body request.UserCreate true "user data"
// @Success 200 {object} response.SwaggerResponse
// @Failure 400,500 {object} response.SwaggerResponse
// @Router /api/v1/user/create [post]
func (h *User) Create(c *gin.Context) {
	var data request.UserCreate
	if err := c.ShouldBindJSON(&data); err != nil {
		_ = c.Error(errors.ErrInvalidParam)
		return
	}

	if err := h.service.Create(c.Request.Context(), &data); err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, response.Success[any](c, nil))
	return
}

// Update godoc
// @Summary update user
// @Description update user by userID
// @Tags user
// @Accept json
// @Produce json
// @Param user body request.UserUpdate true "user data"
// @Success 200 {object} response.SwaggerResponse
// @Failure 400,500 {object} response.SwaggerResponse
// @Router /api/v1/user/update [post]
func (h *User) Update(c *gin.Context) {
	var data request.UserUpdate
	if err := c.ShouldBindJSON(&data); err != nil {
		_ = c.Error(errors.ErrInvalidParam)
		return
	}

	if err := h.service.Update(c.Request.Context(), &data); err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, response.Success[any](c, nil))
	return
}

// Delete godoc
// @Summary delete user
// @Description delete user by userID
// @Tags user
// @Accept json
// @Produce json
// @Param id path string true "user data ID"
// @Success 200 {object} response.SwaggerResponse
// @Failure 400,500 {object} response.SwaggerResponse
// @Router /api/v1/user/delete/:id [delete]
func (h *User) Delete(c *gin.Context) {
	var data request.DataID
	err := c.ShouldBindUri(&data)
	if err != nil {
		_ = c.Error(errors.ErrInvalidParam.WithErr(err))
	}
	if err := h.service.Delete(c.Request.Context(), data.ID); err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, response.Success[any](c, nil))
	return
}

// GetByID godoc
// @Summary get user by userID
// @Description get user by userID
// @Tags user
// @Accept json
// @Produce json
// @Param id path string true "user data ID"
// @Success 200 {object} response.SwaggerResponse{data=response.UserResp}
// @Failure 400,500 {object} response.SwaggerResponse
// @Router /api/v1/user/:id [get]
func (h *User) GetByID(c *gin.Context) {
	var data request.DataID
	err := c.ShouldBindUri(&data)
	if err != nil {
		_ = c.Error(errors.ErrInvalidParam.WithErr(err))
	}
	userResp, err := h.service.GetByID(c.Request.Context(), data.ID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, response.Success[*response.UserResp](c, userResp))
	return
}

// GetByUsername godoc
// @Summary get user by username
// @Description get user by username
// @Tags user
// @Accept json
// @Produce json
// @Param username path string true "user username"
// @Success 200 {object} response.SwaggerResponse{data=response.UserResp}
// @Failure 400,500 {object} response.SwaggerResponse
// @Router /api/v1/user/username/:username [get]
func (h *User) GetByUsername(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		_ = c.Error(errors.ErrInvalidParam)
		return
	}
	userResp, err := h.service.GetByUsername(c.Request.Context(), username)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, response.Success[*response.UserResp](c, userResp))
	return
}

// List godoc
// @Summary get user list
// @Description get user list by query
// @Tags user
// @Accept json
// @Produce json
// @Param user query request.UserQuery false "user data by query"
// @Success 200 {object} response.SwaggerResponse{data=[]response.UserResp}
// @Failure 400,500 {object} response.SwaggerResponse
// @Router /api/v1/user/list [get]
func (h *User) List(c *gin.Context) {
	var data request.UserQuery
	if err := c.ShouldBindQuery(&data); err != nil {
		_ = c.Error(errors.ErrInvalidParam)
		return
	}

	userResp, err := h.service.List(c.Request.Context(), &data)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, response.Success[[]*response.UserResp](c, userResp))
	return
}

// ListByPage godoc
// @Summary get user list by paginated
// @Description get user list by query and paginated
// @Tags user
// @Accept json
// @Produce json
// @Param user query request.UserQuery false "user data by query and paginated"
// @Success 200 {object} response.SwaggerPaginatedData{list=[]response.UserResp}
// @Failure 400,500 {object} response.SwaggerResponse
// @Router /api/v1/user/page [get]
func (h *User) ListByPage(c *gin.Context) {
	var data request.UserQuery
	if err := c.ShouldBindQuery(&data); err != nil {
		_ = c.Error(errors.ErrInvalidParam)
		return
	}

	userResp, err := h.service.ListByPage(c.Request.Context(), &data)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, response.Success[*response.Paginated[response.UserResp]](c, userResp))
	return
}
