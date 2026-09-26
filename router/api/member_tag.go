package api

import (
	"github.com/Kevin-Jii/tower-go/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterMemberTagRoutes(v1 *gin.RouterGroup, c *Controllers) {
	tags := v1.Group("/member-tags")
	tags.Use(middleware.AuthMiddleware(), middleware.StoreBusinessGuard())
	{
		tags.GET("", middleware.Permission("marketing:sms:list"), c.MemberTag.List)
		tags.GET("/members/search", middleware.Permission("marketing:sms:list"), c.MemberTag.SearchMembers)
		tags.POST("", middleware.Permission("marketing:sms:add"), c.MemberTag.Create)
		tags.GET("/:id", middleware.Permission("marketing:sms:list"), c.MemberTag.Get)
		tags.PUT("/:id", middleware.Permission("marketing:sms:edit"), c.MemberTag.Update)
		tags.DELETE("/:id", middleware.Permission("marketing:sms:delete"), c.MemberTag.Delete)
		tags.GET("/:id/members", middleware.Permission("marketing:sms:list"), c.MemberTag.ListTagMembers)
		tags.POST("/:id/members", middleware.Permission("marketing:sms:edit"), c.MemberTag.BindMember)
		tags.DELETE("/:id/members/:member_id", middleware.Permission("marketing:sms:edit"), c.MemberTag.UnbindMember)
	}
	memberTags := v1.Group("/members/:id/tags")
	memberTags.Use(middleware.AuthMiddleware(), middleware.StoreBusinessGuard())
	{
		memberTags.GET("", middleware.Permission("marketing:sms:list"), c.MemberTag.ListMemberTags)
		memberTags.PUT("", middleware.Permission("marketing:sms:edit"), c.MemberTag.AssignMemberTags)
	}
}
