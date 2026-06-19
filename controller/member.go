package controller

import (
	"strconv"

	"github.com/Kevin-Jii/tower-go/middleware"
	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/service"
	"github.com/Kevin-Jii/tower-go/utils/http"

	"github.com/gin-gonic/gin"
)

// MemberController 会员控制器
type MemberController struct {
	service *service.MemberService
}

// NewMemberController 创建会员控制器
func NewMemberController(s *service.MemberService) *MemberController {
	return &MemberController{service: s}
}

// CreateMember 创建会员
// @Summary 创建会员
// @Description 新增一个会员
// @Tags 会员管理
// @Accept json
// @Produce json
// @Param data body model.CreateMemberReq true "会员信息"
// @Success 200 {object} http.Response{data=model.Member}
// @Failure 400 {object} map[string]string
// @Router /members [post]
func (c *MemberController) CreateMember(ctx *gin.Context) {
	var req model.CreateMemberReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		http.Error(ctx, 400, err.Error())
		return
	}
	storeID := middleware.GetStoreID(ctx)
	member, err := c.service.CreateMember(&req, storeID)
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, member)
}

// UpdateMember 更新会员
// @Summary 更新会员
// @Description 根据ID更新会员信息，可局部更新
// @Tags 会员管理
// @Accept json
// @Produce json
// @Param id path int true "会员ID"
// @Param data body model.UpdateMemberReq true "仅传需要修改的字段"
// @Success 200 {object} http.Response{data=model.Member}
// @Failure 400 {object} map[string]string
// @Router /members/{id} [put]
func (c *MemberController) UpdateMember(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(ctx, 400, "invalid id")
		return
	}
	var req model.UpdateMemberReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		http.Error(ctx, 400, err.Error())
		return
	}
	storeID := middleware.GetStoreID(ctx)
	member, err := c.service.UpdateMember(uint(id), &req, storeID, middleware.HQUnboundAdmin(ctx))
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, member)
}

// DeleteMember 删除会员
// @Summary 删除会员
// @Description 根据ID删除会员
// @Tags 会员管理
// @Produce json
// @Param id path int true "会员ID"
// @Success 200 {object} http.Response
// @Failure 400 {object} map[string]string
// @Router /members/{id} [delete]
func (c *MemberController) DeleteMember(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(ctx, 400, "invalid id")
		return
	}
	storeID := middleware.GetStoreID(ctx)
	if err := c.service.DeleteMember(uint(id), storeID, middleware.HQUnboundAdmin(ctx)); err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, nil)
}

// GetMember 获取会员详情
// @Summary 获取会员详情
// @Description 根据ID获取会员信息
// @Tags 会员管理
// @Produce json
// @Param id path int true "会员ID"
// @Success 200 {object} http.Response{data=model.Member}
// @Failure 404 {object} map[string]string
// @Router /members/{id} [get]
func (c *MemberController) GetMember(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(ctx, 400, "invalid id")
		return
	}
	storeID := middleware.GetStoreID(ctx)
	member, err := c.service.GetMember(uint(id), storeID, middleware.HQUnboundAdmin(ctx))
	if err != nil {
		http.Error(ctx, 404, err.Error())
		return
	}
	http.Success(ctx, member)
}

// ListMembers 获取会员列表
// @Summary 获取会员列表
// @Description 获取会员列表，支持关键字模糊查询手机号/UID，支持分页
// @Tags 会员管理
// @Produce json
// @Param keyword query string false "关键字(模糊匹配手机号/UID)"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} http.Response{data=[]model.Member}
// @Router /members [get]
func (c *MemberController) ListMembers(ctx *gin.Context) {
	keyword := ctx.Query("keyword")
	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("page_size", "20")

	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	storeID := middleware.GetStoreID(ctx)
	members, total, err := c.service.ListMembers(keyword, page, pageSize, storeID, middleware.HQUnboundAdmin(ctx))
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.SuccessWithPagination(ctx, members, total, page, pageSize)
}

// ListWineStorages 查询会员存酒
func (c *MemberController) ListWineStorages(ctx *gin.Context) {
	var req model.ListMemberWineStorageReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		http.Error(ctx, 400, err.Error())
		return
	}
	req.Page = http.GetPage(ctx)
	req.PageSize = http.GetPageSize(ctx)

	storeID := middleware.GetStoreID(ctx)
	list, total, err := c.service.ListWineStorages(&req, storeID, middleware.HQUnboundAdmin(ctx))
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.SuccessWithPagination(ctx, list, total, req.Page, req.PageSize)
}

// DepositWine 会员存酒
func (c *MemberController) DepositWine(ctx *gin.Context) {
	var req model.MemberWineAdjustReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		http.Error(ctx, 400, err.Error())
		return
	}
	storeID := middleware.GetStoreID(ctx)
	row, err := c.service.DepositWine(storeID, middleware.GetUserID(ctx), middleware.HQUnboundAdmin(ctx), &req)
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, row)
}

// WithdrawWine 会员取酒
func (c *MemberController) WithdrawWine(ctx *gin.Context) {
	var req model.MemberWineAdjustReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		http.Error(ctx, 400, err.Error())
		return
	}
	storeID := middleware.GetStoreID(ctx)
	row, err := c.service.WithdrawWine(storeID, middleware.GetUserID(ctx), middleware.HQUnboundAdmin(ctx), &req)
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, row)
}

// ListWineTransactions 查询会员存取酒流水
func (c *MemberController) ListWineTransactions(ctx *gin.Context) {
	var req model.ListMemberWineTransactionReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		http.Error(ctx, 400, err.Error())
		return
	}
	req.Page = http.GetPage(ctx)
	req.PageSize = http.GetPageSize(ctx)

	storeID := middleware.GetStoreID(ctx)
	list, total, err := c.service.ListWineTransactions(&req, storeID, middleware.HQUnboundAdmin(ctx))
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.SuccessWithPagination(ctx, list, total, req.Page, req.PageSize)
}

// AdjustBalance 调整余额
// @Summary 调整会员余额
// @Description 使用乐观锁调整会员余额（仅管理员可操作）
// @Tags 会员管理
// @Accept json
// @Produce json
// @Param id path int true "会员ID"
// @Param data body model.AdjustBalanceReq true "调整信息"
// @Success 200 {object} http.Response{data=model.Member}
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /members/{id}/adjust-balance [post]
func (c *MemberController) AdjustBalance(ctx *gin.Context) {
	// 检查是否是管理员
	if !middleware.IsAdmin(ctx) {
		http.Error(ctx, 403, "permission denied: admin only")
		return
	}

	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(ctx, 400, "invalid id")
		return
	}
	var req model.AdjustBalanceReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		http.Error(ctx, 400, err.Error())
		return
	}

	// 获取当前用户的门店ID和用户ID
	storeID := middleware.GetStoreID(ctx)
	userID := middleware.GetUserID(ctx)

	member, err := c.service.AdjustBalance(uint(id), req.Amount, req.Type, req.Remark, req.Version, storeID, userID, middleware.HQUnboundAdmin(ctx))
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, member)
}

// GetMemberByPhone 通过手机号获取会员
// @Summary 通过手机号获取会员
// @Description 根据手机号获取会员信息
// @Tags 会员管理
// @Produce json
// @Param phone query string true "手机号"
// @Success 200 {object} http.Response{data=model.Member}
// @Failure 404 {object} map[string]string
// @Router /members/phone [get]
func (c *MemberController) GetMemberByPhone(ctx *gin.Context) {
	phone := ctx.Query("phone")
	if phone == "" {
		http.Error(ctx, 400, "phone is required")
		return
	}
	storeID := middleware.GetStoreID(ctx)
	member, err := c.service.GetMemberByPhone(phone, storeID, middleware.HQUnboundAdmin(ctx))
	if err != nil {
		http.Error(ctx, 404, err.Error())
		return
	}
	http.Success(ctx, member)
}

// ========== WalletLog 接口 ==========

// ListWalletLogs 查询流水列表
// @Summary 查询流水列表
// @Description 查询会员流水记录
// @Tags 流水管理
// @Produce json
// @Param memberId query int false "会员ID"
// @Param changeType query int false "变动类型 1=充值 2=消费 3=退款 4=调增 5=调减"
// @Param startTime query string false "开始时间"
// @Param endTime query string false "结束时间"
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Success 200 {object} http.Response{data=[]model.WalletLog}
// @Router /wallet-logs [get]
func (c *MemberController) ListWalletLogs(ctx *gin.Context) {
	var req model.ListWalletLogReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		http.Error(ctx, 400, err.Error())
		return
	}
	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("pageSize", "20")
	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	storeID := middleware.GetStoreID(ctx)
	logs, total, err := c.service.ListWalletLogs(&req, page, pageSize, storeID, middleware.HQUnboundAdmin(ctx))
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.SuccessWithPagination(ctx, logs, total, page, pageSize)
}

// ========== RechargeOrder 接口 ==========

// CreateRechargeOrder 创建充值单
// @Summary 创建充值单
// @Description 创建会员充值订单（自动完成支付）
// @Tags 充值管理
// @Accept json
// @Produce json
// @Param data body model.CreateRechargeOrderReq true "充值信息"
// @Success 200 {object} http.Response{data=model.RechargeOrder}
// @Failure 400 {object} map[string]string
// @Router /recharge-orders [post]
func (c *MemberController) CreateRechargeOrder(ctx *gin.Context) {
	var req model.CreateRechargeOrderReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		http.Error(ctx, 400, err.Error())
		return
	}

	// 获取当前用户的门店ID和用户ID
	storeID := middleware.GetStoreID(ctx)
	userID := middleware.GetUserID(ctx)

	order, err := c.service.CreateRechargeOrder(&req, storeID, userID, middleware.HQUnboundAdmin(ctx))
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, order)
}

// GetRechargeOrder 获取充值单详情
// @Summary 获取充值单详情
// @Description 根据ID或单号获取充值单信息
// @Tags 充值管理
// @Produce json
// @Param id path int true "充值单ID"
// @Success 200 {object} http.Response{data=model.RechargeOrder}
// @Failure 404 {object} map[string]string
// @Router /recharge-orders/{id} [get]
func (c *MemberController) GetRechargeOrder(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(ctx, 400, "invalid id")
		return
	}
	storeID := middleware.GetStoreID(ctx)
	order, err := c.service.GetRechargeOrder(uint(id), storeID, middleware.HQUnboundAdmin(ctx))
	if err != nil {
		http.Error(ctx, 404, err.Error())
		return
	}
	http.Success(ctx, order)
}

// ListRechargeOrders 查询充值单列表
// @Summary 查询充值单列表
// @Description 查询会员充值记录
// @Tags 充值管理
// @Produce json
// @Param memberId query int false "会员ID"
// @Param status query int false "支付状态 0=待支付 1=已支付 2=已取消 3=已退款"
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Success 200 {object} http.Response{data=[]model.RechargeOrder}
// @Router /recharge-orders [get]
func (c *MemberController) ListRechargeOrders(ctx *gin.Context) {
	memberIDStr := ctx.Query("memberId")
	var memberID uint
	if memberIDStr != "" {
		id, err := strconv.Atoi(memberIDStr)
		if err == nil {
			memberID = uint(id)
		}
	}
	statusStr := ctx.Query("status")
	var status *model.PayStatusEnum
	if statusStr != "" {
		s := model.PayStatusEnum(statusStr[0] - '0')
		status = &s
	}
	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("pageSize", "20")
	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	storeID := middleware.GetStoreID(ctx)
	orders, total, err := c.service.ListRechargeOrders(memberID, status, page, pageSize, storeID, middleware.HQUnboundAdmin(ctx))
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.SuccessWithPagination(ctx, orders, total, page, pageSize)
}

// PayRechargeOrder 支付充值单
// @Summary 支付充值单
// @Description 支付充值单（余额充值）
// @Tags 充值管理
// @Accept json
// @Produce json
// @Param data body model.PayRechargeOrderReq true "支付信息"
// @Success 200 {object} http.Response{data=model.RechargeOrder}
// @Failure 400 {object} map[string]string
// @Router /recharge-orders/pay [post]
func (c *MemberController) PayRechargeOrder(ctx *gin.Context) {
	var req model.PayRechargeOrderReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		http.Error(ctx, 400, err.Error())
		return
	}

	// 获取当前用户的门店ID和用户ID
	storeID := middleware.GetStoreID(ctx)
	userID := middleware.GetUserID(ctx)

	order, err := c.service.PayRechargeOrder(req.OrderNo, storeID, userID, middleware.HQUnboundAdmin(ctx))
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, order)
}

// CancelRechargeOrder 取消充值单
// @Summary 取消充值单
// @Description 取消待支付的充值单
// @Tags 充值管理
// @Accept json
// @Produce json
// @Param orderNo path string true "充值单号"
// @Success 200 {object} http.Response{data=model.RechargeOrder}
// @Failure 400 {object} map[string]string
// @Router /recharge-orders/{orderNo}/cancel [post]
func (c *MemberController) CancelRechargeOrder(ctx *gin.Context) {
	orderNo := ctx.Param("orderNo")
	if orderNo == "" {
		http.Error(ctx, 400, "orderNo is required")
		return
	}
	storeID := middleware.GetStoreID(ctx)
	order, err := c.service.CancelRechargeOrder(orderNo, storeID, middleware.HQUnboundAdmin(ctx))
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, order)
}

// ListMemberConsumptions 查询会员消费记录
// @Summary 查询会员消费记录
// @Description 按会员查询门店记账消费记录，支持日期筛选与汇总
// @Tags 会员管理
// @Produce json
// @Param id path int true "会员ID"
// @Param start_date query string false "开始日期(YYYY-MM-DD)"
// @Param end_date query string false "结束日期(YYYY-MM-DD)"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} http.Response
// @Router /members/{id}/consumptions [get]
func (c *MemberController) ListMemberConsumptions(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(ctx, 400, "invalid id")
		return
	}

	startDate := ctx.Query("start_date")
	endDate := ctx.Query("end_date")
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	storeID := middleware.GetStoreID(ctx)
	list, total, summary, err := c.service.ListMemberConsumptions(uint(id), startDate, endDate, page, pageSize, storeID, middleware.HQUnboundAdmin(ctx))
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	http.Success(ctx, gin.H{
		"list":      list,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
		"summary":   summary,
	})
}
