package controller

import (
	"github.com/Kevin-Jii/tower-go/middleware"
	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/service"
	"github.com/gin-gonic/gin"

	httpPkg "github.com/Kevin-Jii/tower-go/utils/http"
)

type SmsCampaignController struct {
	svc *service.SmsCampaignService
}

func NewSmsCampaignController(svc *service.SmsCampaignService) *SmsCampaignController {
	return &SmsCampaignController{svc: svc}
}

func (c *SmsCampaignController) Config(ctx *gin.Context) {
	httpPkg.Success(ctx, c.svc.GetConfig())
}

func (c *SmsCampaignController) List(ctx *gin.Context) {
	rows, err := c.svc.List()
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
	row, err := c.svc.GetByID(id)
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
	userID := middleware.GetUserID(ctx)
	row, err := c.svc.Create(&req, userID)
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
	if err := c.svc.Update(id, &req); err != nil {
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
	if err := c.svc.Delete(id); err != nil {
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
	if err := c.svc.SendNow(id, middleware.IsAdmin(ctx), middleware.GetStoreID(ctx)); err != nil {
		httpPkg.ErrorFrom(ctx, err)
		return
	}
	httpPkg.Success(ctx, gin.H{"message": "send started"})
}

func (c *SmsCampaignController) Cancel(ctx *gin.Context) {
	id, ok := httpPkg.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	if err := c.svc.Cancel(id); err != nil {
		httpPkg.ErrorFrom(ctx, err)
		return
	}
	httpPkg.Success(ctx, gin.H{"message": "cancelled"})
}

func (c *SmsCampaignController) Records(ctx *gin.Context) {
	id, ok := httpPkg.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	rows, err := c.svc.ListSendRecords(id)
	if err != nil {
		httpPkg.ErrorFrom(ctx, err)
		return
	}
	httpPkg.Success(ctx, rows)
}
