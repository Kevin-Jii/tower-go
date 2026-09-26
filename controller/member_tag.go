package controller

import (
	"strconv"

	"github.com/Kevin-Jii/tower-go/middleware"
	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/service"
	httpPkg "github.com/Kevin-Jii/tower-go/utils/http"
	"github.com/gin-gonic/gin"
)

type MemberTagController struct{ svc *service.MemberTagService }

func NewMemberTagController(s *service.MemberTagService) *MemberTagController {
	return &MemberTagController{svc: s}
}

func tagScope(ctx *gin.Context) (uint, bool) {
	sid := middleware.ResolveQueryStoreID(ctx, "store_id")
	return sid, middleware.HQUnboundAdmin(ctx) && sid == 0
}

func (c *MemberTagController) List(ctx *gin.Context) {
	sid, all := tagScope(ctx)
	rows, err := c.svc.List(sid, all)
	if err != nil {
		httpPkg.ErrorFrom(ctx, err)
		return
	}
	httpPkg.Success(ctx, rows)
}
func (c *MemberTagController) Get(ctx *gin.Context) {
	id, ok := httpPkg.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	sid, all := tagScope(ctx)
	row, err := c.svc.Get(id, sid, all)
	if err != nil {
		httpPkg.ErrorFrom(ctx, err)
		return
	}
	httpPkg.Success(ctx, row)
}
func (c *MemberTagController) Create(ctx *gin.Context) {
	var req model.UpsertMemberTagReq
	if !httpPkg.BindJSON(ctx, &req) {
		return
	}
	sid := middleware.ResolveQueryStoreID(ctx, "store_id")
	if middleware.HQUnboundAdmin(ctx) && sid == 0 && req.StoreID > 0 {
		sid = req.StoreID
	}
	row, err := c.svc.Create(&req, sid)
	if err != nil {
		httpPkg.ErrorFrom(ctx, err)
		return
	}
	httpPkg.Success(ctx, row)
}
func (c *MemberTagController) Update(ctx *gin.Context) {
	id, ok := httpPkg.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	var req model.UpsertMemberTagReq
	if !httpPkg.BindJSON(ctx, &req) {
		return
	}
	sid, all := tagScope(ctx)
	row, err := c.svc.Update(id, sid, all, &req)
	if err != nil {
		httpPkg.ErrorFrom(ctx, err)
		return
	}
	httpPkg.Success(ctx, row)
}
func (c *MemberTagController) Delete(ctx *gin.Context) {
	id, ok := httpPkg.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	sid, all := tagScope(ctx)
	if err := c.svc.Delete(id, sid, all); err != nil {
		httpPkg.ErrorFrom(ctx, err)
		return
	}
	httpPkg.Success(ctx, nil)
}
func (c *MemberTagController) SearchMembers(ctx *gin.Context) {
	sid, all := tagScope(ctx)
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "20"))
	rows, err := c.svc.SearchMembers(ctx.Query("keyword"), sid, all, limit)
	if err != nil {
		httpPkg.ErrorFrom(ctx, err)
		return
	}
	httpPkg.Success(ctx, rows)
}

func (c *MemberTagController) ListTagMembers(ctx *gin.Context) {
	id, ok := httpPkg.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	sid, all := tagScope(ctx)
	rows, err := c.svc.ListTagMembers(id, sid, all)
	if err != nil {
		httpPkg.ErrorFrom(ctx, err)
		return
	}
	httpPkg.Success(ctx, rows)
}
func (c *MemberTagController) BindMember(ctx *gin.Context) {
	id, ok := httpPkg.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	var req model.BindMemberTagReq
	if !httpPkg.BindJSON(ctx, &req) {
		return
	}
	sid, all := tagScope(ctx)
	row, err := c.svc.BindMember(id, req.MemberID, sid, all)
	if err != nil {
		httpPkg.ErrorFrom(ctx, err)
		return
	}
	httpPkg.Success(ctx, row)
}
func (c *MemberTagController) UnbindMember(ctx *gin.Context) {
	id, ok := httpPkg.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	memberID, ok := httpPkg.ParseUintParam(ctx, "member_id")
	if !ok {
		return
	}
	sid, all := tagScope(ctx)
	if err := c.svc.UnbindMember(id, memberID, sid, all); err != nil {
		httpPkg.ErrorFrom(ctx, err)
		return
	}
	httpPkg.Success(ctx, nil)
}
func (c *MemberTagController) ListMemberTags(ctx *gin.Context) {
	id, ok := httpPkg.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	sid, all := tagScope(ctx)
	rows, err := c.svc.ListMemberTags(id, sid, all)
	if err != nil {
		httpPkg.ErrorFrom(ctx, err)
		return
	}
	httpPkg.Success(ctx, rows)
}
func (c *MemberTagController) AssignMemberTags(ctx *gin.Context) {
	id, ok := httpPkg.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	var req model.AssignMemberTagsReq
	if !httpPkg.BindJSON(ctx, &req) {
		return
	}
	sid, all := tagScope(ctx)
	if err := c.svc.AssignMemberTags(id, sid, all, req.TagIDs); err != nil {
		httpPkg.ErrorFrom(ctx, err)
		return
	}
	httpPkg.Success(ctx, gin.H{"message": "updated"})
}
