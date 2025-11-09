package controller

import (
	"strconv"
	"tower-go/middleware"
	"tower-go/model"
	"tower-go/service"
	"tower-go/utils/http"

	"github.com/gin-gonic/gin"
)

type DishController struct {
	dishService *service.DishService
}

func NewDishController(dishService *service.DishService) *DishController {
	return &DishController{dishService: dishService}
}

// CreateDish godoc
// @Summary 创建菜品
// @Description 创建新菜品（自动关联当前门店）
// @Tags dishes
// @Accept json
// @Produce json
// @Security Bearer
// @Param dish body model.CreateDishReq true "菜品信息"
// @Success 200 {object} utils.Response
// @Router /dishes [post]
func (c *DishController) CreateDish(ctx *gin.Context) {
	storeID := middleware.GetStoreID(ctx)

	var req model.CreateDishReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		http.Error(ctx, 400, err.Error())
		return
	}

	if err := c.dishService.CreateDish(storeID, &req); err != nil {
		status := 500
		if err.Error() == "dish name already exists in this category" {
			status = 409
		}
		http.Error(ctx, status, err.Error())
		return
	}

	http.Success(ctx, nil)
}

// GetDish godoc
// @Summary 获取菜品详情
// @Description 获取菜品详细信息
// @Tags dishes
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "菜品ID"
// @Success 200 {object} utils.Response{data=model.Dish}
// @Router /dishes/{id} [get]
func (c *DishController) GetDish(ctx *gin.Context) {
	storeID := middleware.GetStoreID(ctx)

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		http.Error(ctx, 400, "Invalid dish ID")
		return
	}

	dish, err := c.dishService.GetDish(uint(id), storeID)
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	http.Success(ctx, dish)
}

// ListDishes godoc
// @Summary 菜品列表
// @Description 获取当前门店的菜品列表，支持分页
// @Tags dishes
// @Accept json
// @Produce json
// @Security Bearer
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Param category_id query int false "分类ID"
// @Success 200 {object} utils.Response{data=[]model.Dish} "分页 meta: total,page,page_size,page_count,has_more"
// @Router /dishes [get]
func (c *DishController) ListDishes(ctx *gin.Context) {
	storeID := middleware.GetStoreID(ctx)
	categoryIDStr := ctx.Query("category_id")
	if categoryIDStr != "" {
		cid, err := strconv.ParseUint(categoryIDStr, 10, 32)
		if err != nil {
			http.Error(ctx, 400, "invalid category_id")
			return
		}
		dishes, err := c.dishService.ListDishesByCategory(storeID, uint(cid))
		if err != nil {
			http.Error(ctx, 500, err.Error())
			return
		}
		http.Success(ctx, dishes)
		return
	}

	// 否则返回分页列表
	page := http.GetPage(ctx)
	pageSize := http.GetPageSize(ctx)

	dishes, total, err := c.dishService.ListDishes(storeID, page, pageSize)
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	http.SuccessWithPagination(ctx, dishes, total, page, pageSize)
}

// UpdateDish godoc
// @Summary 更新菜品信息
// @Description 更新菜品信息（含上下架）
// @Tags dishes
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "菜品ID"
// @Param dish body model.UpdateDishReq true "菜品信息"
// @Success 200 {object} utils.Response
// @Router /dishes/{id} [put]
func (c *DishController) UpdateDish(ctx *gin.Context) {
	storeID := middleware.GetStoreID(ctx)

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		http.Error(ctx, 400, "Invalid dish ID")
		return
	}

	var req model.UpdateDishReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		http.Error(ctx, 400, err.Error())
		return
	}

	if err := c.dishService.UpdateDish(uint(id), storeID, &req); err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	http.Success(ctx, nil)
}

// DeleteDish godoc
// @Summary 删除菜品
// @Description 删除菜品
// @Tags dishes
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "菜品ID"
// @Success 200 {object} utils.Response
// @Router /dishes/{id} [delete]
func (c *DishController) DeleteDish(ctx *gin.Context) {
	storeID := middleware.GetStoreID(ctx)

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		http.Error(ctx, 400, "Invalid dish ID")
		return
	}

	if err := c.dishService.DeleteDish(uint(id), storeID); err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	http.Success(ctx, nil)
}

// DeleteDishForStore godoc
// @Summary 删除指定门店的菜品
// @Tags dishes
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "门店ID"
// @Param did path int true "菜品ID"
// @Success 200 {object} utils.Response
// @Router /stores/{id}/dishes/{did} [delete]
func (c *DishController) DeleteDishForStore(ctx *gin.Context) {
	storeID, ok := http.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	dishID, ok := http.ParseUintParam(ctx, "did")
	if !ok {
		return
	}
	currentStoreID := middleware.GetStoreID(ctx)
	if currentStoreID != storeID && !middleware.IsAdmin(ctx) {
		http.Error(ctx, 403, "forbidden: cross-store delete")
		return
	}
	if err := c.dishService.DeleteDish(dishID, storeID); err != nil {
		status := 500
		if err.Error() == "dish not found" {
			status = 404
		}
		http.Error(ctx, status, err.Error())
		return
	}
	http.Success(ctx, nil)
}

// UpdateDishForStore godoc
// @Summary 更新指定门店的菜品
// @Tags dishes
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "门店ID"
// @Param did path int true "菜品ID"
// @Param dish body model.UpdateDishReq true "菜品信息"
// @Success 200 {object} utils.Response
// @Router /stores/{id}/dishes/{did} [put]
func (c *DishController) UpdateDishForStore(ctx *gin.Context) {
	storeID, ok := http.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	dishID, ok := http.ParseUintParam(ctx, "did")
	if !ok {
		return
	}
	currentStoreID := middleware.GetStoreID(ctx)
	if currentStoreID != storeID && !middleware.IsAdmin(ctx) {
		http.Error(ctx, 403, "forbidden: cross-store update")
		return
	}
	var req model.UpdateDishReq
	if !http.BindJSON(ctx, &req) {
		return
	}
	if err := c.dishService.UpdateDish(dishID, storeID, &req); err != nil {
		status := 500
		if err.Error() == "dish not found" {
			status = 404
		} else if err.Error() == "dish name already exists in this category" {
			status = 409
		}
		http.Error(ctx, status, err.Error())
		return
	}
	http.Success(ctx, nil)
}
