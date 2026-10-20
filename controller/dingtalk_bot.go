package controller

import (
	"github.com/Kevin-Jii/tower-go/middleware"
	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/pkg/apicode"
	"github.com/Kevin-Jii/tower-go/service"

	"github.com/gin-gonic/gin"

	httpPkg "github.com/Kevin-Jii/tower-go/utils/http"
)

type DingTalkBotController struct {
	svc *service.DingTalkService
}

func NewDingTalkBotController(svc *service.DingTalkService) *DingTalkBotController {
	return &DingTalkBotController{svc: svc}
}

// CreateBot godoc
// @Summary 创建钉钉机器人配置
// @Tags 钉钉机器人
// @Accept json
// @Produce json
// @Security Bearer
// @Param bot body model.CreateDingTalkBotReq true "机器人配置"
// @Success 200 {object} http.Response{data=model.DingTalkBot}
// @Router /dingtalk-bots [post]
func (c *DingTalkBotController) CreateBot(ctx *gin.Context) {
	if !middleware.IsAdmin(ctx) {
		httpPkg.ErrorApp(ctx, apicode.OperationDenied)
		return
	}

	var req model.CreateDingTalkBotReq
	if !httpPkg.BindJSON(ctx, &req) {
		return
	}

	bot, err := c.svc.CreateBot(&req)
	if err != nil {
		switch err.Error() {
		case "webhook already exists":
			httpPkg.ErrorApp(ctx, apicode.Conflict)
		case "webhook is required for webhook type",
			"clientID and clientSecret are required for stream type",
			"card msg type only supports stream bot",
			"card_msg_key is required when msg_type is card":
			httpPkg.ErrorApp(ctx, apicode.InvalidParameter)
		default:
			httpPkg.ErrorFrom(ctx, err)
		}
		return
	}

	httpPkg.Success(ctx, bot)
}

// GetBot godoc
// @Summary 获取钉钉机器人详情
// @Tags 钉钉机器人
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "机器人ID"
// @Success 200 {object} http.Response{data=model.DingTalkBot}
// @Router /dingtalk-bots/{id} [get]
func (c *DingTalkBotController) GetBot(ctx *gin.Context) {
	if !middleware.IsAdmin(ctx) {
		httpPkg.ErrorApp(ctx, apicode.OperationDenied)
		return
	}

	id, ok := httpPkg.ParseUintParam(ctx, "id")
	if !ok {
		return
	}

	bot, err := c.svc.GetBot(id)
	if err != nil {
		httpPkg.ErrorFrom(ctx, err)
		return
	}

	httpPkg.Success(ctx, bot)
}

// ListBots godoc
// @Summary 获取钉钉机器人列表
// @Tags 钉钉机器人
// @Accept json
// @Produce json
// @Security Bearer
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} http.Response{data=[]model.DingTalkBot}
// @Router /dingtalk-bots [get]
func (c *DingTalkBotController) ListBots(ctx *gin.Context) {
	if !middleware.IsAdmin(ctx) {
		httpPkg.ErrorApp(ctx, apicode.OperationDenied)
		return
	}

	page := httpPkg.GetPage(ctx)
	pageSize := httpPkg.GetPageSize(ctx)

	bots, total, err := c.svc.ListBots(page, pageSize)
	if err != nil {
		httpPkg.ErrorFrom(ctx, err)
		return
	}

	// 自定义返回结构，去掉嵌套 store，直接返回 store_code 和 store_name
	var result []map[string]interface{}
	for _, bot := range bots {
		m := httpPkg.StructToMap(bot)
		if bot.Store != nil {
			m["store_code"] = bot.Store.StoreCode
			m["store_name"] = bot.Store.Name
		} else {
			m["store_code"] = ""
			m["store_name"] = ""
		}
		// 移除嵌套 store 字段
		delete(m, "store")
		result = append(result, m)
	}
	httpPkg.SuccessWithPagination(ctx, result, total, page, pageSize)
}

// UpdateBot godoc
// @Summary 更新钉钉机器人配置
// @Tags 钉钉机器人
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "机器人ID"
// @Param bot body model.UpdateDingTalkBotReq true "机器人配置"
// @Success 200 {object} http.Response
// @Router /dingtalk-bots/{id} [put]
func (c *DingTalkBotController) UpdateBot(ctx *gin.Context) {
	if !middleware.IsAdmin(ctx) {
		httpPkg.ErrorApp(ctx, apicode.OperationDenied)
		return
	}

	id, ok := httpPkg.ParseUintParam(ctx, "id")
	if !ok {
		return
	}

	var req model.UpdateDingTalkBotReq
	if !httpPkg.BindJSON(ctx, &req) {
		return
	}

	if err := c.svc.UpdateBot(id, &req); err != nil {
		switch err.Error() {
		case "webhook already exists":
			httpPkg.ErrorApp(ctx, apicode.Conflict)
		case "no fields to update",
			"card msg type only supports stream bot",
			"card_msg_key is required when msg_type is card":
			httpPkg.ErrorApp(ctx, apicode.InvalidParameter)
		default:
			httpPkg.ErrorFrom(ctx, err)
		}
		return
	}

	// 返回最新机器人详情，保证 robot_code 字段
	updatedBot, err := c.svc.GetBot(id)
	if err != nil {
		httpPkg.ErrorFrom(ctx, err)
		return
	}
	httpPkg.Success(ctx, updatedBot)
}

// DeleteBot godoc
// @Summary 删除钉钉机器人配置
// @Tags 钉钉机器人
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "机器人ID"
// @Success 200 {object} http.Response
// @Router /dingtalk-bots/{id} [delete]
func (c *DingTalkBotController) DeleteBot(ctx *gin.Context) {
	if !middleware.IsAdmin(ctx) {
		httpPkg.ErrorApp(ctx, apicode.OperationDenied)
		return
	}

	id, ok := httpPkg.ParseUintParam(ctx, "id")
	if !ok {
		return
	}

	if err := c.svc.DeleteBot(id); err != nil {
		httpPkg.ErrorFrom(ctx, err)
		return
	}

	// 删除后返回空对象，结构一致
	httpPkg.Success(ctx, gin.H{})
}

// TestBot godoc
// @Summary 测试钉钉机器人连接
// @Tags 钉钉机器人
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "机器人ID"
// @Success 200 {object} http.Response
// @Router /dingtalk-bots/{id}/test [post]
func (c *DingTalkBotController) TestBot(ctx *gin.Context) {
	if !middleware.IsAdmin(ctx) {
		httpPkg.ErrorApp(ctx, apicode.OperationDenied)
		return
	}

	id, ok := httpPkg.ParseUintParam(ctx, "id")
	if !ok {
		return
	}

	// 先获取机器人配置并做基础校验（避免把错误请求发送到钉钉）
	bot, err := c.svc.GetBot(id)
	if err != nil {
		httpPkg.ErrorFrom(ctx, err)
		return
	}

	if bot.BotType == "stream" && bot.RobotCode == "" {
		httpPkg.ErrorApp(ctx, apicode.InvalidParameter)
		return
	}

	if err := c.svc.TestBot(id); err != nil {
		httpPkg.ErrorFrom(ctx, err)
		return
	}

	// 返回机器人详情，保证 robot_code 字段
	httpPkg.Success(ctx, gin.H{"message": "测试消息发送成功", "robot_code": bot.RobotCode})
}

// TestStreamBotCallback godoc
// @Summary 测试 Stream 机器人接收回调
// @Tags 钉钉机器人
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "机器人ID"
// @Success 200 {object} http.Response
// @Router /dingtalk-bots/{id}/test-callback [post]
func (c *DingTalkBotController) TestStreamBotCallback(ctx *gin.Context) {
	if !middleware.IsAdmin(ctx) {
		httpPkg.ErrorApp(ctx, apicode.OperationDenied)
		return
	}

	id, ok := httpPkg.ParseUintParam(ctx, "id")
	if !ok {
		return
	}

	// 获取机器人配置
	bot, err := c.svc.GetBot(id)
	if err != nil {
		httpPkg.ErrorFrom(ctx, err)
		return
	}

	if bot.BotType != "stream" {
		httpPkg.ErrorApp(ctx, apicode.InvalidParameter)
		return
	}

	// 检查 Stream 客户端是否运行
	streamClient := c.svc.GetStreamClient()
	if !streamClient.IsRunning() {
		httpPkg.ErrorApp(ctx, apicode.ExternalServiceFailed)
		return
	}

	_, exists := streamClient.GetClient(id)
	if !exists {
		httpPkg.ErrorApp(ctx, apicode.ExternalServiceFailed)
		return
	}

	httpPkg.Success(ctx, gin.H{
		"message":     "Stream 机器人回调处理器已启用",
		"bot_id":      id,
		"bot_name":    bot.Name,
		"client_type": "stream",
		"status":      "connected",
		"note":        "现在可以在钉钉中向该机器人发送消息以测试回调",
	})
}
