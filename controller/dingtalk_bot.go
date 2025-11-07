package controller

import (
	"net/http"
	"tower-go/middleware"
	"tower-go/model"
	"tower-go/service"

	"github.com/gin-gonic/gin"

	"tower-go/utils"
)

type DingTalkBotController struct {
	svc *service.DingTalkService
}

func NewDingTalkBotController(svc *service.DingTalkService) *DingTalkBotController {
	return &DingTalkBotController{svc: svc}
}

// CreateBot godoc
// @Summary 创建钉钉机器人配置
// @Tags dingtalk-bots
// @Accept json
// @Produce json
// @Security Bearer
// @Param bot body model.CreateDingTalkBotReq true "机器人配置"
// @Success 200 {object} utils.StandardResponse{data=model.DingTalkBot}
// @Router /dingtalk-bots [post]
func (c *DingTalkBotController) CreateBot(ctx *gin.Context) {
	if !middleware.IsAdmin(ctx) {
		utils.Error(ctx, http.StatusForbidden, "only admin can create bot")
		return
	}

	var req model.CreateDingTalkBotReq
	if !utils.BindJSON(ctx, &req) {
		return
	}

	bot, err := c.svc.CreateBot(&req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "webhook already exists" {
			status = http.StatusConflict
		} else if err.Error() == "webhook is required for webhook type" ||
			err.Error() == "clientID and clientSecret are required for stream type" {
			status = http.StatusBadRequest
		}
		utils.Error(ctx, status, err.Error())
		return
	}

	utils.Success(ctx, bot)
}

// GetBot godoc
// @Summary 获取钉钉机器人详情
// @Tags dingtalk-bots
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "机器人ID"
// @Success 200 {object} utils.StandardResponse{data=model.DingTalkBot}
// @Router /dingtalk-bots/{id} [get]
func (c *DingTalkBotController) GetBot(ctx *gin.Context) {
	if !middleware.IsAdmin(ctx) {
		utils.Error(ctx, http.StatusForbidden, "only admin can view bots")
		return
	}

	id, ok := utils.ParseUintParam(ctx, "id")
	if !ok {
		return
	}

	bot, err := c.svc.GetBot(id)
	if err != nil {
		utils.Error(ctx, http.StatusNotFound, "bot not found")
		return
	}

	utils.Success(ctx, bot)
}

// ListBots godoc
// @Summary 获取钉钉机器人列表
// @Tags dingtalk-bots
// @Accept json
// @Produce json
// @Security Bearer
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} utils.StandardResponse{data=[]model.DingTalkBot}
// @Router /dingtalk-bots [get]
func (c *DingTalkBotController) ListBots(ctx *gin.Context) {
	if !middleware.IsAdmin(ctx) {
		utils.Error(ctx, http.StatusForbidden, "only admin can view bots")
		return
	}

	page := utils.GetPage(ctx)
	pageSize := utils.GetPageSize(ctx)

	bots, _, err := c.svc.ListBots(page, pageSize)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	// 自定义返回结构，去掉嵌套 store，直接返回 store_code 和 store_name
	var result []map[string]interface{}
	for _, bot := range bots {
		m := utils.StructToMap(bot)
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
	utils.Success(ctx, result)
}

// UpdateBot godoc
// @Summary 更新钉钉机器人配置
// @Tags dingtalk-bots
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "机器人ID"
// @Param bot body model.UpdateDingTalkBotReq true "机器人配置"
// @Success 200 {object} utils.StandardResponse
// @Router /dingtalk-bots/{id} [put]
func (c *DingTalkBotController) UpdateBot(ctx *gin.Context) {
	if !middleware.IsAdmin(ctx) {
		utils.Error(ctx, http.StatusForbidden, "only admin can update bot")
		return
	}

	id, ok := utils.ParseUintParam(ctx, "id")
	if !ok {
		return
	}

	var req model.UpdateDingTalkBotReq
	if !utils.BindJSON(ctx, &req) {
		return
	}

	if err := c.svc.UpdateBot(id, &req); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "webhook already exists" {
			status = http.StatusConflict
		} else if err.Error() == "no fields to update" {
			status = http.StatusBadRequest
		}
		utils.Error(ctx, status, err.Error())
		return
	}

	// 返回最新机器人详情，保证 robot_code 字段
	updatedBot, err := c.svc.GetBot(id)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "failed to fetch updated bot")
		return
	}
	utils.Success(ctx, updatedBot)
}

// DeleteBot godoc
// @Summary 删除钉钉机器人配置
// @Tags dingtalk-bots
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "机器人ID"
// @Success 200 {object} utils.StandardResponse
// @Router /dingtalk-bots/{id} [delete]
func (c *DingTalkBotController) DeleteBot(ctx *gin.Context) {
	if !middleware.IsAdmin(ctx) {
		utils.Error(ctx, http.StatusForbidden, "only admin can delete bot")
		return
	}

	id, ok := utils.ParseUintParam(ctx, "id")
	if !ok {
		return
	}

	if err := c.svc.DeleteBot(id); err != nil {
		utils.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	// 删除后返回空对象，结构一致
	utils.Success(ctx, gin.H{})
}

// TestBot godoc
// @Summary 测试钉钉机器人连接
// @Tags dingtalk-bots
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "机器人ID"
// @Success 200 {object} utils.StandardResponse
// @Router /dingtalk-bots/{id}/test [post]
func (c *DingTalkBotController) TestBot(ctx *gin.Context) {
	if !middleware.IsAdmin(ctx) {
		utils.Error(ctx, http.StatusForbidden, "only admin can test bot")
		return
	}

	id, ok := utils.ParseUintParam(ctx, "id")
	if !ok {
		return
	}

	// 先获取机器人配置并做基础校验（避免把错误请求发送到钉钉）
	bot, err := c.svc.GetBot(id)
	if err != nil {
		utils.Error(ctx, http.StatusNotFound, "bot not found")
		return
	}

	if bot.BotType == "stream" && bot.RobotCode == "" {
		utils.Error(ctx, http.StatusBadRequest, "robot_code is empty for stream bot: please set the correct robot_code (钉钉分配的机器人编码) in bot configuration before testing")
		return
	}

	if err := c.svc.TestBot(id); err != nil {
		utils.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	// 返回机器人详情，保证 robot_code 字段
	utils.Success(ctx, gin.H{"message": "test message sent successfully", "robot_code": bot.RobotCode})
}
