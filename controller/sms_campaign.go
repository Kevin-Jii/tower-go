package controller

import (
	"github.com/Kevin-Jii/tower-go/middleware"
	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/pkg/apicode"
	"github.com/Kevin-Jii/tower-go/service"
	httpPkg "github.com/Kevin-Jii/tower-go/utils/http"
	"github.com/gin-gonic/gin"
)

type SmsCampaignController struct{ svc *service.SmsCampaignService }

func NewSmsCampaignController(svc *service.SmsCampaignService) *SmsCampaignController {
	return &SmsCampaignController{svc: svc}
}
func campaignScope(ctx *gin.Context) (uint, bool) {
	sid := middleware.ResolveQueryStoreID(ctx, "store_id")
	return sid, middleware.HQUnboundAdmin(ctx) && sid == 0
}
func (c *SmsCampaignController) Config(ctx *gin.Context) {
	sid := middleware.ResolveQueryStoreID(ctx, "store_id")
	fallback := c.svc.GetStoreSignName(sid)
	httpPkg.Success(ctx, c.svc.GetConfig(sid, fallback))
}
func (c *SmsCampaignController) List(ctx *gin.Context) {
	sid, all := campaignScope(ctx)
	rows, err := c.svc.List(sid, all)
	if err != nil {
		httpPkg.ErrorFrom(ctx, err)
		return
	}
	httpPkg.Success(ctx, rows)
}
func (c *SmsCampaignController) Get(ctx *gin.Context) {
	id, ok := httpPkg.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	sid, all := campaignScope(ctx)
	row, err := c.svc.GetByID(id, sid, all)
	if err != nil {
		httpPkg.Error(ctx, 404, "活动不存在")
		return
	}
	httpPkg.Success(ctx, row)
}
func (c *SmsCampaignController) Create(ctx *gin.Context) {
	var req model.CreateSmsCampaignReq
	if !httpPkg.BindJSON(ctx, &req) {
		return
	}
	if req.ScheduledAt != nil && !requireSmsSendPermission(ctx) {
		return
	}
	effective := middleware.ResolveQueryStoreID(ctx, "store_id")
	if effective == 0 && req.OwnerStoreID > 0 {
		effective = req.OwnerStoreID
	}
	row, err := c.svc.Create(&req, middleware.GetUserID(ctx), effective, middleware.HQUnboundAdmin(ctx))
	if err != nil {
		httpPkg.ErrorFrom(ctx, err)
		return
	}
	httpPkg.Success(ctx, row)
}
func (c *SmsCampaignController) Update(ctx *gin.Context) {
	id, ok := httpPkg.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	var req model.UpdateSmsCampaignReq
	if !httpPkg.BindJSON(ctx, &req) {
		return
	}
	if req.ScheduledAt != nil && !requireSmsSendPermission(ctx) {
		return
	}
	sid, all := campaignScope(ctx)
	if err := c.svc.Update(id, &req, sid, all); err != nil {
		httpPkg.ErrorFrom(ctx, err)
		return
	}
	httpPkg.Success(ctx, gin.H{"message": "updated"})
}
func (c *SmsCampaignController) Delete(ctx *gin.Context) {
	id, ok := httpPkg.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	sid, all := campaignScope(ctx)
	if err := c.svc.Delete(id, sid, all); err != nil {
		httpPkg.ErrorFrom(ctx, err)
		return
	}
	httpPkg.Success(ctx, gin.H{"message": "deleted"})
}
func (c *SmsCampaignController) Send(ctx *gin.Context) {
	id, ok := httpPkg.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	sid, all := campaignScope(ctx)
	if err := c.svc.SendNow(id, sid, all); err != nil {
		httpPkg.ErrorFrom(ctx, err)
		return
	}
	httpPkg.Success(ctx, gin.H{"message": "send completed"})
}
func (c *SmsCampaignController) Cancel(ctx *gin.Context) {
	id, ok := httpPkg.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	sid, all := campaignScope(ctx)
	if err := c.svc.Cancel(id, sid, all); err != nil {
		httpPkg.ErrorFrom(ctx, err)
		return
	}
	httpPkg.Success(ctx, gin.H{"message": "cancelled"})
}
func requireSmsSendPermission(ctx *gin.Context) bool {
	allowed, err := middleware.HasPermission(ctx, "marketing:sms:send")
	if err != nil {
		httpPkg.ErrorApp(ctx, apicode.PermissionLoadFailed)
		return false
	}
	if !allowed {
		httpPkg.ErrorApp(ctx, apicode.PermissionDenied)
		return false
	}
	return true
}

func (c *SmsCampaignController) Records(ctx *gin.Context) {
	id, ok := httpPkg.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	sid, all := campaignScope(ctx)
	rows, err := c.svc.ListSendRecords(id, sid, all)
	if err != nil {
		httpPkg.ErrorFrom(ctx, err)
		return
	}
	httpPkg.Success(ctx, rows)
}

func (c *SmsCampaignController) RetryRecord(ctx *gin.Context) {
	campaignID, ok := httpPkg.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	recordID, ok := httpPkg.ParseUintParam(ctx, "record_id")
	if !ok {
		return
	}
	sid, all := campaignScope(ctx)
	if err := c.svc.RetryFailedRecord(campaignID, recordID, sid, all); err != nil {
		// Retry errors are actionable SMS/provider messages and must be visible to
		// the operator as well as persisted on the failed send record.
		httpPkg.Error(ctx, apicode.InvalidParameter.Num, err.Error())
		return
	}
	httpPkg.Success(ctx, gin.H{"message": "resent"})
}
