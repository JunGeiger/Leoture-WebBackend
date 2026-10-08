package handler

import (
	"LeotureWeb/internal/errors"
	"LeotureWeb/internal/service/user_role"
	"LeotureWeb/internal/types/request"
	"LeotureWeb/internal/types/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserRole struct {
	service user_role.Service
}

func newUserRoleHandler(svc user_role.Service) *UserRole {
	return &UserRole{service: svc}
}

// Update godoc
// @Summary create userRole
// @Description create user role correlations
// @Tags user
// @Accept json
// @Produce json
// @Param user body request.UserRoles true "user-role data"
// @Success 200 {object} response.SwaggerResponse
// @Failure 400,500 {object} response.SwaggerResponse
// @Router /api/v1/user-role/update [post]
func (h *UserRole) Update(c *gin.Context) {
	var data request.UserRoles
	if err := c.ShouldBindJSON(&data); err != nil {
		_ = c.Error(errors.ErrInvalidParam.WithErr(err))
		return
	}

	if err := h.service.Update(c.Request.Context(), &data); err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, response.Success[any](c, nil))
	return
}

// GetRolesByUserID godoc
// @Summary get role list
// @Description get role list by userID
// @Tags role
// @Accept json
// @Produce json
// @Param userID path string true "user data ID"
// @Success 200 {object} response.SwaggerResponse[data=[]response.RoleResp]
// @Failure 400,500 {object} response.SwaggerResponse
// @Router /api/v1/user-role/roles/:userID [get]
func (h *UserRole) GetRolesByUserID(c *gin.Context) {
	var data request.DataID
	err := c.ShouldBindUri(&data)
	if err != nil {
		_ = c.Error(errors.ErrInvalidParam.WithErr(err))
	}
	roleResp, err := h.service.GetRolesByUserID(c.Request.Context(), data.ID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, response.Success[[]*response.RoleResp](c, roleResp))
	return
}
