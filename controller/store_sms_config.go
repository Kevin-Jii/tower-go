package controller

import (
	"github.com/Kevin-Jii/tower-go/middleware"
	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/service"
	httpPkg "github.com/Kevin-Jii/tower-go/utils/http"
	"github.com/gin-gonic/gin"
)

type StoreSmsConfigController struct {
	svc *service.StoreSmsConfigService
}

func NewStoreSmsConfigController(s *service.StoreSmsConfigService) *StoreSmsConfigController {
	return &StoreSmsConfigController{svc: s}
}

// resolveStoreID 返回目标门店 ID：优先路径参数，再 query store_id，再 token storeID。
// 总部未绑定 admin 可显式选 store_id；门店用户强制使用 token 中的 store_id。
func (c *StoreSmsConfigController) resolveStoreID(ctx *gin.Context) (uint, bool) {
	storeID := middleware.ResolveQueryStoreID(ctx, "store_id")
	if middleware.HQUnboundAdmin(ctx) {
		return storeID, true
	}
	tid := middleware.GetStoreID(ctx)
	if storeID != 0 && storeID != tid {
		httpPkg.Error(ctx, 403, "门店用户只能配置所属门店")
		return 0, false
	}
	return tid, true
}

func (c *StoreSmsConfigController) Get(ctx *gin.Context) {
	storeID, ok := c.resolveStoreID(ctx)
	if !ok {
		return
	}
	if storeID == 0 {
		httpPkg.Error(ctx, 400, "请指定门店 ID（?store_id=）")
		return
	}
	row, err := c.svc.GetByStoreID(storeID)
	if err != nil {
		httpPkg.ErrorFrom(ctx, err)
		return
	}
	httpPkg.Success(ctx, row)
}

func (c *StoreSmsConfigController) Upsert(ctx *gin.Context) {
	storeID, ok := c.resolveStoreID(ctx)
	if !ok {
		return
	}
	if storeID == 0 {
		httpPkg.Error(ctx, 400, "请指定门店 ID（?store_id=）")
		return
	}
	var req model.UpsertStoreSmsConfigReq
	if !httpPkg.BindJSON(ctx, &req) {
		return
	}
	row, err := c.svc.Upsert(storeID, &req, middleware.GetUserID(ctx))
	if err != nil {
		httpPkg.ErrorFrom(ctx, err)
		return
	}
	httpPkg.Success(ctx, row)
}

func (c *StoreSmsConfigController) Test(ctx *gin.Context) {
	var req model.StoreSmsConfigTestReq
	if !httpPkg.BindJSON(ctx, &req) {
		return
	}
	msg, err := c.svc.TestConnection(&req)
	if err != nil {
		httpPkg.ErrorFrom(ctx, err)
		return
	}
	httpPkg.Success(ctx, gin.H{"message": msg})
}
