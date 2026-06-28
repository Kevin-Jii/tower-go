package controller

import (
	"github.com/Kevin-Jii/tower-go/middleware"
	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/service"
	"github.com/Kevin-Jii/tower-go/utils/excelxml"
	"github.com/Kevin-Jii/tower-go/utils/http"
	"github.com/gin-gonic/gin"
)

type StoreExpenseController struct {
	storeExpenseService *service.StoreExpenseService
}

func NewStoreExpenseController(storeExpenseService *service.StoreExpenseService) *StoreExpenseController {
	return &StoreExpenseController{storeExpenseService: storeExpenseService}
}

func (c *StoreExpenseController) Create(ctx *gin.Context) {
	var req model.CreateStoreExpenseReq
	if !http.BindJSON(ctx, &req) {
		return
	}
	record, err := c.storeExpenseService.Create(middleware.GetStoreID(ctx), middleware.GetUserID(ctx), &req, middleware.HQUnboundAdmin(ctx))
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, record)
}

func (c *StoreExpenseController) List(ctx *gin.Context) {
	var req model.ListStoreExpenseReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		http.Error(ctx, 400, err.Error())
		return
	}
	if !middleware.HQUnboundAdmin(ctx) {
		req.StoreID = middleware.GetStoreID(ctx)
	}
	list, total, err := c.storeExpenseService.List(ctx.Request.Context(), &req)
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.SuccessWithPagination(ctx, list, total, req.Page, req.PageSize)
}

func (c *StoreExpenseController) Export(ctx *gin.Context) {
	var req model.ListStoreExpenseReq
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

	list, _, err := c.storeExpenseService.List(ctx.Request.Context(), &req)
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	rows := make([][]interface{}, 0, len(list))
	for _, item := range list {
		rows = append(rows, []interface{}{
			item.ExpenseDate,
			item.ExpenseNo,
			storeName(item.Store),
			item.CategoryName,
			formatAmount(item.Amount),
			userName(item.Operator, item.OperatorName),
			item.Remark,
			item.CreatedAt,
		})
	}
	data := excelxml.Build([]excelxml.Sheet{{
		Name:    "门店支出",
		Headers: []string{"支出日期", "支出单号", "门店", "分类", "金额", "操作人", "备注", "创建时间"},
		Rows:    rows,
	}})
	http.File(ctx, data, excelxml.Filename("store-expenses-"+date))
}

func (c *StoreExpenseController) Get(ctx *gin.Context) {
	id, ok := http.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	record, err := c.storeExpenseService.Get(id, middleware.GetStoreID(ctx), middleware.HQUnboundAdmin(ctx))
	if err != nil {
		http.Error(ctx, 404, "未找到该支出记录")
		return
	}
	http.Success(ctx, record)
}

func (c *StoreExpenseController) Update(ctx *gin.Context) {
	id, ok := http.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	var req model.UpdateStoreExpenseReq
	if !http.BindJSON(ctx, &req) {
		return
	}
	record, err := c.storeExpenseService.Update(id, middleware.GetStoreID(ctx), middleware.GetUserID(ctx), &req, middleware.HQUnboundAdmin(ctx))
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, record)
}

func (c *StoreExpenseController) Delete(ctx *gin.Context) {
	id, ok := http.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	if err := c.storeExpenseService.Delete(id, middleware.GetStoreID(ctx), middleware.HQUnboundAdmin(ctx)); err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, nil)
}

func (c *StoreExpenseController) Stats(ctx *gin.Context) {
	var req model.ListStoreExpenseReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		http.Error(ctx, 400, err.Error())
		return
	}
	if !middleware.HQUnboundAdmin(ctx) {
		req.StoreID = middleware.GetStoreID(ctx)
	}
	stats, err := c.storeExpenseService.Stats(&req)
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, stats)
}
