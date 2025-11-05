package bootstrap

import (
	"tower-go/config"
	"tower-go/utils"

	"go.uber.org/zap"
)

func InitRedisCache() func() {
	redisCfg := config.GetRedisConfig()
	if err := utils.InitRedis(redisCfg); err != nil {
		utils.LogWarn("Redis 连接失败，缓存禁用", zap.Error(err))
		return func() {}
	}
	if utils.IsRedisEnabled() {
		utils.LogInfo("Redis 缓存已启用")
	}
	return func() { _ = utils.CloseRedis() }
}
