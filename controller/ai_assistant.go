package controller

import (
	"strconv"

	"github.com/Kevin-Jii/tower-go/middleware"
	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/pkg/apicode"
	"github.com/Kevin-Jii/tower-go/service"
	httpPkg "github.com/Kevin-Jii/tower-go/utils/http"
	"github.com/gin-gonic/gin"
)

type AIAssistantController struct{ svc *service.AIAssistantService }

func NewAIAssistantController(s *service.AIAssistantService) *AIAssistantController {
	return &AIAssistantController{svc: s}
}
func (c *AIAssistantController) GetConfig(ctx *gin.Context) {
	if !middleware.HQUnboundAdmin(ctx) {
		httpPkg.ErrorApp(ctx, apicode.OperationDenied)
		return
	}
	v, e := c.svc.GetConfig()
	if e != nil {
		httpPkg.ErrorFrom(ctx, e)
		return
	}
	httpPkg.Success(ctx, v)
}
func (c *AIAssistantController) SaveConfig(ctx *gin.Context) {
	if !middleware.HQUnboundAdmin(ctx) {
		httpPkg.ErrorApp(ctx, apicode.OperationDenied)
		return
	}
	var req model.AIAssistantConfigRequest
	if !httpPkg.BindJSON(ctx, &req) {
		return
	}
	v, e := c.svc.SaveConfig(&req, middleware.GetUserID(ctx))
	if e != nil {
		httpPkg.ErrorFrom(ctx, e)
		return
	}
	httpPkg.Success(ctx, v)
}
func (c *AIAssistantController) TestConfig(ctx *gin.Context) {
	if !middleware.HQUnboundAdmin(ctx) {
		httpPkg.ErrorApp(ctx, apicode.OperationDenied)
		return
	}
	var req model.AIAssistantConfigRequest
	if !httpPkg.BindJSON(ctx, &req) {
		return
	}
	if e := c.svc.TestConfig(&req); e != nil {
		httpPkg.ErrorFrom(ctx, e)
		return
	}
	httpPkg.Success(ctx, gin.H{"message": "连接成功"})
}
func (c *AIAssistantController) storeID(ctx *gin.Context) uint {
	return middleware.ResolveQueryStoreID(ctx, "store_id")
}
func (c *AIAssistantController) Conversations(ctx *gin.Context) {
	rows, e := c.svc.Conversations(middleware.GetUserID(ctx), c.storeID(ctx))
	if e != nil {
		httpPkg.ErrorFrom(ctx, e)
		return
	}
	httpPkg.Success(ctx, rows)
}
func (c *AIAssistantController) CreateConversation(ctx *gin.Context) {
	row, e := c.svc.CreateConversation(middleware.GetUserID(ctx), c.storeID(ctx))
	if e != nil {
		httpPkg.ErrorFrom(ctx, e)
		return
	}
	httpPkg.Success(ctx, row)
}
func (c *AIAssistantController) Messages(ctx *gin.Context) {
	id, e := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if e != nil {
		httpPkg.ErrorApp(ctx, apicode.InvalidID)
		return
	}
	rows, e := c.svc.Messages(uint(id), middleware.GetUserID(ctx), c.storeID(ctx))
	if e != nil {
		httpPkg.ErrorFrom(ctx, e)
		return
	}
	httpPkg.Success(ctx, rows)
}
func (c *AIAssistantController) Rename(ctx *gin.Context) {
	id, e := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if e != nil {
		httpPkg.ErrorApp(ctx, apicode.InvalidID)
		return
	}
	var req struct {
		Title string `json:"title"`
	}
	if !httpPkg.BindJSON(ctx, &req) {
		return
	}
	if e = c.svc.Rename(uint(id), middleware.GetUserID(ctx), c.storeID(ctx), req.Title); e != nil {
		httpPkg.ErrorFrom(ctx, e)
		return
	}
	httpPkg.Success(ctx, gin.H{"updated": true})
}
func (c *AIAssistantController) Delete(ctx *gin.Context) {
	id, e := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if e != nil {
		httpPkg.ErrorApp(ctx, apicode.InvalidID)
		return
	}
	if e = c.svc.Delete(uint(id), middleware.GetUserID(ctx), c.storeID(ctx)); e != nil {
		httpPkg.ErrorFrom(ctx, e)
		return
	}
	httpPkg.Success(ctx, gin.H{"deleted": true})
}
func (c *AIAssistantController) Chat(ctx *gin.Context) {
	var req model.AIAssistantChatRequest
	if !httpPkg.BindJSON(ctx, &req) {
		return
	}
	sid := middleware.ResolveQueryStoreID(ctx, "store_id")
	if req.StoreID > 0 {
		if !middleware.HQUnboundAdmin(ctx) && req.StoreID != middleware.GetStoreID(ctx) {
			httpPkg.ErrorApp(ctx, apicode.OperationDenied)
			return
		}
		if middleware.HQUnboundAdmin(ctx) {
			sid = req.StoreID
		}
	}
	if sid == 0 && !middleware.HQUnboundAdmin(ctx) {
		httpPkg.ErrorApp(ctx, apicode.StoreRequired)
		return
	}
	msg, e := c.svc.Chat(ctx.Request.Context(), &req, middleware.GetUserID(ctx), sid)
	if e != nil {
		httpPkg.ErrorFrom(ctx, e)
		return
	}
	httpPkg.Success(ctx, msg)
}
