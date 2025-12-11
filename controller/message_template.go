package controller

import (
	"github.com/Kevin-Jii/tower-go/middleware"
	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/service"
	"github.com/gin-gonic/gin"

	httpPkg "github.com/Kevin-Jii/tower-go/utils/http"
)

type MessageTemplateController struct {
	svc *service.MessageTemplateService
}

func NewMessageTemplateController(svc *service.MessageTemplateService) *MessageTemplateController {
	return &MessageTemplateController{svc: svc}
}

// List godoc
// @Summary 获取消息模板列表
// @Tags 消息模板
// @Produce json
// @Security Bearer
// @Success 200 {object} http.Response{data=[]model.MessageTemplate}
// @Router /message-templates [get]
func (c *MessageTemplateController) List(ctx *gin.Context) {
	if !middleware.IsAdmin(ctx) {
		httpPkg.Error(ctx, 403, "only admin can view templates")
		return
	}

	templates, err := c.svc.List()
	if err != nil {
		httpPkg.Error(ctx, 500, err.Error())
		return
	}

	httpPkg.Success(ctx, templates)
}

// Get godoc
// @Summary 获取消息模板详情
// @Tags 消息模板
// @Produce json
// @Security Bearer
// @Param id path int true "模板ID"
// @Success 200 {object} http.Response{data=model.MessageTemplate}
// @Router /message-templates/{id} [get]
func (c *MessageTemplateController) Get(ctx *gin.Context) {
	if !middleware.IsAdmin(ctx) {
		httpPkg.Error(ctx, 403, "only admin can view templates")
		return
	}

	id, ok := httpPkg.ParseUintParam(ctx, "id")
	if !ok {
		return
	}

	template, err := c.svc.GetByID(id)
	if err != nil {
		httpPkg.Error(ctx, 404, "template not found")
		return
	}

	httpPkg.Success(ctx, template)
}

// Create godoc
// @Summary 创建消息模板
// @Tags 消息模板
// @Accept json
// @Produce json
// @Security Bearer
// @Param template body model.CreateMessageTemplateReq true "模板信息"
// @Success 200 {object} http.Response{data=model.MessageTemplate}
// @Router /message-templates [post]
func (c *MessageTemplateController) Create(ctx *gin.Context) {
	if !middleware.IsAdmin(ctx) {
		httpPkg.Error(ctx, 403, "only admin can create templates")
		return
	}

	var req model.CreateMessageTemplateReq
	if !httpPkg.BindJSON(ctx, &req) {
		return
	}

	template, err := c.svc.Create(&req)
	if err != nil {
		httpPkg.Error(ctx, 500, err.Error())
		return
	}

	httpPkg.Success(ctx, template)
}

// Update godoc
// @Summary 更新消息模板
// @Tags 消息模板
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "模板ID"
// @Param template body model.UpdateMessageTemplateReq true "模板信息"
// @Success 200 {object} http.Response
// @Router /message-templates/{id} [put]
func (c *MessageTemplateController) Update(ctx *gin.Context) {
	if !middleware.IsAdmin(ctx) {
		httpPkg.Error(ctx, 403, "only admin can update templates")
		return
	}

	id, ok := httpPkg.ParseUintParam(ctx, "id")
	if !ok {
		return
	}

	var req model.UpdateMessageTemplateReq
	if !httpPkg.BindJSON(ctx, &req) {
		return
	}

	if err := c.svc.Update(id, &req); err != nil {
		httpPkg.Error(ctx, 500, err.Error())
		return
	}

	httpPkg.Success(ctx, gin.H{"message": "updated"})
}

// Delete godoc
// @Summary 删除消息模板
// @Tags 消息模板
// @Produce json
// @Security Bearer
// @Param id path int true "模板ID"
// @Success 200 {object} http.Response
// @Router /message-templates/{id} [delete]
func (c *MessageTemplateController) Delete(ctx *gin.Context) {
	if !middleware.IsAdmin(ctx) {
		httpPkg.Error(ctx, 403, "only admin can delete templates")
		return
	}

	id, ok := httpPkg.ParseUintParam(ctx, "id")
	if !ok {
		return
	}

	if err := c.svc.Delete(id); err != nil {
		httpPkg.Error(ctx, 500, err.Error())
		return
	}

	httpPkg.Success(ctx, gin.H{"message": "deleted"})
}
