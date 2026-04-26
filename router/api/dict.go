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
		dictTypes.POST("", middleware.Permission("system:dict:type:add"), c.Dict.CreateType)
		dictTypes.GET("", middleware.Permission("system:dict:list"), c.Dict.ListTypes)
		dictTypes.GET("/:id", middleware.Permission("system:dict:list"), c.Dict.GetType)
		dictTypes.PUT("/:id", middleware.Permission("system:dict:type:edit"), c.Dict.UpdateType)
		dictTypes.DELETE("/:id", middleware.Permission("system:dict:type:delete"), c.Dict.DeleteType)
	}

	// 字典数据
	dictData := v1.Group("/dict-data")
	dictData.Use(middleware.AuthMiddleware())
	{
		dictData.POST("", middleware.Permission("system:dict:data:add"), c.Dict.CreateData)
		dictData.GET("", middleware.Permission("system:dict:list"), c.Dict.ListDataByType)
		dictData.GET("/:id", middleware.Permission("system:dict:list"), c.Dict.GetData)
		dictData.PUT("/:id", middleware.Permission("system:dict:data:edit"), c.Dict.UpdateData)
		dictData.DELETE("/:id", middleware.Permission("system:dict:data:delete"), c.Dict.DeleteData)
	}

	// 获取所有字典（用于前端缓存）
	v1.GET("/dicts", middleware.AuthMiddleware(), middleware.Permission("system:dict:list"), c.Dict.GetAllDict)
}
