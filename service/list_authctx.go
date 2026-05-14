package service

import (
	"context"

	"github.com/Kevin-Jii/tower-go/internal/authctx"
	"github.com/Kevin-Jii/tower-go/model"
)

// applyListRBACFromContext 将 HTTP 链路注入的 AuthContext 写入列表请求（P2：与 Controller 解耦）。
// 无 AuthContext 时（如钉钉机器人直连 module）不修改 req，沿用调用方已设字段。
func applyListRBACFromContextToPurchaseOrder(ctx context.Context, req *model.ListPurchaseOrderReq) {
	if req == nil {
		return
	}
	if ac := authctx.FromContext(ctx); ac != nil {
		req.DataScope = ac.EffectiveDataScope
		req.UserID = ac.UserID
		req.RoleCode = ac.RoleCode
	}
}

func applyListRBACFromContextToStoreAccount(ctx context.Context, req *model.ListStoreAccountReq) {
	if req == nil {
		return
	}
	if ac := authctx.FromContext(ctx); ac != nil {
		req.DataScope = ac.EffectiveDataScope
		req.UserID = ac.UserID
		req.RoleCode = ac.RoleCode
	}
}

func applyListRBACFromContextToInventory(ctx context.Context, req *model.ListInventoryReq) {
	if req == nil {
		return
	}
	if ac := authctx.FromContext(ctx); ac != nil {
		req.DataScope = ac.EffectiveDataScope
		req.UserID = ac.UserID
		req.RoleCode = ac.RoleCode
	}
}

func applyListRBACFromContextToInventoryOrder(ctx context.Context, req *model.ListInventoryOrderReq) {
	if req == nil {
		return
	}
	if ac := authctx.FromContext(ctx); ac != nil {
		req.DataScope = ac.EffectiveDataScope
		req.UserID = ac.UserID
		req.RoleCode = ac.RoleCode
	}
}
