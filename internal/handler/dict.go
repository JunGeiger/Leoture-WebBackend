package handler

import (
	"LeotureWeb/internal/errors"
	"LeotureWeb/internal/service/dict"
	"LeotureWeb/internal/types/request"
	"LeotureWeb/internal/types/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Dict 字典处理器
type Dict struct {
	service dict.Service
}

func newDictHandler(svc dict.Service) *Dict {
	return &Dict{service: svc}
}

// Create godoc
// @Summary create dict
// @Description create dictionary item
// @Tags dict
// @Accept json
// @Produce json
// @Param dict body request.DictCreate true "dict data"
// @Success 200 {object} response.SwaggerResponse
// @Failure 400,500 {object} response.SwaggerResponse
// @Router /api/v1/dict/create [post]
func (h *Dict) Create(c *gin.Context) {
	var data request.DictCreate
	if err := c.ShouldBindJSON(&data); err != nil {
		_ = c.Error(errors.ErrInvalidParam)
		return
	}
	if err := h.service.Create(c.Request.Context(), &data); err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, response.Success[any](c, nil))
}

// Update godoc
// @Summary update dict
// @Description update dictionary item by ID
// @Tags dict
// @Accept json
// @Produce json
// @Param dict body request.DictUpdate true "dict data"
// @Success 200 {object} response.SwaggerResponse
// @Failure 400,500 {object} response.SwaggerResponse
// @Router /api/v1/dict/update [post]
func (h *Dict) Update(c *gin.Context) {
	var data request.DictUpdate
	if err := c.ShouldBindJSON(&data); err != nil {
		_ = c.Error(errors.ErrInvalidParam)
		return
	}
	if err := h.service.Update(c.Request.Context(), &data); err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, response.Success[any](c, nil))
}

// Delete godoc
// @Summary delete dict
// @Description delete dictionary item by ID
// @Tags dict
// @Accept json
// @Produce json
// @Param id path string true "dict ID"
// @Success 200 {object} response.SwaggerResponse
// @Failure 400,500 {object} response.SwaggerResponse
// @Router /api/v1/dict/delete/:id [delete]
func (h *Dict) Delete(c *gin.Context) {
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
}

// GetByID godoc
// @Summary get dict by ID
// @Description get dictionary item by ID
// @Tags dict
// @Accept json
// @Produce json
// @Param id path string true "dict ID"
// @Success 200 {object} response.SwaggerResponse{data=response.DictResp}
// @Failure 400,500 {object} response.SwaggerResponse
// @Router /api/v1/dict/:id [get]
func (h *Dict) GetByID(c *gin.Context) {
	var data request.DataID
	err := c.ShouldBindUri(&data)
	if err != nil {
		_ = c.Error(errors.ErrInvalidParam.WithErr(err))
	}
	resp, err := h.service.GetByID(c.Request.Context(), data.ID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, response.Success[*response.DictResp](c, resp))
}

// GetByKey godoc
// @Summary get dict by category and key
// @Description get dictionary item by category and key
// @Tags dict
// @Accept json
// @Produce json
// @Param category path string true "dict category"
// @Param key path string true "dict key"
// @Success 200 {object} response.SwaggerResponse{data=response.DictResp}
// @Failure 400,500 {object} response.SwaggerResponse
// @Router /api/v1/dict/key/{category}/{key} [get]
func (h *Dict) GetByKey(c *gin.Context) {
	category := c.Param("category")
	key := c.Param("key")
	if category == "" || key == "" {
		_ = c.Error(errors.ErrInvalidParam)
		return
	}
	resp, err := h.service.GetByKey(c.Request.Context(), category, key)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, response.Success[*response.DictResp](c, resp))
}

// List godoc
// @Summary get dict list
// @Description get dictionary list by query
// @Tags dict
// @Accept json
// @Produce json
// @Param dict query request.DictQuery false "query params"
// @Success 200 {object} response.SwaggerResponse{data=[]response.DictResp}
// @Failure 400,500 {object} response.SwaggerResponse
// @Router /api/v1/dict/list [get]
func (h *Dict) List(c *gin.Context) {
	var data request.DictQuery
	if err := c.ShouldBindQuery(&data); err != nil {
		_ = c.Error(errors.ErrInvalidParam)
		return
	}
	resp, err := h.service.List(c.Request.Context(), &data)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, response.Success[[]*response.DictResp](c, resp))
}

// ListByPage godoc
// @Summary get dict list paginated
// @Description get dictionary list by query with pagination
// @Tags dict
// @Accept json
// @Produce json
// @Param dict query request.DictQuery false "query params"
// @Success 200 {object} response.SwaggerPaginatedData{list=[]response.DictResp}
// @Failure 400,500 {object} response.SwaggerResponse
// @Router /api/v1/dict/page [get]
func (h *Dict) ListByPage(c *gin.Context) {
	var data request.DictQuery
	if err := c.ShouldBindQuery(&data); err != nil {
		_ = c.Error(errors.ErrInvalidParam)
		return
	}
	resp, err := h.service.ListByPage(c.Request.Context(), &data)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, response.Success[*response.Paginated[response.DictResp]](c, resp))
}
