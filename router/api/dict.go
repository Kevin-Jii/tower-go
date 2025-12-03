package api

import (
	"github.com/Kevin-Jii/tower-go/middleware"
	"github.com/gin-gonic/gin"
)

// RegisterDictRoutes 注册数据字典路由
func RegisterDictRoutes(v1 *gin.RouterGroup, c *Controllers) {
	// 字典类型
	dictTypes := v1.Group("/dict-types")
	dictTypes.Use(middleware.AuthMiddleware())
	{
		dictTypes.POST("", c.Dict.CreateType)
		dictTypes.GET("", c.Dict.ListTypes)
		dictTypes.GET("/:id", c.Dict.GetType)
		dictTypes.PUT("/:id", c.Dict.UpdateType)
		dictTypes.DELETE("/:id", c.Dict.DeleteType)
	}

	// 字典数据
	dictData := v1.Group("/dict-data")
	dictData.Use(middleware.AuthMiddleware())
	{
		dictData.POST("", c.Dict.CreateData)
		dictData.GET("", c.Dict.ListDataByType)
		dictData.GET("/:id", c.Dict.GetData)
		dictData.PUT("/:id", c.Dict.UpdateData)
		dictData.DELETE("/:id", c.Dict.DeleteData)
	}

	// 获取所有字典（用于前端缓存）
	v1.GET("/dicts", middleware.AuthMiddleware(), c.Dict.GetAllDict)
}
