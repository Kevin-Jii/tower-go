package api

import (
	"github.com/Kevin-Jii/tower-go/controller"
	"github.com/Kevin-Jii/tower-go/middleware"
	"github.com/gin-gonic/gin"
)

// CRUDController 通用 CRUD 适配接口
type CRUDController interface {
	Create(*gin.Context)
	List(*gin.Context)
	Get(*gin.Context)
	Update(*gin.Context)
	Delete(*gin.Context)
}

type dishCategoryCRUDAdapter struct {
	inner *controller.DishCategoryController
}

func (a dishCategoryCRUDAdapter) Create(ctx *gin.Context) { a.inner.CreateCategory(ctx) }
func (a dishCategoryCRUDAdapter) List(ctx *gin.Context)   { a.inner.ListCategories(ctx) }
func (a dishCategoryCRUDAdapter) Get(ctx *gin.Context)    { a.inner.ListCategoriesForStore(ctx) }
func (a dishCategoryCRUDAdapter) Update(ctx *gin.Context) { a.inner.UpdateCategory(ctx) }
func (a dishCategoryCRUDAdapter) Delete(ctx *gin.Context) { a.inner.DeleteCategory(ctx) }

func registerCRUD(group *gin.RouterGroup, ctrl CRUDController) {
	group.POST("", ctrl.Create)
	group.GET("", ctrl.List)
	group.PUT(":id", ctrl.Update)
	group.DELETE(":id", ctrl.Delete)
}

// RegisterDishRoutes 注册菜品管理路由
func RegisterDishRoutes(v1 *gin.RouterGroup, c *Controllers) {
	dishes := v1.Group("/dishes")
	dishes.Use(middleware.StoreAuthMiddleware())
	{
		dishes.POST("", c.Dish.CreateDish)
		dishes.GET("", c.Dish.ListDishes)
		dishes.GET("/:id", c.Dish.GetDish)
		dishes.PUT("/:id", c.Dish.UpdateDish)
		dishes.DELETE("/:id", c.Dish.DeleteDish)
		dishes.GET("/by-category", c.DishCategory.ListCategoriesWithDishes)
	}

	cats := v1.Group("/dish-categories")
	cats.Use(middleware.StoreAuthMiddleware())
	{
		registerCRUD(cats, dishCategoryCRUDAdapter{inner: c.DishCategory})
		cats.POST("/reorder", c.DishCategory.ReorderCategories)
	}
}
