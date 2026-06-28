package controller

import (
	"github.com/Kevin-Jii/tower-go/middleware"
	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/service"
	"github.com/Kevin-Jii/tower-go/utils/excelxml"
	"github.com/Kevin-Jii/tower-go/utils/http"
	"github.com/gin-gonic/gin"
)

type B2BController struct {
	service *service.B2BService
}

func NewB2BController(s *service.B2BService) *B2BController {
	return &B2BController{service: s}
}

func (c *B2BController) CreateCustomer(ctx *gin.Context) {
	var req model.CreateB2BCustomerReq
	if !http.BindJSON(ctx, &req) {
		return
	}
	storeID := middleware.GetStoreID(ctx)
	if middleware.HQUnboundAdmin(ctx) && req.StoreID > 0 {
		storeID = req.StoreID
	}
	customer, err := c.service.CreateCustomer(storeID, &req)
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, customer)
}

func (c *B2BController) UpdateCustomer(ctx *gin.Context) {
	id, ok := http.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	var req model.UpdateB2BCustomerReq
	if !http.BindJSON(ctx, &req) {
		return
	}
	customer, err := c.service.UpdateCustomer(id, middleware.GetStoreID(ctx), middleware.HQUnboundAdmin(ctx), &req)
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, customer)
}

func (c *B2BController) ListCustomers(ctx *gin.Context) {
	var req model.ListB2BCustomerReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		http.Error(ctx, 400, err.Error())
		return
	}
	if !middleware.HQUnboundAdmin(ctx) {
		req.StoreID = middleware.GetStoreID(ctx)
	}
	rows, total, err := c.service.ListCustomers(ctx.Request.Context(), &req)
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.SuccessWithPagination(ctx, rows, total, req.Page, req.PageSize)
}

func (c *B2BController) UpsertPrice(ctx *gin.Context) {
	var req model.UpsertB2BPriceReq
	if !http.BindJSON(ctx, &req) {
		return
	}
	if err := c.service.UpsertPrice(middleware.GetStoreID(ctx), middleware.HQUnboundAdmin(ctx), &req); err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, nil)
}

func (c *B2BController) ListPrices(ctx *gin.Context) {
	var req model.ListB2BPriceReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		http.Error(ctx, 400, err.Error())
		return
	}
	if !middleware.HQUnboundAdmin(ctx) {
		req.StoreID = middleware.GetStoreID(ctx)
	}
	rows, total, err := c.service.ListPrices(&req)
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.SuccessWithPagination(ctx, rows, total, req.Page, req.PageSize)
}

func (c *B2BController) DeletePrice(ctx *gin.Context) {
	id, ok := http.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	if err := c.service.DeletePrice(id); err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, nil)
}

func (c *B2BController) CreateSupplyOrder(ctx *gin.Context) {
	var req model.CreateB2BSupplyOrderReq
	if !http.BindJSON(ctx, &req) {
		return
	}
	order, err := c.service.CreateSupplyOrder(middleware.GetStoreID(ctx), middleware.GetUserID(ctx), middleware.HQUnboundAdmin(ctx), &req)
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, order)
}

func (c *B2BController) ListSupplyOrders(ctx *gin.Context) {
	var req model.ListB2BSupplyOrderReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		http.Error(ctx, 400, err.Error())
		return
	}
	if !middleware.HQUnboundAdmin(ctx) {
		req.StoreID = middleware.GetStoreID(ctx)
	}
	rows, total, err := c.service.ListSupplyOrders(&req)
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.SuccessWithPagination(ctx, rows, total, req.Page, req.PageSize)
}

func (c *B2BController) ExportSupplyOrders(ctx *gin.Context) {
	var req model.ListB2BSupplyOrderReq
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
		req.StoreID = middleware.GetStoreID(ctx)
	}

	list, _, err := c.service.ListSupplyOrders(&req)
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	rows := make([][]interface{}, 0, len(list))
	for _, item := range list {
		rows = append(rows, []interface{}{
			item.OrderDate,
			item.OrderNo,
			item.CustomerName,
			b2bDeliveryLabel(item.DeliveryStatus),
			b2bPaymentLabel(item.PaymentStatus),
			formatAmount(item.TotalAmount),
			formatAmount(item.PaidAmount),
			formatAmount(item.UnpaidAmount),
			formatAmount(item.CostAmount),
			formatAmount(item.ProfitAmount),
			b2bItemsText(item.Items),
			item.OperatorName,
			item.Remark,
			item.CreatedAt,
		})
	}

	data := excelxml.Build([]excelxml.Sheet{{
		Name:    "B2B供货",
		Headers: []string{"供货日期", "供货单号", "客户", "配送状态", "收款状态", "应收金额", "已收金额", "未收金额", "成本金额", "毛利金额", "商品明细", "操作人", "备注", "创建时间"},
		Rows:    rows,
	}})
	http.File(ctx, data, excelxml.Filename("b2b-supply-orders-"+date))
}

func (c *B2BController) GetSupplyOrder(ctx *gin.Context) {
	id, ok := http.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	order, err := c.service.GetSupplyOrder(id, middleware.GetStoreID(ctx), middleware.HQUnboundAdmin(ctx))
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, order)
}

func (c *B2BController) UpdateSupplyOrderDelivery(ctx *gin.Context) {
	id, ok := http.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	var req model.UpdateB2BSupplyOrderDeliveryReq
	if !http.BindJSON(ctx, &req) {
		return
	}
	order, err := c.service.UpdateSupplyOrderDelivery(id, middleware.GetStoreID(ctx), middleware.HQUnboundAdmin(ctx), &req)
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, order)
}

func (c *B2BController) UpdateSupplyOrderPayment(ctx *gin.Context) {
	id, ok := http.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	var req model.UpdateB2BSupplyOrderPaymentReq
	if !http.BindJSON(ctx, &req) {
		return
	}
	order, err := c.service.UpdateSupplyOrderPayment(id, middleware.GetStoreID(ctx), middleware.HQUnboundAdmin(ctx), &req)
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, order)
}
