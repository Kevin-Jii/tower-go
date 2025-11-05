package controller

import (
	"net/http"
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
	if !utils.BindJSON(ctx, &req) {
		return
	}
	cat, err := c.svc.Create(storeID, &req)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Success(ctx, cat)
}

// CreateCategoryForStore godoc
// @Summary 指定门店创建菜品分类
// @Tags dish-categories
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "门店ID"
// @Param data body model.CreateDishCategoryReq true "分类数据"
// @Success 200 {object} utils.StandardResponse{data=model.DishCategory}
// @Router /stores/{id}/dish-categories [post]
func (c *DishCategoryController) CreateCategoryForStore(ctx *gin.Context) {
	storeID, ok := utils.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	currentStoreID := middleware.GetStoreID(ctx)
	if currentStoreID != storeID && !middleware.IsAdmin(ctx) {
		utils.Error(ctx, http.StatusForbidden, "forbidden: cross-store create")
		return
	}
	var req model.CreateDishCategoryReq
	if !utils.BindJSON(ctx, &req) {
		return
	}
	cat, err := c.svc.Create(storeID, &req)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Success(ctx, cat)
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
	id, ok := utils.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	var req model.UpdateDishCategoryReq
	if !utils.BindJSON(ctx, &req) {
		return
	}
	if err := c.svc.Update(id, storeID, &req); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "category not found" {
			status = http.StatusNotFound
		} else if err.Error() == "category name already exists in this store" {
			status = http.StatusConflict
		}
		utils.Error(ctx, status, err.Error())
		return
	}
	utils.Success(ctx, nil)
}

// UpdateCategoryForStore 门店作用域下更新分类
// @Summary 更新指定门店的菜品分类
// @Tags dish-categories
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "门店ID"
// @Param cid path int true "分类ID"
// @Param data body model.UpdateDishCategoryReq true "分类更新数据"
// @Success 200 {object} utils.StandardResponse
// @Router /stores/{id}/dish-categories/{cid} [put]
func (c *DishCategoryController) UpdateCategoryForStore(ctx *gin.Context) {
	storeID, ok := utils.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	catID, ok := utils.ParseUintParam(ctx, "cid")
	if !ok {
		return
	}
	currentStoreID := middleware.GetStoreID(ctx)
	if currentStoreID != storeID && !middleware.IsAdmin(ctx) {
		utils.Error(ctx, http.StatusForbidden, "forbidden: cross-store update")
		return
	}
	var req model.UpdateDishCategoryReq
	if !utils.BindJSON(ctx, &req) {
		return
	}
	if err := c.svc.Update(catID, storeID, &req); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "category not found" {
			status = http.StatusNotFound
		} else if err.Error() == "category name already exists in this store" {
			status = http.StatusConflict
		}
		utils.Error(ctx, status, err.Error())
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
	id, ok := utils.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	if err := c.svc.Delete(id, storeID); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "category has dishes" {
			status = http.StatusConflict
		} else if err.Error() == "category not found" {
			status = http.StatusNotFound
		}
		utils.Error(ctx, status, err.Error())
		return
	}
	utils.Success(ctx, nil)
}

// DeleteCategoryForStore godoc
// @Summary 删除指定门店的菜品分类
// @Tags dish-categories
// @Security Bearer
// @Param id path int true "门店ID"
// @Param cid path int true "分类ID"
// @Success 200 {object} utils.StandardResponse
// @Router /stores/{id}/dish-categories/{cid} [delete]
func (c *DishCategoryController) DeleteCategoryForStore(ctx *gin.Context) {
	storeID, ok := utils.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	catID, ok := utils.ParseUintParam(ctx, "cid")
	if !ok {
		return
	}
	currentStoreID := middleware.GetStoreID(ctx)
	if currentStoreID != storeID && !middleware.IsAdmin(ctx) {
		utils.Error(ctx, http.StatusForbidden, "forbidden: cross-store delete")
		return
	}
	if err := c.svc.Delete(catID, storeID); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "category has dishes" {
			status = http.StatusConflict
		} else if err.Error() == "category not found" {
			status = http.StatusNotFound
		}
		utils.Error(ctx, status, err.Error())
		return
	}
	utils.Success(ctx, nil)
}

// ListDishesForStoreCategory godoc
// @Summary 门店分类下菜品列表
// @Tags dish-categories
// @Security Bearer
// @Param id path int true "门店ID"
// @Param cid path int true "分类ID"
// @Success 200 {object} utils.StandardResponse{data=[]model.Dish}
// @Router /stores/{id}/dish-categories/{cid}/dishes [get]
func (c *DishCategoryController) ListDishesForStoreCategory(ctx *gin.Context) {
	storeID, ok := utils.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	catID, ok := utils.ParseUintParam(ctx, "cid")
	if !ok {
		return
	}
	currentStoreID := middleware.GetStoreID(ctx)
	if currentStoreID != storeID && !middleware.IsAdmin(ctx) {
		utils.Error(ctx, http.StatusForbidden, "forbidden: cross-store access")
		return
	}
	dishes, err := c.svc.ListDishesForCategory(storeID, catID)
	if err != nil {
		utils.Error(ctx, http.StatusNotFound, err.Error())
		return
	}
	utils.Success(ctx, dishes)
}

// CreateDishForStoreCategory godoc
// @Summary 在门店指定分类下创建菜品
// @Tags dishes
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "门店ID"
// @Param cid path int true "分类ID"
// @Param data body model.CreateDishReq true "菜品数据 (category_id 忽略, 以路径分类为准)"
// @Success 200 {object} utils.StandardResponse{data=model.Dish}
// @Router /stores/{id}/dish-categories/{cid}/dishes [post]
func (c *DishCategoryController) CreateDishForStoreCategory(ctx *gin.Context) {
	storeID, ok := utils.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	catID, ok := utils.ParseUintParam(ctx, "cid")
	if !ok {
		return
	}
	currentStoreID := middleware.GetStoreID(ctx)
	if currentStoreID != storeID && !middleware.IsAdmin(ctx) {
		utils.Error(ctx, http.StatusForbidden, "forbidden: cross-store create")
		return
	}
	var req model.CreateDishReq
	if !utils.BindJSON(ctx, &req) {
		return
	}
	dish, err := c.svc.CreateDishInCategory(storeID, catID, &req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "category not found" {
			status = http.StatusNotFound
		} else if err.Error() == "dish name already exists in this category" {
			status = http.StatusConflict
		}
		utils.Error(ctx, status, err.Error())
		return
	}
	utils.Success(ctx, dish)
}

// UpdateDishForStoreCategory godoc
// @Summary 在门店指定分类下更新菜品
// @Tags dishes
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "门店ID"
// @Param cid path int true "分类ID"
// @Param did path int true "菜品ID"
// @Param data body model.UpdateDishReq true "菜品更新数据（忽略 category_id 防止跨分类）"
// @Success 200 {object} utils.StandardResponse
// @Router /stores/{id}/dish-categories/{cid}/dishes/{did} [put]
func (c *DishCategoryController) UpdateDishForStoreCategory(ctx *gin.Context) {
	storeID, ok := utils.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	catID, ok := utils.ParseUintParam(ctx, "cid")
	if !ok {
		return
	}
	dishID, ok := utils.ParseUintParam(ctx, "did")
	if !ok {
		return
	}
	currentStoreID := middleware.GetStoreID(ctx)
	if currentStoreID != storeID && !middleware.IsAdmin(ctx) {
		utils.Error(ctx, http.StatusForbidden, "forbidden: cross-store update")
		return
	}
	var req model.UpdateDishReq
	if !utils.BindJSON(ctx, &req) {
		return
	}
	if err := c.svc.UpdateDishInCategory(storeID, catID, dishID, &req); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "category not found" || err.Error() == "dish not found" {
			status = http.StatusNotFound
		} else if err.Error() == "dish name already exists in this category" {
			status = http.StatusConflict
		} else if err.Error() == "dish does not belong to this category" {
			status = http.StatusBadRequest
		}
		utils.Error(ctx, status, err.Error())
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

// ListCategoriesForStore godoc
// @Summary 指定门店分类列表
// @Tags dish-categories
// @Security Bearer
// @Param id path int true "门店ID"
// @Success 200 {object} utils.StandardResponse{data=[]model.DishCategory}
// @Router /stores/{id}/dish-categories [get]
func (c *DishCategoryController) ListCategoriesForStore(ctx *gin.Context) {
	storeID, ok := utils.ParseUintParam(ctx, "id")
	if !ok {
		return
	}

	// 权限: 如果请求的 store 与当前 token 中的 store 不同, 仅允许总部管理员访问
	currentStoreID := middleware.GetStoreID(ctx)
	if currentStoreID != storeID && !middleware.IsAdmin(ctx) {
		utils.Error(ctx, http.StatusForbidden, "forbidden: cross-store access")
		return
	}

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
	if !utils.BindJSON(ctx, &req) {
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
