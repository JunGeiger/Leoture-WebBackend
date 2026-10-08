package handler

import (
	"LeotureWeb/internal/errors"
	"LeotureWeb/internal/service/menu"
	"LeotureWeb/internal/types/request"
	"LeotureWeb/internal/types/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Menu 菜单处理器
type Menu struct {
	service menu.Service
}

func newMenuHandler(svc menu.Service) *Menu {
	return &Menu{service: svc}
}

// Create godoc
// @Summary 创建菜单
// @Description 创建新菜单
// @Tags menu
// @Accept json
// @Produce json
// @Param menu body request.MenuCreate true "菜单数据"
// @Success 200 {object} response.SwaggerResponse
// @Failure 400,500 {object} response.SwaggerResponse
// @Router /api/v1/menu/create [post]
func (h *Menu) Create(c *gin.Context) {
	var data request.MenuCreate
	if err := c.ShouldBindJSON(&data); err != nil {
		_ = c.Error(errors.ErrInvalidParam.WithErr(err))
		return
	}

	if err := h.service.Create(c.Request.Context(), &data); err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, response.Success[any](c, nil))
}

// Update godoc
// @Summary 更新菜单
// @Description 更新菜单信息
// @Tags menu
// @Accept json
// @Produce json
// @Param menu body request.MenuUpdate true "菜单数据"
// @Success 200 {object} response.SwaggerResponse
// @Failure 400,500 {object} response.SwaggerResponse
// @Router /api/v1/menu/update [post]
func (h *Menu) Update(c *gin.Context) {
	var data request.MenuUpdate
	if err := c.ShouldBindJSON(&data); err != nil {
		_ = c.Error(errors.ErrInvalidParam.WithErr(err))
		return
	}

	if err := h.service.Update(c.Request.Context(), &data); err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, response.Success[any](c, nil))
}

// Delete godoc
// @Summary 删除菜单
// @Description 软删除菜单
// @Tags menu
// @Accept json
// @Produce json
// @Param id path string true "菜单ID"
// @Success 200 {object} response.SwaggerResponse
// @Failure 400,500 {object} response.SwaggerResponse
// @Router /api/v1/menu/delete/{id} [delete]
func (h *Menu) Delete(c *gin.Context) {
	var data request.DataID
	if err := c.ShouldBindUri(&data); err != nil {
		_ = c.Error(errors.ErrInvalidParam.WithErr(err))
		return
	}

	if err := h.service.Delete(c.Request.Context(), data.ID); err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, response.Success[any](c, nil))
}

// GetByID godoc
// @Summary 获取菜单详情
// @Description 根据ID获取菜单详情
// @Tags menu
// @Accept json
// @Produce json
// @Param id path string true "菜单ID"
// @Success 200 {object} response.SwaggerResponse{data=response.MenuResp}
// @Failure 400,500 {object} response.SwaggerResponse
// @Router /api/v1/menu/{id} [get]
func (h *Menu) GetByID(c *gin.Context) {
	var data request.DataID
	if err := c.ShouldBindUri(&data); err != nil {
		_ = c.Error(errors.ErrInvalidParam.WithErr(err))
		return
	}

	menuResp, err := h.service.GetByID(c.Request.Context(), data.ID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, response.Success[*response.MenuResp](c, menuResp))
}

// List godoc
// @Summary 获取菜单列表
// @Description 根据条件获取菜单列表
// @Tags menu
// @Accept json
// @Produce json
// @Param menu query request.MenuQuery false "查询条件"
// @Success 200 {object} response.SwaggerResponse{data=[]response.MenuResp}
// @Failure 400,500 {object} response.SwaggerResponse
// @Router /api/v1/menu/list [get]
func (h *Menu) List(c *gin.Context) {
	var data request.MenuQuery
	if err := c.ShouldBindQuery(&data); err != nil {
		_ = c.Error(errors.ErrInvalidParam.WithErr(err))
		return
	}

	menuResp, err := h.service.List(c.Request.Context(), &data)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, response.Success[[]*response.MenuResp](c, menuResp))
}

// ListByPage godoc
// @Summary 分页获取菜单列表
// @Description 根据条件分页获取菜单列表
// @Tags menu
// @Accept json
// @Produce json
// @Param menu query request.MenuQuery false "查询条件"
// @Success 200 {object} response.SwaggerPaginatedData{data=[]response.MenuResp}
// @Failure 400,500 {object} response.SwaggerResponse
// @Router /api/v1/menu/page [get]
func (h *Menu) ListByPage(c *gin.Context) {
	var data request.MenuQuery
	if err := c.ShouldBindQuery(&data); err != nil {
		_ = c.Error(errors.ErrInvalidParam.WithErr(err))
		return
	}

	menuResp, err := h.service.ListByPage(c.Request.Context(), &data)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, response.Success[*response.Paginated[response.MenuResp]](c, menuResp))
}

// Tree godoc
// @Summary 分页获取菜单列表
// @Description 根据条件分页获取菜单列表
// @Tags menu
// @Accept json
// @Produce json
// @Param menuID query request.DataID false "菜单ID"
// @Success 200 {object} response.SwaggerPaginatedData{list=[]response.MenuResp}
// @Failure 400,500 {object} response.SwaggerResponse
// @Router /api/v1/menu/tree [get]
func (h *Menu) Tree(c *gin.Context) {
	treeResp, err := h.service.Tree(c.Request.Context())
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, response.Success[[]*response.MenuTreeResp](c, treeResp))
}
