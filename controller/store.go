package controller

import (
	"github.com/Kevin-Jii/tower-go/middleware"
	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/service"
	"github.com/Kevin-Jii/tower-go/utils/http"

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
// @Tags 门店管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param store body model.CreateStoreReq true "门店信息"
// @Success 200 {object} http.Response
// @Router /stores [post]
func (c *StoreController) CreateStore(ctx *gin.Context) {
	// 管理员校验
	if !http.RequireAdmin(ctx) {
		return
	}

	var req model.CreateStoreReq
	if !http.BindJSON(ctx, &req) {
		return
	}

	if err := c.storeService.CreateStore(&req); err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	http.Success(ctx, nil)
}

// GetStore godoc
// @Summary 获取门店详情
// @Description 获取门店详细信息
// @Tags 门店管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "门店ID"
// @Success 200 {object} http.Response{data=model.Store}
// @Router /stores/{id} [get]
func (c *StoreController) GetStore(ctx *gin.Context) {
	id, ok := http.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	store, err := c.storeService.GetStore(id)
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	http.Success(ctx, store)
}

// ListStores godoc
// @Summary 门店列表
// @Description 获取全部门店列表（不分页）
// @Tags 门店管理
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} http.Response{data=[]model.Store} "返回全部门店数据，meta 包含 total"
// @Router /stores [get]
func (c *StoreController) ListStores(ctx *gin.Context) {
	stores, total, err := c.storeService.ListStores()
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	// 使用分页格式返回，但不实际分页（page=1, pageSize=total）
	http.SuccessWithPagination(ctx, stores, total, 1, int(total))
}

// ListAllStores godoc
// @Summary 全部门店（无分页）
// @Description 返回全部门店列表。默认仅总部管理员可访问，如需开放可去掉权限判断
// @Tags 门店管理
// @Security Bearer
// @Produce json
// @Success 200 {object} http.Response{data=[]model.Store}
// @Router /stores/all [get]
func (c *StoreController) ListAllStores(ctx *gin.Context) {
	// 权限限制：仅 admin，避免普通门店账号看到其他门店
	if !middleware.IsAdmin(ctx) {
		http.Error(ctx, 403, "Only admin can list all stores")
		return
	}
	stores, _, err := c.storeService.ListStores()
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, stores)
}

// UpdateStore godoc
// @Summary 更新门店信息
// @Description 更新门店信息（仅总部管理员）
// @Tags 门店管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "门店ID"
// @Param store body model.UpdateStoreReq true "门店信息"
// @Success 200 {object} http.Response
// @Router /stores/{id} [put]
func (c *StoreController) UpdateStore(ctx *gin.Context) {
	if !http.RequireAdmin(ctx) {
		return
	}
	id, ok := http.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	var req model.UpdateStoreReq
	if !http.BindJSON(ctx, &req) {
		return
	}
	if err := c.storeService.UpdateStore(id, &req); err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	http.Success(ctx, nil)
}

// DeleteStore godoc
// @Summary 删除门店
// @Description 删除门店（仅总部管理员）
// @Tags 门店管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "门店ID"
// @Success 200 {object} http.Response
// @Router /stores/{id} [delete]
func (c *StoreController) DeleteStore(ctx *gin.Context) {
	if !http.RequireAdmin(ctx) {
		return
	}
	id, ok := http.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	if err := c.storeService.DeleteStore(id); err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	http.Success(ctx, nil)
}
