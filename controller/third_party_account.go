package controller

import (
	"strconv"

	"github.com/Kevin-Jii/tower-go/middleware"
	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/service"
	httpPkg "github.com/Kevin-Jii/tower-go/utils/http"
	"github.com/gin-gonic/gin"
)

type ThirdPartyAccountController struct {
	svc *service.ThirdPartyAccountService
}

func NewThirdPartyAccountController(svc *service.ThirdPartyAccountService) *ThirdPartyAccountController {
	return &ThirdPartyAccountController{svc: svc}
}

func (c *ThirdPartyAccountController) List(ctx *gin.Context) {
	if !middleware.IsAdmin(ctx) {
		httpPkg.Error(ctx, 403, "only admin can view")
		return
	}
	rows, err := c.svc.List(ctx.Query("keyword"))
	if err != nil {
		httpPkg.Error(ctx, 500, err.Error())
		return
	}
	httpPkg.Success(ctx, rows)
}

func (c *ThirdPartyAccountController) Get(ctx *gin.Context) {
	if !middleware.IsAdmin(ctx) {
		httpPkg.Error(ctx, 403, "only admin can view")
		return
	}
	id, ok := httpPkg.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	row, err := c.svc.GetByID(id)
	if err != nil {
		httpPkg.Error(ctx, 404, "account not found")
		return
	}
	httpPkg.Success(ctx, row)
}

func (c *ThirdPartyAccountController) Create(ctx *gin.Context) {
	if !middleware.IsAdmin(ctx) {
		httpPkg.Error(ctx, 403, "only admin can create")
		return
	}
	var req model.CreateThirdPartyAccountReq
	if !httpPkg.BindJSON(ctx, &req) {
		return
	}
	row, err := c.svc.Create(&req)
	if err != nil {
		httpPkg.Error(ctx, 500, err.Error())
		return
	}
	httpPkg.Success(ctx, row)
}

func (c *ThirdPartyAccountController) Update(ctx *gin.Context) {
	if !middleware.IsAdmin(ctx) {
		httpPkg.Error(ctx, 403, "only admin can update")
		return
	}
	id, ok := httpPkg.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	var req model.UpdateThirdPartyAccountReq
	if !httpPkg.BindJSON(ctx, &req) {
		return
	}
	if err := c.svc.Update(id, &req); err != nil {
		httpPkg.Error(ctx, 500, err.Error())
		return
	}
	httpPkg.Success(ctx, gin.H{"message": "updated"})
}

func (c *ThirdPartyAccountController) Delete(ctx *gin.Context) {
	if !middleware.IsAdmin(ctx) {
		httpPkg.Error(ctx, 403, "only admin can delete")
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

func (c *ThirdPartyAccountController) TestLogin(ctx *gin.Context) {
	if !middleware.IsAdmin(ctx) {
		httpPkg.Error(ctx, 403, "only admin can test")
		return
	}
	id, ok := httpPkg.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	res, err := c.svc.TestLogin(id)
	if err != nil {
		httpPkg.Error(ctx, 500, err.Error())
		return
	}
	httpPkg.Success(ctx, res)
}

func (c *ThirdPartyAccountController) SyncLatestOrders(ctx *gin.Context) {
	if !middleware.IsAdmin(ctx) {
		httpPkg.Error(ctx, 403, "only admin can sync")
		return
	}
	id, ok := httpPkg.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	res, err := c.svc.SyncLatestOrders(id)
	if err != nil {
		httpPkg.Error(ctx, 500, err.Error())
		return
	}
	httpPkg.Success(ctx, res)
}

func (c *ThirdPartyAccountController) ListSyncedOrders(ctx *gin.Context) {
	if !middleware.IsAdmin(ctx) {
		httpPkg.Error(ctx, 403, "only admin can view")
		return
	}
	id, ok := httpPkg.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))
	rows, total, err := c.svc.ListSyncedOrders(id, page, pageSize)
	if err != nil {
		httpPkg.Error(ctx, 500, err.Error())
		return
	}
	httpPkg.SuccessWithPagination(ctx, rows, total, page, pageSize)
}
