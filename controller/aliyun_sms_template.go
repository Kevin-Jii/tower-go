package controller

import (
	"github.com/Kevin-Jii/tower-go/middleware"
	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/service"
	httpPkg "github.com/Kevin-Jii/tower-go/utils/http"
	"github.com/gin-gonic/gin"
)

type AliyunSmsTemplateController struct {
	svc *service.AliyunSmsTemplateService
}

func NewAliyunSmsTemplateController(s *service.AliyunSmsTemplateService) *AliyunSmsTemplateController {
	return &AliyunSmsTemplateController{svc: s}
}

func (c *AliyunSmsTemplateController) List(ctx *gin.Context) {
	rows, err := c.svc.List(ctx.Query("keyword"), ctx.Query("audit_status"))
	if err != nil {
		httpPkg.ErrorFrom(ctx, err)
		return
	}
	httpPkg.Success(ctx, rows)
}

func (c *AliyunSmsTemplateController) ListApproved(ctx *gin.Context) {
	rows, err := c.svc.ListApproved()
	if err != nil {
		httpPkg.ErrorFrom(ctx, err)
		return
	}
	httpPkg.Success(ctx, rows)
}

func (c *AliyunSmsTemplateController) Create(ctx *gin.Context) {
	var req model.CreateAliyunSmsTemplateReq
	if !httpPkg.BindJSON(ctx, &req) {
		return
	}
	row, err := c.svc.Create(&req, middleware.GetUserID(ctx))
	if err != nil {
		httpPkg.ErrorFrom(ctx, err)
		return
	}
	httpPkg.Success(ctx, row)
}

func (c *AliyunSmsTemplateController) Refresh(ctx *gin.Context) {
	var req model.RefreshAliyunSmsTemplateReq
	if !httpPkg.BindJSON(ctx, &req) {
		return
	}
	row, err := c.svc.Refresh(req.TemplateCode)
	if err != nil {
		httpPkg.ErrorFrom(ctx, err)
		return
	}
	httpPkg.Success(ctx, row)
}

func (c *AliyunSmsTemplateController) Delete(ctx *gin.Context) {
	code := ctx.Param("code")
	if code == "" {
		httpPkg.Error(ctx, 400, "缺少模板 CODE")
		return
	}
	if err := c.svc.Delete(code); err != nil {
		httpPkg.ErrorFrom(ctx, err)
		return
	}
	httpPkg.Success(ctx, gin.H{"message": "deleted"})
}
