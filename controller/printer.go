package controller

import (
	"github.com/Kevin-Jii/tower-go/middleware"
	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/service"
	"github.com/Kevin-Jii/tower-go/utils/http"

	"github.com/gin-gonic/gin"
)

type PrinterController struct {
	printerService *service.PrinterService
}

func NewPrinterController(printerService *service.PrinterService) *PrinterController {
	return &PrinterController{printerService: printerService}
}

// BindPrinter godoc
// @Summary 绑定打印机到门店
// @Description 将打印机绑定到指定门店，同时推送到芯烨云
// @Tags 打印机管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param body body model.BindPrinterReq true "打印机信息"
// @Success 200 {object} http.Response
// @Router /printers/bind [post]
func (c *PrinterController) BindPrinter(ctx *gin.Context) {
	// 权限校验
	if !http.RequireAdmin(ctx) {
		return
	}

	var req model.BindPrinterReq
	if !http.BindJSON(ctx, &req) {
		return
	}

	if err := c.printerService.BindPrinter(&req); err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	http.Success(ctx, nil)
}

// UnbindPrinter godoc
// @Summary 解绑打印机
// @Description 从门店解绑打印机，同时从芯烨云移除
// @Tags 打印机管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "打印机ID"
// @Success 200 {object} http.Response
// @Router /printers/{id} [delete]
func (c *PrinterController) UnbindPrinter(ctx *gin.Context) {
	if !http.RequireAdmin(ctx) {
		return
	}

	id, ok := http.ParseUintParam(ctx, "id")
	if !ok {
		return
	}

	if err := c.printerService.UnbindPrinter(id); err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	http.Success(ctx, nil)
}

// UpdatePrinter godoc
// @Summary 更新打印机信息
// @Description 更新打印机信息
// @Tags 打印机管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "打印机ID"
// @Param body body model.UpdatePrinterReq true "打印机信息"
// @Success 200 {object} http.Response
// @Router /printers/{id} [put]
func (c *PrinterController) UpdatePrinter(ctx *gin.Context) {
	if !http.RequireAdmin(ctx) {
		return
	}

	id, ok := http.ParseUintParam(ctx, "id")
	if !ok {
		return
	}

	var req model.UpdatePrinterReq
	if !http.BindJSON(ctx, &req) {
		return
	}

	if err := c.printerService.UpdatePrinter(id, &req); err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	http.Success(ctx, nil)
}

// GetPrinter godoc
// @Summary 获取打印机详情
// @Description 获取打印机详细信息及在线状态
// @Tags 打印机管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "打印机ID"
// @Success 200 {object} http.Response{data=model.PrinterResp}
// @Router /printers/{id} [get]
func (c *PrinterController) GetPrinter(ctx *gin.Context) {
	id, ok := http.ParseUintParam(ctx, "id")
	if !ok {
		return
	}

	resp, err := c.printerService.GetPrinterWithStatus(id)
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	http.Success(ctx, resp)
}

// ListPrinters godoc
// @Summary 门店打印机列表
// @Description 获取指定门店的打印机列表
// @Tags 打印机管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param store_id query int true "门店ID"
// @Success 200 {object} http.Response{data=[]model.Printer}
// @Router /printers [get]
func (c *PrinterController) ListPrinters(ctx *gin.Context) {
	storeID, ok := http.ParseUintQuery(ctx, "store_id")
	if !ok {
		http.Error(ctx, 400, "store_id is required")
		return
	}

	printers, err := c.printerService.ListPrintersByStore(storeID)
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	http.Success(ctx, printers)
}

// ListAllPrinters godoc
// @Summary 所有打印机列表
// @Description 获取所有打印机（仅管理员）
// @Tags 打印机管理
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} http.Response
// @Router /printers/all [get]
func (c *PrinterController) ListAllPrinters(ctx *gin.Context) {
	if !middleware.IsAdmin(ctx) {
		http.Error(ctx, 403, "admin only")
		return
	}

	printers, total, err := c.printerService.ListAllPrinters()
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	http.SuccessWithPagination(ctx, printers, total, 1, int(total))
}

// GetStoreDefaultPrinter godoc
// @Summary 获取门店默认打印机
// @Description 获取指定门店的默认打印机
// @Tags 打印机管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param store_id query int true "门店ID"
// @Success 200 {object} http.Response{data=model.Printer}
// @Router /printers/default [get]
func (c *PrinterController) GetStoreDefaultPrinter(ctx *gin.Context) {
	storeID, ok := http.ParseUintQuery(ctx, "store_id")
	if !ok {
		http.Error(ctx, 400, "store_id is required")
		return
	}

	printer, err := c.printerService.GetDefaultPrinter(storeID)
	if err != nil {
		http.Error(ctx, 404, "no default printer")
		return
	}

	http.Success(ctx, printer)
}

// QueryPrinterStatus godoc
// @Summary 查询打印机在线状态
// @Description 查询指定打印机的在线状态
// @Tags 打印机管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param sn query string true "打印机SN"
// @Success 200 {object} http.Response
// @Router /printers/status [get]
func (c *PrinterController) QueryPrinterStatus(ctx *gin.Context) {
	sn := ctx.Query("sn")
	if sn == "" {
		http.Error(ctx, 400, "sn is required")
		return
	}

	status, err := c.printerService.QueryPrinterStatus(sn)
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	http.Success(ctx, gin.H{
		"sn":     sn,
		"status": status,
		"online": status == 1,
	})
}

// BatchQueryStatus godoc
// @Summary 批量查询打印机状态
// @Description 批量查询指定门店下所有打印机的状态
// @Tags 打印机管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param store_id query int true "门店ID"
// @Success 200 {object} http.Response
// @Router /printers/status/batch [get]
func (c *PrinterController) BatchQueryStatus(ctx *gin.Context) {
	storeID, ok := http.ParseUintQuery(ctx, "store_id")
	if !ok {
		http.Error(ctx, 400, "store_id is required")
		return
	}

	statuses, err := c.printerService.BatchQueryStatus(storeID)
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	http.Success(ctx, statuses)
}

// TestPrintReq 测试打印请求
type TestPrintReq struct {
	Content string `json:"content" binding:"required"`
	Copies  int    `json:"copies"`
}

// TestPrint godoc
// @Summary 测试打印
// @Description 向指定打印机发送测试打印
// @Tags 打印机管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "打印机ID"
// @Param body body TestPrintReq true "打印内容"
// @Success 200 {object} http.Response
// @Router /printers/{id}/test [post]
func (c *PrinterController) TestPrint(ctx *gin.Context) {
	if !http.RequireAdmin(ctx) {
		return
	}

	id, ok := http.ParseUintParam(ctx, "id")
	if !ok {
		return
	}

	var req TestPrintReq
	if !http.BindJSON(ctx, &req) {
		return
	}

	if req.Copies <= 0 {
		req.Copies = 1
	}

	orderId, err := c.printerService.TestPrint(id, req.Content, req.Copies)
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	http.Success(ctx, gin.H{"order_id": orderId})
}