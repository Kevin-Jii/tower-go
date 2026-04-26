package controller

import (
	"strconv"
	"time"

	"github.com/Kevin-Jii/tower-go/middleware"
	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/pkg/apicode"
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
// @Description 创建记账单，支持多个商品；items.price 不传或传0时，后端按 items.unit 自动选价（单位含“瓶”取 bottle_price，含“箱”取 case_price）；items.amount 不传或传0时自动按 price*quantity 计算
// @Tags 门店记账
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
// @Tags 门店记账
// @Produce json
// @Security Bearer
// @Param id path int true "记账ID"
// @Success 200 {object} http.Response{data=model.StoreAccount}
// @Router /store-accounts/{id} [get]
func (c *StoreAccountController) Get(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		http.ErrorApp(ctx, apicode.InvalidID)
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
// @Tags 门店记账
// @Produce json
// @Security Bearer
// @Param store_id query int false "门店ID"
// @Param channel query string false "渠道来源"
// @Param order_no query string false "订单编号"
// @Param tag_code query string false "标签编码"
// @Param start_date query string false "开始日期"
// @Param end_date query string false "结束日期"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} http.Response{data=[]model.StoreAccount}
// @Router /store-accounts [get]
func (c *StoreAccountController) List(ctx *gin.Context) {
	storeID := middleware.GetStoreID(ctx)

	var req model.ListStoreAccountReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		http.Error(ctx, 400, err.Error())
		return
	}

	req.DataScope, req.UserID, req.RoleCode = middleware.ListRBAC(ctx)
	// 非管理员只能查看自己门店；管理员可按 query 传 store_id
	if !middleware.IsAdmin(ctx) {
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
// @Description 更新记账信息，仅限创建后24小时内可修改（管理员不受限制）
// @Tags 门店记账
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
		http.ErrorApp(ctx, apicode.InvalidID)
		return
	}

	// 获取记账记录
	account, err := c.storeAccountService.Get(uint(id))
	if err != nil {
		http.ErrorApp(ctx, apicode.StoreAccountGone)
		return
	}

	// 非管理员检查24小时限制
	roleCode := middleware.GetRoleCode(ctx)
	if roleCode != model.RoleCodeAdmin && roleCode != model.RoleCodeSuperAdmin {
		if time.Since(account.CreatedAt) > 24*time.Hour {
			http.ErrorApp(ctx, apicode.StoreAccountEditTimeout)
			return
		}
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
// @Tags 门店记账
// @Produce json
// @Security Bearer
// @Param id path int true "记账ID"
// @Success 200 {object} http.Response
// @Router /store-accounts/{id} [delete]
func (c *StoreAccountController) Delete(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		http.ErrorApp(ctx, apicode.InvalidID)
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
// @Tags 门店记账
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
