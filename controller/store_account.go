package controller

import (
	"strconv"

	"github.com/Kevin-Jii/tower-go/middleware"
	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/pkg/apicode"
	"github.com/Kevin-Jii/tower-go/service"
	"github.com/Kevin-Jii/tower-go/utils/excelxml"
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
// @Description 创建记账单，支持多个商品；items.product_id=0 时为自定义明细，items.product_name 必填且 items.price 必填，不扣库存。系统商品时 items.price 不传或传0则按 items.unit 自动选价；items.amount 不传或传0时按 price*quantity 计算
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

	account, err := c.storeAccountService.GetScoped(uint(id), middleware.GetStoreID(ctx), middleware.HQUnboundAdmin(ctx))
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
// @Param payment_status query int false "支付状态 1=已支付 2=未支付"
// @Param member_keyword query string false "会员搜索（手机号/姓名/会员ID）"
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

	if !middleware.HQUnboundAdmin(ctx) {
		req.StoreID = storeID
	}

	middleware.AttachAuthContextToHTTPRequest(ctx)

	list, total, err := c.storeAccountService.List(ctx.Request.Context(), &req)
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	http.SuccessWithPagination(ctx, list, total, req.Page, req.PageSize)
}

// Export 导出单日记账 Excel。
func (c *StoreAccountController) Export(ctx *gin.Context) {
	storeID := middleware.GetStoreID(ctx)

	var req model.ListStoreAccountReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		http.Error(ctx, 400, err.Error())
		return
	}
	startDate, endDate, date := exportDateQuery(ctx.Query("date"), req.StartDate, req.EndDate)
	if date == "" {
		http.Error(ctx, 400, "请选择导出日期")
		return
	}
	req.StartDate = startDate
	req.EndDate = endDate
	req.Page = 1
	req.PageSize = exportPageSize

	if !middleware.HQUnboundAdmin(ctx) {
		req.StoreID = storeID
	}
	middleware.AttachAuthContextToHTTPRequest(ctx)

	list, _, err := c.storeAccountService.List(ctx.Request.Context(), &req)
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	rows := make([][]interface{}, 0, len(list))
	for _, item := range list {
		rows = append(rows, []interface{}{
			item.AccountDate,
			item.AccountNo,
			storeName(item.Store),
			item.Channel,
			item.OrderNo,
			memberName(item.Member),
			storeAccountPaymentLabel(item.PaymentStatus),
			formatAmount(item.TotalAmount),
			formatAmount(item.OtherExpenseAmount),
			formatAmount(item.ErrandFee),
			formatAmount(item.RoundAmount),
			formatAmount(item.GiftWineCostAmount),
			formatAmount(item.NetIncomeAmount),
			accountItemsText(item.Items),
			accountConsumablesText(item.Consumables),
			userName(item.Operator, ""),
			item.Remark,
			item.CreatedAt,
		})
	}

	data := excelxml.Build([]excelxml.Sheet{{
		Name:    "记账",
		Headers: []string{"日期", "记账编号", "门店", "渠道", "订单号", "会员", "支付状态", "销售额", "其他支出", "跑腿费", "抹零", "赠酒成本", "净收入", "商品明细", "消耗品", "操作人", "备注", "创建时间"},
		Rows:    rows,
	}})
	http.File(ctx, data, excelxml.Filename("store-accounts-"+date))
}

// Update godoc
// @Summary 更新记账
// @Description 更新记账信息，仅允许当前营业日内修改
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

	var req model.UpdateStoreAccountReq
	if !http.BindJSON(ctx, &req) {
		return
	}

	if err := c.storeAccountService.UpdateScoped(uint(id), middleware.GetStoreID(ctx), middleware.HQUnboundAdmin(ctx), &req); err != nil {
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
	http.Error(ctx, 400, "记账记录不允许删除")
}

// BindConsumables godoc
// @Summary 绑定记账单消耗品
// @Tags 门店记账
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "记账ID"
// @Param body body model.BindStoreAccountConsumablesReq true "消耗品信息"
// @Success 200 {object} http.Response
// @Router /store-accounts/{id}/consumables [post]
func (c *StoreAccountController) BindConsumables(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		http.ErrorApp(ctx, apicode.InvalidID)
		return
	}
	var req model.BindStoreAccountConsumablesReq
	if !http.BindJSON(ctx, &req) {
		return
	}
	if err := c.storeAccountService.BindConsumablesScoped(uint(id), middleware.GetStoreID(ctx), middleware.HQUnboundAdmin(ctx), &req); err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, nil)
}

func (c *StoreAccountController) CreateConsumableProduct(ctx *gin.Context) {
	storeID := middleware.GetStoreID(ctx)

	var req model.UpsertStoreAccountConsumableProductReq
	if !http.BindJSON(ctx, &req) {
		return
	}
	product, err := c.storeAccountService.CreateConsumableProduct(storeID, &req, middleware.HQUnboundAdmin(ctx))
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, product)
}

func (c *StoreAccountController) ListConsumableProducts(ctx *gin.Context) {
	storeID := middleware.GetStoreID(ctx)

	var req model.ListStoreAccountConsumableProductReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		http.Error(ctx, 400, err.Error())
		return
	}
	if !middleware.HQUnboundAdmin(ctx) {
		req.StoreID = storeID
	}
	list, total, err := c.storeAccountService.ListConsumableProducts(ctx.Request.Context(), &req)
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.SuccessWithPagination(ctx, list, total, req.Page, req.PageSize)
}

func (c *StoreAccountController) UpdateConsumableProduct(ctx *gin.Context) {
	id, ok := http.ParseUintParam(ctx, "id")
	if !ok {
		return
	}

	var req model.UpsertStoreAccountConsumableProductReq
	if !http.BindJSON(ctx, &req) {
		return
	}
	product, err := c.storeAccountService.UpdateConsumableProduct(id, middleware.GetStoreID(ctx), &req, middleware.HQUnboundAdmin(ctx))
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, product)
}

func (c *StoreAccountController) DeleteConsumableProduct(ctx *gin.Context) {
	id, ok := http.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	if err := c.storeAccountService.DeleteConsumableProduct(id, middleware.GetStoreID(ctx), middleware.HQUnboundAdmin(ctx)); err != nil {
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
	queryStoreID := middleware.ResolveQueryStoreID(ctx, "store_id")
	startDate := ctx.Query("start_date")
	endDate := ctx.Query("end_date")

	stats, err := c.storeAccountService.GetStats(queryStoreID, startDate, endDate)
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	http.Success(ctx, stats)
}
