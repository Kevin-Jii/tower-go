package controller

import (
	"strings"

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

func templateScope(ctx *gin.Context) (storeID uint, allStores bool) {
	storeID = middleware.ResolveQueryStoreID(ctx, "store_id")
	return storeID, middleware.HQUnboundAdmin(ctx) && storeID == 0
}

// Template codes are unique only inside one Aliyun account. HQ mutations must
// therefore carry the row owner explicitly, including owner_store_id=0.
func templateMutationScope(ctx *gin.Context) (storeID uint, ok bool) {
	if !middleware.HQUnboundAdmin(ctx) {
		return middleware.GetStoreID(ctx), true
	}
	if _, exists := ctx.GetQuery("owner_store_id"); !exists {
		httpPkg.Error(ctx, 400, "总部操作模板时必须指定 owner_store_id")
		return 0, false
	}
	return middleware.ResolveQueryStoreID(ctx, "owner_store_id"), true
}

func (c *AliyunSmsTemplateController) List(ctx *gin.Context) {
	storeID, allStores := templateScope(ctx)
	if middleware.HQUnboundAdmin(ctx) {
		if _, exists := ctx.GetQuery("owner_store_id"); exists {
			storeID = middleware.ResolveQueryStoreID(ctx, "owner_store_id")
			allStores = false
		}
	}
	rows, err := c.svc.List(storeID, allStores, ctx.Query("keyword"), ctx.Query("audit_status"))
	if err != nil {
		httpPkg.ErrorFrom(ctx, err)
		return
	}
	httpPkg.Success(ctx, rows)
}

func (c *AliyunSmsTemplateController) ListApproved(ctx *gin.Context) {
	storeID, allStores := templateScope(ctx)
	rows, err := c.svc.ListApproved(storeID, allStores)
	if err != nil {
		httpPkg.ErrorFrom(ctx, err)
		return
	}
	httpPkg.Success(ctx, rows)
}

func (c *AliyunSmsTemplateController) ListSignatures(ctx *gin.Context) {
	storeID, _ := templateScope(ctx)
	approvedOnly := strings.TrimSpace(ctx.Query("approved_only")) != "false"
	rows, err := c.svc.ListSignatures(storeID, approvedOnly)
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
	ownerStoreID := middleware.ResolveQueryStoreID(ctx, "store_id")
	if middleware.HQUnboundAdmin(ctx) && req.OwnerStoreID > 0 {
		ownerStoreID = req.OwnerStoreID
	}
	row, err := c.svc.Create(&req, ownerStoreID, middleware.GetUserID(ctx))
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
	storeID, ok := templateMutationScope(ctx)
	if !ok {
		return
	}
	row, err := c.svc.Refresh(req.TemplateCode, storeID, false)
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
	storeID, ok := templateMutationScope(ctx)
	if !ok {
		return
	}
	if err := c.svc.Delete(code, storeID, false); err != nil {
		httpPkg.ErrorFrom(ctx, err)
		return
	}
	httpPkg.Success(ctx, gin.H{"message": "deleted"})
}
