package controller

import (
	"tower-go/middleware"
	"tower-go/model"
	"tower-go/service"
	httpPkg "tower-go/utils/http"

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
// @Success 200 {object} utils.Response
// @Router /dish-categories [post]
func (c *DishCategoryController) CreateCategory(ctx *gin.Context) {
	storeID := middleware.GetStoreID(ctx)
	var req model.CreateDishCategoryReq
	if !httpPkg.BindJSON(ctx, &req) {
		return
	}
	cat, err := c.svc.Create(storeID, &req)
	if err != nil {
		httpPkg.Error(ctx, 500, err.Error())
		return
	}
	httpPkg.Success(ctx, cat)
}

// CreateCategoryForStore godoc
// @Summary 指定门店创建菜品分类
// @Tags dish-categories
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "门店ID"
// @Param data body model.CreateDishCategoryReq true "分类数据"
// @Success 200 {object} utils.Response{data=model.DishCategory}
// @Router /stores/{id}/dish-categories [post]
func (c *DishCategoryController) CreateCategoryForStore(ctx *gin.Context) {
	storeID, ok := httpPkg.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	currentStoreID := middleware.GetStoreID(ctx)
	if currentStoreID != storeID && !middleware.IsAdmin(ctx) {
		httpPkg.Error(ctx, 403, "forbidden: cross-store create")
		return
	}
	var req model.CreateDishCategoryReq
	if !httpPkg.BindJSON(ctx, &req) {
		return
	}
	cat, err := c.svc.Create(storeID, &req)
	if err != nil {
		httpPkg.Error(ctx, 500, err.Error())
		return
	}
	httpPkg.Success(ctx, cat)
}

// UpdateCategory godoc
// @Summary 更新菜品分类
// @Tags dish-categories
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "分类ID"
// @Param data body model.UpdateDishCategoryReq true "分类数据"
// @Success 200 {object} utils.Response
// @Router /dish-categories/{id} [put]
func (c *DishCategoryController) UpdateCategory(ctx *gin.Context) {
	storeID := middleware.GetStoreID(ctx)
	id, ok := httpPkg.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	var req model.UpdateDishCategoryReq
	if !httpPkg.BindJSON(ctx, &req) {
		return
	}
	if err := c.svc.Update(id, storeID, &req); err != nil {
		status := 500
		if err.Error() == "category not found" {
			status = 404
		} else if err.Error() == "category name already exists in this store" {
			status = 409
		}
		httpPkg.Error(ctx, status, err.Error())
		return
	}
	httpPkg.Success(ctx, nil)
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
// @Success 200 {object} utils.Response
// @Router /stores/{id}/dish-categories/{cid} [put]
func (c *DishCategoryController) UpdateCategoryForStore(ctx *gin.Context) {
	storeID, ok := httpPkg.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	catID, ok := httpPkg.ParseUintParam(ctx, "cid")
	if !ok {
		return
	}
	currentStoreID := middleware.GetStoreID(ctx)
	if currentStoreID != storeID && !middleware.IsAdmin(ctx) {
		httpPkg.Error(ctx, 403, "forbidden: cross-store update")
		return
	}
	var req model.UpdateDishCategoryReq
	if !httpPkg.BindJSON(ctx, &req) {
		return
	}
	if err := c.svc.Update(catID, storeID, &req); err != nil {
		status := 500
		if err.Error() == "category not found" {
			status = 404
		} else if err.Error() == "category name already exists in this store" {
			status = 409
		}
		httpPkg.Error(ctx, status, err.Error())
		return
	}
	httpPkg.Success(ctx, nil)
}

// DeleteCategory godoc
// @Summary 删除菜品分类
// @Tags dish-categories
// @Security Bearer
// @Param id path int true "分类ID"
// @Success 200 {object} utils.Response
// @Router /dish-categories/{id} [delete]
func (c *DishCategoryController) DeleteCategory(ctx *gin.Context) {
	storeID := middleware.GetStoreID(ctx)
	id, ok := httpPkg.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	if err := c.svc.Delete(id, storeID); err != nil {
		status := 500
		if err.Error() == "category has dishes" {
			status = 409
		} else if err.Error() == "category not found" {
			status = 404
		}
		httpPkg.Error(ctx, status, err.Error())
		return
	}
	httpPkg.Success(ctx, nil)
}

// DeleteCategoryForStore godoc
// @Summary 删除指定门店的菜品分类
// @Tags dish-categories
// @Security Bearer
// @Param id path int true "门店ID"
// @Param cid path int true "分类ID"
// @Success 200 {object} utils.Response
// @Router /stores/{id}/dish-categories/{cid} [delete]
func (c *DishCategoryController) DeleteCategoryForStore(ctx *gin.Context) {
	storeID, ok := httpPkg.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	catID, ok := httpPkg.ParseUintParam(ctx, "cid")
	if !ok {
		return
	}
	currentStoreID := middleware.GetStoreID(ctx)
	if currentStoreID != storeID && !middleware.IsAdmin(ctx) {
		httpPkg.Error(ctx, 403, "forbidden: cross-store delete")
		return
	}
	if err := c.svc.Delete(catID, storeID); err != nil {
		status := 500
		if err.Error() == "category has dishes" {
			status = 409
		} else if err.Error() == "category not found" {
			status = 404
		}
		httpPkg.Error(ctx, status, err.Error())
		return
	}
	httpPkg.Success(ctx, nil)
}

// ListDishesForStoreCategory godoc
// @Summary 门店分类下菜品列表
// @Tags dish-categories
// @Security Bearer
// @Param id path int true "门店ID"
// @Param cid path int true "分类ID"
// @Success 200 {object} utils.Response{data=[]model.Dish}
// @Router /stores/{id}/dish-categories/{cid}/dishes [get]
func (c *DishCategoryController) ListDishesForStoreCategory(ctx *gin.Context) {
	storeID, ok := httpPkg.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	catID, ok := httpPkg.ParseUintParam(ctx, "cid")
	if !ok {
		return
	}
	currentStoreID := middleware.GetStoreID(ctx)
	if currentStoreID != storeID && !middleware.IsAdmin(ctx) {
		httpPkg.Error(ctx, 403, "forbidden: cross-store access")
		return
	}
	dishes, err := c.svc.ListDishesForCategory(storeID, catID)
	if err != nil {
		httpPkg.Error(ctx, 404, err.Error())
		return
	}
	httpPkg.Success(ctx, dishes)
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
// @Success 200 {object} utils.Response{data=model.Dish}
// @Router /stores/{id}/dish-categories/{cid}/dishes [post]
func (c *DishCategoryController) CreateDishForStoreCategory(ctx *gin.Context) {
	storeID, ok := httpPkg.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	catID, ok := httpPkg.ParseUintParam(ctx, "cid")
	if !ok {
		return
	}
	currentStoreID := middleware.GetStoreID(ctx)
	if currentStoreID != storeID && !middleware.IsAdmin(ctx) {
		httpPkg.Error(ctx, 403, "forbidden: cross-store create")
		return
	}
	var req model.CreateDishReq
	if !httpPkg.BindJSON(ctx, &req) {
		return
	}
	dish, err := c.svc.CreateDishInCategory(storeID, catID, &req)
	if err != nil {
		status := 500
		if err.Error() == "category not found" {
			status = 404
		} else if err.Error() == "dish name already exists in this category" {
			status = 409
		}
		httpPkg.Error(ctx, status, err.Error())
		return
	}
	httpPkg.Success(ctx, dish)
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
// @Success 200 {object} utils.Response
// @Router /stores/{id}/dish-categories/{cid}/dishes/{did} [put]
func (c *DishCategoryController) UpdateDishForStoreCategory(ctx *gin.Context) {
	storeID, ok := httpPkg.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	catID, ok := httpPkg.ParseUintParam(ctx, "cid")
	if !ok {
		return
	}
	dishID, ok := httpPkg.ParseUintParam(ctx, "did")
	if !ok {
		return
	}
	currentStoreID := middleware.GetStoreID(ctx)
	if currentStoreID != storeID && !middleware.IsAdmin(ctx) {
		httpPkg.Error(ctx, 403, "forbidden: cross-store update")
		return
	}
	var req model.UpdateDishReq
	if !httpPkg.BindJSON(ctx, &req) {
		return
	}
	if err := c.svc.UpdateDishInCategory(storeID, catID, dishID, &req); err != nil {
		status := 500
		if err.Error() == "category not found" || err.Error() == "dish not found" {
			status = 404
		} else if err.Error() == "dish name already exists in this category" {
			status = 409
		} else if err.Error() == "dish does not belong to this category" {
			status = 400
		}
		httpPkg.Error(ctx, status, err.Error())
		return
	}
	httpPkg.Success(ctx, nil)
}

// DeleteDishForStoreCategory godoc
// @Summary 删除门店指定分类下的菜品
// @Tags dishes
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "门店ID"
// @Param cid path int true "分类ID"
// @Param did path int true "菜品ID"
// @Success 200 {object} utils.Response
// @Router /stores/{id}/dish-categories/{cid}/dishes/{did} [delete]
func (c *DishCategoryController) DeleteDishForStoreCategory(ctx *gin.Context) {
	storeID, ok := httpPkg.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	catID, ok := httpPkg.ParseUintParam(ctx, "cid")
	if !ok {
		return
	}
	dishID, ok := httpPkg.ParseUintParam(ctx, "did")
	if !ok {
		return
	}
	currentStoreID := middleware.GetStoreID(ctx)
	if currentStoreID != storeID && !middleware.IsAdmin(ctx) {
		httpPkg.Error(ctx, 403, "forbidden: cross-store delete")
		return
	}
	if err := c.svc.DeleteDishInCategory(storeID, catID, dishID); err != nil {
		status := 500
		if err.Error() == "category not found" || err.Error() == "dish not found" {
			status = 404
		} else if err.Error() == "dish does not belong to this category" {
			status = 400
		}
		httpPkg.Error(ctx, status, err.Error())
		return
	}
	httpPkg.Success(ctx, nil)
}

// ListCategories godoc
// @Summary 分类列表
// @Tags dish-categories
// @Security Bearer
// @Success 200 {object} utils.Response{data=[]model.DishCategory}
// @Router /dish-categories [get]
func (c *DishCategoryController) ListCategories(ctx *gin.Context) {
	storeID := middleware.GetStoreID(ctx)
	cats, err := c.svc.List(storeID)
	if err != nil {
		httpPkg.Error(ctx, 500, err.Error())
		return
	}
	httpPkg.Success(ctx, cats)
}

// ListCategoriesForStore godoc
// @Summary 指定门店分类列表
// @Tags dish-categories
// @Security Bearer
// @Param id path int true "门店ID"
// @Success 200 {object} utils.Response{data=[]model.DishCategory}
// @Router /stores/{id}/dish-categories [get]
func (c *DishCategoryController) ListCategoriesForStore(ctx *gin.Context) {
	storeID, ok := httpPkg.ParseUintParam(ctx, "id")
	if !ok {
		return
	}

	// 权限: 如果请求的 store 与当前 token 中的 store 不同, 仅允许总部管理员访问
	currentStoreID := middleware.GetStoreID(ctx)
	if currentStoreID != storeID && !middleware.IsAdmin(ctx) {
		httpPkg.Error(ctx, 403, "forbidden: cross-store access")
		return
	}

	cats, err := c.svc.List(storeID)
	if err != nil {
		httpPkg.Error(ctx, 500, err.Error())
		return
	}
	httpPkg.Success(ctx, cats)
}

// ListCategoriesWithDishes godoc
// @Summary 分类及菜品列表
// @Tags dish-categories
// @Security Bearer
// @Success 200 {object} utils.Response{data=[]model.DishCategory}
// @Router /dishes/by-category [get]
func (c *DishCategoryController) ListCategoriesWithDishes(ctx *gin.Context) {
	storeID := middleware.GetStoreID(ctx)
	cats, err := c.svc.ListWithDishes(storeID)
	if err != nil {
		httpPkg.Error(ctx, 500, err.Error())
		return
	}
	httpPkg.Success(ctx, cats)
}

// ReorderCategories godoc
// @Summary 分类排序
// @Tags dish-categories
// @Accept json
// @Produce json
// @Security Bearer
// @Param data body model.ReorderDishCategoriesReq true "排序数据"
// @Success 200 {object} utils.Response
// @Router /dish-categories/reorder [post]
func (c *DishCategoryController) ReorderCategories(ctx *gin.Context) {
	storeID := middleware.GetStoreID(ctx)
	var req model.ReorderDishCategoriesReq
	if !httpPkg.BindJSON(ctx, &req) {
		return
	}
	if len(req.Items) == 0 {
		httpPkg.Success(ctx, nil)
		return
	}
	if err := c.svc.Reorder(storeID, req.Items); err != nil {
		httpPkg.Error(ctx, 500, err.Error())
		return
	}
	httpPkg.Success(ctx, nil)
}
