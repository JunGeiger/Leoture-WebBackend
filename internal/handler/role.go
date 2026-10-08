package handler

import (
	"LeotureWeb/internal/errors"
	"LeotureWeb/internal/service/role"
	"LeotureWeb/internal/types/request"
	"LeotureWeb/internal/types/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Role struct {
	service role.Service
}

func newRoleHandler(svc role.Service) *Role {
	return &Role{service: svc}
}

// Create godoc
// @Summary create role
// @Description create role
// @Tags role
// @Accept json
// @Produce json
// @Param role body request.RoleCreate true "role data"
// @Success 200 {object} response.SwaggerResponse
// @Failure 400,500 {object} response.SwaggerResponse
// @Router /api/v1/role/create [post]
func (h *Role) Create(c *gin.Context) {
	var data request.RoleCreate
	if err := c.ShouldBindJSON(&data); err != nil {
		_ = c.Error(errors.ErrInvalidParam.WithErr(err))
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
// @Summary update role
// @Description update role by roleID
// @Tags role
// @Accept json
// @Produce json
// @Param role body request.RoleUpdate true "role data"
// @Success 200 {object} response.SwaggerResponse
// @Failure 400,500 {object} response.SwaggerResponse
// @Router /api/v1/role/update [post]
func (h *Role) Update(c *gin.Context) {
	var data request.RoleUpdate
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
// @Summary delete role
// @Description delete role by roleID
// @Tags role
// @Accept json
// @Produce json
// @Param id path string true "role data ID"
// @Success 200 {object} response.SwaggerResponse
// @Failure 400,500 {object} response.SwaggerResponse
// @Router /api/v1/role/delete/:id [delete]
func (h *Role) Delete(c *gin.Context) {
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
// @Summary get role by roleID
// @Description get role by roleID
// @Tags role
// @Accept json
// @Produce json
// @Param id path string true "role data ID"
// @Success 200 {object} response.SwaggerResponse{data=response.RoleResp}
// @Failure 400,500 {object} response.SwaggerResponse
// @Router /api/v1/role/:id [get]
func (h *Role) GetByID(c *gin.Context) {
	var data request.DataID
	err := c.ShouldBindUri(&data)
	if err != nil {
		_ = c.Error(errors.ErrInvalidParam.WithErr(err))
	}
	roleResp, err := h.service.GetByID(c.Request.Context(), data.ID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, response.Success[*response.RoleResp](c, roleResp))
	return
}

// GetByCode godoc
// @Summary get role by code
// @Description get role by code
// @Tags role
// @Accept json
// @Produce json
// @Param code path string true "role data code"
// @Success 200 {object} response.SwaggerResponse{data=response.RoleResp}
// @Failure 400,500 {object} response.SwaggerResponse
// @Router /api/v1/role/code/:code [get]
func (h *Role) GetByCode(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		_ = c.Error(errors.ErrInvalidParam)
		return
	}
	roleResp, err := h.service.GetByCode(c.Request.Context(), code)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, response.Success[*response.RoleResp](c, roleResp))
	return
}

// List godoc
// @Summary get role list
// @Description get role list by query
// @Tags role
// @Accept json
// @Produce json
// @Param role query request.RoleQuery false "role data by query"
// @Success 200 {object} response.SwaggerResponse{data=[]response.RoleResp}
// @Failure 400,500 {object} response.SwaggerResponse
// @Router /api/v1/role/list [get]
func (h *Role) List(c *gin.Context) {
	var data request.RoleQuery
	if err := c.ShouldBindQuery(&data); err != nil {
		_ = c.Error(errors.ErrInvalidParam)
		return
	}

	roleResp, err := h.service.List(c.Request.Context(), &data)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, response.Success[[]*response.RoleResp](c, roleResp))
	return
}

// ListByPage godoc
// @Summary get role list by paginated
// @Description get role list by query and paginated
// @Tags role
// @Accept json
// @Produce json
// @Param role query request.RoleQuery false "role data by query and paginated"
// @Success 200 {object} response.SwaggerPaginatedData{list=[]response.RoleResp}
// @Failure 400,500 {object} response.SwaggerResponse
// @Router /api/v1/role/page [get]
func (h *Role) ListByPage(c *gin.Context) {
	var data request.RoleQuery
	if err := c.ShouldBindQuery(&data); err != nil {
		_ = c.Error(errors.ErrInvalidParam)
		return
	}

	roleResp, err := h.service.ListByPage(c.Request.Context(), &data)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, response.Success[*response.Paginated[response.RoleResp]](c, roleResp))
	return
}
