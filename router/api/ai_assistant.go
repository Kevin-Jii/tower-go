package api

import (
	"github.com/Kevin-Jii/tower-go/middleware"
	"github.com/Kevin-Jii/tower-go/model"
	"github.com/gin-gonic/gin"
)

func RegisterAIAssistantRoutes(v1 *gin.RouterGroup, c *Controllers) {
	g := v1.Group("/ai-assistant", middleware.AuthMiddleware())
	g.GET("/config", c.AIAssistant.GetConfig)
	g.PUT("/config", c.AIAssistant.SaveConfig)
	g.POST("/config/test", c.AIAssistant.TestConfig)
	useAssistant := middleware.PermissionOrRoles("ai:assistant:use", model.RoleCodeStoreAdmin)
	g.GET("/conversations", useAssistant, c.AIAssistant.Conversations)
	g.POST("/conversations", useAssistant, c.AIAssistant.CreateConversation)
	g.GET("/conversations/:id/messages", useAssistant, c.AIAssistant.Messages)
	g.PUT("/conversations/:id", useAssistant, c.AIAssistant.Rename)
	g.DELETE("/conversations/:id", useAssistant, c.AIAssistant.Delete)
	g.POST("/chat", useAssistant, c.AIAssistant.Chat)
}
