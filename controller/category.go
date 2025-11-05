package controller

import (
	"net/http"
	"strconv"
	"tower-go/middleware"
	"tower-go/model"
	"tower-go/service"
	"tower-go/utils"

	"github.com/gin-gonic/gin"
)

type DishCategoryController struct {
	svc *service.DishCategoryService
}

func NewDishCategoryController(s *service.DishCategoryService) *DishCategoryController {
	return &DishCategoryController{svc: s}
}

// CreateCategory godoc
// @Summary 创建菜品分类
// @Tags dish-categories
// @Accept json
// @Produce json
// @Security Bearer
// @Param data body model.CreateDishCategoryReq true "分类数据"
// @Success 200 {object} utils.StandardResponse
// @Router /dish-categories [post]
func (c *DishCategoryController) CreateCategory(ctx *gin.Context) {
	storeID := middleware.GetStoreID(ctx)
	var req model.CreateDishCategoryReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if err := c.svc.Create(storeID, &req); err != nil {
		utils.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Success(ctx, nil)
}

// UpdateCategory godoc
// @Summary 更新菜品分类
// @Tags dish-categories
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "分类ID"
// @Param data body model.UpdateDishCategoryReq true "分类数据"
// @Success 200 {object} utils.StandardResponse
// @Router /dish-categories/{id} [put]
func (c *DishCategoryController) UpdateCategory(ctx *gin.Context) {
	storeID := middleware.GetStoreID(ctx)
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "invalid id")
		return
	}
	var req model.UpdateDishCategoryReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if err := c.svc.Update(uint(id), storeID, &req); err != nil {
		utils.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Success(ctx, nil)
}

// DeleteCategory godoc
// @Summary 删除菜品分类
// @Tags dish-categories
// @Security Bearer
// @Param id path int true "分类ID"
// @Success 200 {object} utils.StandardResponse
// @Router /dish-categories/{id} [delete]
func (c *DishCategoryController) DeleteCategory(ctx *gin.Context) {
	storeID := middleware.GetStoreID(ctx)
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "invalid id")
		return
	}
	if err := c.svc.Delete(uint(id), storeID); err != nil {
		utils.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Success(ctx, nil)
}

// ListCategories godoc
// @Summary 分类列表
// @Tags dish-categories
// @Security Bearer
// @Success 200 {object} utils.StandardResponse{data=[]model.DishCategory}
// @Router /dish-categories [get]
func (c *DishCategoryController) ListCategories(ctx *gin.Context) {
	storeID := middleware.GetStoreID(ctx)
	cats, err := c.svc.List(storeID)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Success(ctx, cats)
}

// ListCategoriesWithDishes godoc
// @Summary 分类及菜品列表
// @Tags dish-categories
// @Security Bearer
// @Success 200 {object} utils.StandardResponse{data=[]model.DishCategory}
// @Router /dishes/by-category [get]
func (c *DishCategoryController) ListCategoriesWithDishes(ctx *gin.Context) {
	storeID := middleware.GetStoreID(ctx)
	cats, err := c.svc.ListWithDishes(storeID)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Success(ctx, cats)
}

// ReorderCategories godoc
// @Summary 分类排序
// @Tags dish-categories
// @Accept json
// @Produce json
// @Security Bearer
// @Param data body model.ReorderDishCategoriesReq true "排序数据"
// @Success 200 {object} utils.StandardResponse
// @Router /dish-categories/reorder [post]
func (c *DishCategoryController) ReorderCategories(ctx *gin.Context) {
	storeID := middleware.GetStoreID(ctx)
	var req model.ReorderDishCategoriesReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if len(req.Items) == 0 {
		utils.Success(ctx, nil)
		return
	}
	if err := c.svc.Reorder(storeID, req.Items); err != nil {
		utils.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Success(ctx, nil)
}
