package handler

import (
	"LeotureWeb/internal/errors"
	"LeotureWeb/internal/service/role_menu"
	"LeotureWeb/internal/types/request"
	"LeotureWeb/internal/types/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

type RoleMenu struct {
	service role_menu.Service
}

func newRoleMenuHandler(svc role_menu.Service) *RoleMenu {
	return &RoleMenu{service: svc}
}

// Update godoc
// @Summary create roleMenu
// @Description create role menu correlations
// @Tags user
// @Accept json
// @Produce json
// @Param user body request.RoleMenus true "role-menu data"
// @Success 200 {object} response.SwaggerResponse
// @Failure 400,500 {object} response.SwaggerResponse
// @Router /api/v1/role-menu/update [post]
func (h *RoleMenu) Update(c *gin.Context) {
	var data request.RoleMenus
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

// GetMenusByRoleID godoc
// @Summary get menu list
// @Description get menu list by roleID
// @Tags role
// @Accept json
// @Produce json
// @Param userID path string true "role data ID"
// @Success 200 {object} response.SwaggerResponse{data=[]response.MenuTreeResp}
// @Failure 400,500 {object} response.SwaggerResponse
// @Router /api/v1/role-menu/menus/:roleID [get]
func (h *RoleMenu) GetMenusByRoleID(c *gin.Context) {
	var data request.DataID
	err := c.ShouldBindUri(&data)
	if err != nil {
		_ = c.Error(errors.ErrInvalidParam.WithErr(err))
	}
	roleResp, err := h.service.GetMenusByRoleID(c.Request.Context(), data.ID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, response.Success[[]*response.MenuTreeResp](c, roleResp))
	return
}
