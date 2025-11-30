package api

import (
	"github.com/Kevin-Jii/tower-go/middleware"
	"github.com/gin-gonic/gin"
)

// RegisterStoreRoutes 注册门店管理路由
func RegisterStoreRoutes(v1 *gin.RouterGroup, c *Controllers) {
	stores := v1.Group("/stores")
	stores.Use(middleware.AuthMiddleware())
	{
		stores.POST("", c.Store.CreateStore)
		stores.GET("", c.Store.ListStores)
		stores.GET("/all", c.Store.ListAllStores)
		// 门店分类相关
		stores.POST("/:id/dish-categories", c.DishCategory.CreateCategoryForStore)
		stores.GET("/:id/dish-categories", c.DishCategory.ListCategoriesForStore)
		stores.POST("/:id/dish-categories/:cid/dishes", c.DishCategory.CreateDishForStoreCategory)
		stores.GET("/:id/dish-categories/:cid/dishes", c.DishCategory.ListDishesForStoreCategory)
		stores.PUT("/:id/dish-categories/:cid/dishes/:did", c.DishCategory.UpdateDishForStoreCategory)
		stores.DELETE("/:id/dish-categories/:cid/dishes/:did", c.DishCategory.DeleteDishForStoreCategory)
		stores.PUT("/:id/dish-categories/:cid", c.DishCategory.UpdateCategoryForStore)
		stores.DELETE("/:id/dish-categories/:cid", c.DishCategory.DeleteCategoryForStore)
		stores.PUT("/:id/dishes/:did", c.Dish.UpdateDishForStore)
		stores.DELETE("/:id/dishes/:did", c.Dish.DeleteDishForStore)
		stores.GET("/:id", c.Store.GetStore)
		stores.PUT("/:id", c.Store.UpdateStore)
		stores.DELETE("/:id", c.Store.DeleteStore)
	}
}
