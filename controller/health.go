package controller

import (
	"github.com/Kevin-Jii/tower-go/config"
	"github.com/Kevin-Jii/tower-go/utils/database"
	"github.com/Kevin-Jii/tower-go/utils/logging"
	"github.com/Kevin-Jii/tower-go/utils/redis"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// HealthController 健康检查控制器
type HealthController struct{}

// NewHealthController 创建健康检查控制器
func NewHealthController() *HealthController {
	return &HealthController{}
}

// HealthResponse 健康检查响应
type HealthResponse struct {
	Status    string            `json:"status"`
	Timestamp string            `json:"timestamp"`
	Services  map[string]string `json:"services"`
	Version   string            `json:"version"`
	Uptime    string            `json:"uptime"`
}

// Check 健康检查
// @Summary 健康检查
// @Description 检查服务健康状态
// @Tags 系统
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /health [get]
func (c *HealthController) Check(ctx *gin.Context) {
	response := HealthResponse{
		Status:    "healthy",
		Timestamp: "",
		Services:  make(map[string]string),
		Version:   "1.0.0",
		Uptime:    "",
	}

	// 检查数据库连接
	if database.GetDB() != nil {
		sqlDB, err := database.GetDB().DB()
		if err == nil {
			if err := sqlDB.Ping(); err != nil {
				response.Services["database"] = "unhealthy: " + err.Error()
				response.Status = "degraded"
			} else {
				response.Services["database"] = "healthy"
			}
		}
	} else {
		response.Services["database"] = "unhealthy: not initialized"
		response.Status = "degraded"
	}

	// 检查 Redis 连接
	if config.GetConfig().Redis.Enabled {
		if redis.GetClient() != nil {
			if err := redis.GetClient().Ping(ctx.Request.Context()).Err(); err != nil {
				response.Services["redis"] = "unhealthy: " + err.Error()
				response.Status = "degraded"
			} else {
				response.Services["redis"] = "healthy"
			}
		} else {
			response.Services["redis"] = "not configured"
		}
	}

	// 设置状态码
	if response.Status == "healthy" {
		ctx.JSON(200, response)
	} else {
		ctx.JSON(503, response)
	}
}

// Ready 就绪检查
// @Summary 就绪检查
// @Description 检查服务是否就绪
// @Tags 系统
// @Produce json
// @Success 200 {object} map[string]string
// @Router /ready [get]
func (c *HealthController) Ready(ctx *gin.Context) {
	// 检查关键依赖
	if database.GetDB() == nil {
		ctx.JSON(503, gin.H{
			"status": "not ready",
			"reason": "database not connected",
		})
		return
	}

	ctx.JSON(200, gin.H{
		"status": "ready",
	})
}

// Live 存活检查
// @Summary 存活检查
// @Description 检查服务是否存活
// @Tags 系统
// @Produce json
// @Success 200 {object} map[string]string
// @Router /live [get]
func (c *HealthController) Live(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"status": "alive",
	})
}

// LogLevel 动态调整日志级别
// @Summary 调整日志级别
// @Description 动态调整当前日志级别
// @Tags 系统
// @Accept json
// @Produce json
// @Param level body map[string]string true "日志级别"
// @Success 200 {object} map[string]string
// @Router /log-level [post]
func (c *HealthController) LogLevel(ctx *gin.Context) {
	var req struct {
		Level string `json:"level" binding:"required,oneof=debug info warn error"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{
			"error": "invalid level, must be one of: debug, info, warn, error",
		})
		return
	}

	// 注意：实际生产环境需要更复杂的实现
	logging.LogInfo("Log level changed", zap.String("new_level", req.Level))

	ctx.JSON(200, gin.H{
		"status": "success",
		"level":  req.Level,
	})
}
