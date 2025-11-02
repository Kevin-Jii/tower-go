package controller

import (
	"net/http"
	"strconv"
	"tower-go/middleware"
	"tower-go/model"
	"tower-go/service"
	"tower-go/utils"

	"github.com/gin-gonic/gin"
)

type StoreController struct {
	storeService *service.StoreService
}

func NewStoreController(storeService *service.StoreService) *StoreController {
	return &StoreController{storeService: storeService}
}

// CreateStore godoc
// @Summary 创建门店
// @Description 创建新门店（仅总部管理员）
// @Tags stores
// @Accept json
// @Produce json
// @Security Bearer
// @Param store body model.CreateStoreReq true "门店信息"
// @Success 200 {object} utils.StandardResponse
// @Router /stores [post]
func (c *StoreController) CreateStore(ctx *gin.Context) {
	// 检查是否是总部管理员
	if !middleware.IsAdmin(ctx) {
		utils.Error(ctx, http.StatusForbidden, "Only admin can create stores")
		return
	}

	var req model.CreateStoreReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := c.storeService.CreateStore(&req); err != nil {
		utils.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(ctx, nil)
}

// GetStore godoc
// @Summary 获取门店详情
// @Description 获取门店详细信息
// @Tags stores
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "门店ID"
// @Success 200 {object} utils.StandardResponse{data=model.Store}
// @Router /stores/{id} [get]
func (c *StoreController) GetStore(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "Invalid store ID")
		return
	}

	store, err := c.storeService.GetStore(uint(id))
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(ctx, store)
}

// ListStores godoc
// @Summary 门店列表
// @Description 获取门店列表，支持分页
// @Tags stores
// @Accept json
// @Produce json
// @Security Bearer
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} utils.StandardResponse{data=[]model.Store} "分页 meta: total,page,page_size,page_count,has_more"
// @Router /stores [get]
func (c *StoreController) ListStores(ctx *gin.Context) {
	page := utils.GetPage(ctx)
	pageSize := utils.GetPageSize(ctx)

	stores, total, err := c.storeService.ListStores(page, pageSize)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessWithPagination(ctx, stores, total, page, pageSize)
}

// UpdateStore godoc
// @Summary 更新门店信息
// @Description 更新门店信息（仅总部管理员）
// @Tags stores
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "门店ID"
// @Param store body model.UpdateStoreReq true "门店信息"
// @Success 200 {object} utils.StandardResponse
// @Router /stores/{id} [put]
func (c *StoreController) UpdateStore(ctx *gin.Context) {
	// 检查是否是总部管理员
	if !middleware.IsAdmin(ctx) {
		utils.Error(ctx, http.StatusForbidden, "Only admin can update stores")
		return
	}

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "Invalid store ID")
		return
	}

	var req model.UpdateStoreReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := c.storeService.UpdateStore(uint(id), &req); err != nil {
		utils.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(ctx, nil)
}

// DeleteStore godoc
// @Summary 删除门店
// @Description 删除门店（仅总部管理员）
// @Tags stores
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "门店ID"
// @Success 200 {object} utils.StandardResponse
// @Router /stores/{id} [delete]
func (c *StoreController) DeleteStore(ctx *gin.Context) {
	// 检查是否是总部管理员
	if !middleware.IsAdmin(ctx) {
		utils.Error(ctx, http.StatusForbidden, "Only admin can delete stores")
		return
	}

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "Invalid store ID")
		return
	}

	if err := c.storeService.DeleteStore(uint(id)); err != nil {
		utils.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(ctx, nil)
}
