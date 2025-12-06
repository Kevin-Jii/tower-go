package controller

import (
	"strconv"

	"github.com/Kevin-Jii/tower-go/middleware"
	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/service"
	"github.com/Kevin-Jii/tower-go/utils/http"
	"github.com/gin-gonic/gin"
)

type StoreAccountController struct {
	storeAccountService *service.StoreAccountService
}

func NewStoreAccountController(storeAccountService *service.StoreAccountService) *StoreAccountController {
	return &StoreAccountController{storeAccountService: storeAccountService}
}

// Create godoc
// @Summary 创建记账
// @Tags store-account
// @Accept json
// @Produce json
// @Security Bearer
// @Param account body model.CreateStoreAccountReq true "记账信息"
// @Success 200 {object} http.Response{data=model.StoreAccount}
// @Router /store-accounts [post]
func (c *StoreAccountController) Create(ctx *gin.Context) {
	storeID := middleware.GetStoreID(ctx)
	userID := middleware.GetUserID(ctx)

	var req model.CreateStoreAccountReq
	if !http.BindJSON(ctx, &req) {
		return
	}

	account, err := c.storeAccountService.Create(storeID, userID, &req)
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	http.Success(ctx, account)
}

// Get godoc
// @Summary 获取记账详情
// @Tags store-account
// @Produce json
// @Security Bearer
// @Param id path int true "记账ID"
// @Success 200 {object} http.Response{data=model.StoreAccount}
// @Router /store-accounts/{id} [get]
func (c *StoreAccountController) Get(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		http.Error(ctx, 400, "无效的ID")
		return
	}

	account, err := c.storeAccountService.Get(uint(id))
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	http.Success(ctx, account)
}

// List godoc
// @Summary 记账列表
// @Tags store-account
// @Produce json
// @Security Bearer
// @Param store_id query int false "门店ID"
// @Param product_id query int false "商品ID"
// @Param channel query string false "销售渠道"
// @Param order_source query string false "订单来源"
// @Param order_no query string false "订单编号"
// @Param start_date query string false "开始日期"
// @Param end_date query string false "结束日期"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} http.Response{data=[]model.StoreAccount}
// @Router /store-accounts [get]
func (c *StoreAccountController) List(ctx *gin.Context) {
	storeID := middleware.GetStoreID(ctx)
	roleCode := middleware.GetRoleCode(ctx)

	var req model.ListStoreAccountReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		http.Error(ctx, 400, err.Error())
		return
	}

	// 非管理员只能查看自己门店
	if roleCode != model.RoleCodeAdmin && roleCode != model.RoleCodeSuperAdmin {
		req.StoreID = storeID
	}

	list, total, err := c.storeAccountService.List(&req)
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	http.SuccessWithPagination(ctx, list, total, req.Page, req.PageSize)
}

// Update godoc
// @Summary 更新记账
// @Tags store-account
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "记账ID"
// @Param account body model.UpdateStoreAccountReq true "记账信息"
// @Success 200 {object} http.Response
// @Router /store-accounts/{id} [put]
func (c *StoreAccountController) Update(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		http.Error(ctx, 400, "无效的ID")
		return
	}

	var req model.UpdateStoreAccountReq
	if !http.BindJSON(ctx, &req) {
		return
	}

	if err := c.storeAccountService.Update(uint(id), &req); err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	http.Success(ctx, nil)
}

// Delete godoc
// @Summary 删除记账
// @Tags store-account
// @Produce json
// @Security Bearer
// @Param id path int true "记账ID"
// @Success 200 {object} http.Response
// @Router /store-accounts/{id} [delete]
func (c *StoreAccountController) Delete(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		http.Error(ctx, 400, "无效的ID")
		return
	}

	if err := c.storeAccountService.Delete(uint(id)); err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	http.Success(ctx, nil)
}

// Stats godoc
// @Summary 记账统计
// @Tags store-account
// @Produce json
// @Security Bearer
// @Param store_id query int false "门店ID"
// @Param start_date query string false "开始日期"
// @Param end_date query string false "结束日期"
// @Success 200 {object} http.Response
// @Router /store-accounts/stats [get]
func (c *StoreAccountController) Stats(ctx *gin.Context) {
	storeID := middleware.GetStoreID(ctx)
	roleCode := middleware.GetRoleCode(ctx)

	queryStoreID := uint(0)
	if s := ctx.Query("store_id"); s != "" {
		if id, err := strconv.ParseUint(s, 10, 32); err == nil {
			queryStoreID = uint(id)
		}
	}

	// 非管理员只能查看自己门店
	if roleCode != model.RoleCodeAdmin && roleCode != model.RoleCodeSuperAdmin {
		queryStoreID = storeID
	}

	startDate := ctx.Query("start_date")
	endDate := ctx.Query("end_date")

	stats, err := c.storeAccountService.GetStats(queryStoreID, startDate, endDate)
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	http.Success(ctx, stats)
}
