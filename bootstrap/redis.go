package bootstrap

import (
	"tower-go/config"
	"tower-go/utils/redis"
	"tower-go/utils/logging"

	"go.uber.org/zap"
)

func InitRedisCache() func() {
	redisCfg := config.GetRedisConfig()
	if err := redis.InitRedis(redisCfg); err != nil {
		logging.LogWarn("Redis 连接失败，缓存禁用", zap.Error(err))
		return func() {}
	}
	if redis.IsRedisEnabled() {
		logging.LogInfo("Redis 缓存已启用")
	}
	return func() { _ = redis.CloseRedis() }
}
