package controller

import (
	"strconv"

	"github.com/Kevin-Jii/tower-go/middleware"
	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/service"
	"github.com/Kevin-Jii/tower-go/utils/http"
	"github.com/gin-gonic/gin"
)

type DictController struct {
	dictService *service.DictService
}

func NewDictController(dictService *service.DictService) *DictController {
	return &DictController{dictService: dictService}
}

// ========== 字典类型 ==========

// CreateType godoc
// @Summary 创建字典类型
// @Tags 字典类型
// @Accept json
// @Produce json
// @Security Bearer
// @Param dict body model.CreateDictTypeReq true "字典类型信息"
// @Success 200 {object} http.Response
// @Router /dict-types [post]
func (c *DictController) CreateType(ctx *gin.Context) {
	if !middleware.IsAdmin(ctx) {
		http.Error(ctx, 403, "仅管理员可以创建字典类型")
		return
	}

	var req model.CreateDictTypeReq
	if !http.BindJSON(ctx, &req) {
		return
	}

	if err := c.dictService.CreateType(&req); err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, nil)
}

// GetType godoc
// @Summary 获取字典类型详情
// @Tags 字典类型
// @Produce json
// @Security Bearer
// @Param id path int true "字典类型ID"
// @Success 200 {object} http.Response{data=model.DictType}
// @Router /dict-types/{id} [get]
func (c *DictController) GetType(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		http.Error(ctx, 400, "无效的ID")
		return
	}

	dictType, err := c.dictService.GetType(uint(id))
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, dictType)
}

// ListTypes godoc
// @Summary 获取字典类型列表
// @Tags 字典类型
// @Produce json
// @Security Bearer
// @Param keyword query string false "关键词"
// @Param status query int false "状态"
// @Success 200 {object} http.Response{data=[]model.DictType}
// @Router /dict-types [get]
func (c *DictController) ListTypes(ctx *gin.Context) {
	keyword := ctx.Query("keyword")
	var status *int8
	if s := ctx.Query("status"); s != "" {
		if v, err := strconv.ParseInt(s, 10, 8); err == nil {
			sv := int8(v)
			status = &sv
		}
	}

	types, err := c.dictService.ListTypes(keyword, status)
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, types)
}

// UpdateType godoc
// @Summary 更新字典类型
// @Tags 字典类型
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "字典类型ID"
// @Param dict body model.UpdateDictTypeReq true "字典类型信息"
// @Success 200 {object} http.Response
// @Router /dict-types/{id} [put]
func (c *DictController) UpdateType(ctx *gin.Context) {
	if !middleware.IsAdmin(ctx) {
		http.Error(ctx, 403, "仅管理员可以更新字典类型")
		return
	}

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		http.Error(ctx, 400, "无效的ID")
		return
	}

	var req model.UpdateDictTypeReq
	if !http.BindJSON(ctx, &req) {
		return
	}

	if err := c.dictService.UpdateType(uint(id), &req); err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, nil)
}

// DeleteType godoc
// @Summary 删除字典类型
// @Tags 字典类型
// @Produce json
// @Security Bearer
// @Param id path int true "字典类型ID"
// @Success 200 {object} http.Response
// @Router /dict-types/{id} [delete]
func (c *DictController) DeleteType(ctx *gin.Context) {
	if !middleware.IsAdmin(ctx) {
		http.Error(ctx, 403, "仅管理员可以删除字典类型")
		return
	}

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		http.Error(ctx, 400, "无效的ID")
		return
	}

	if err := c.dictService.DeleteType(uint(id)); err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, nil)
}

// ========== 字典数据 ==========

// CreateData godoc
// @Summary 创建字典数据
// @Tags 字典数据
// @Accept json
// @Produce json
// @Security Bearer
// @Param dict body model.CreateDictDataReq true "字典数据信息"
// @Success 200 {object} http.Response
// @Router /dict-data [post]
func (c *DictController) CreateData(ctx *gin.Context) {
	if !middleware.IsAdmin(ctx) {
		http.Error(ctx, 403, "仅管理员可以创建字典数据")
		return
	}

	var req model.CreateDictDataReq
	if !http.BindJSON(ctx, &req) {
		return
	}

	if err := c.dictService.CreateData(&req); err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, nil)
}

// GetData godoc
// @Summary 获取字典数据详情
// @Tags 字典数据
// @Produce json
// @Security Bearer
// @Param id path int true "字典数据ID"
// @Success 200 {object} http.Response{data=model.DictData}
// @Router /dict-data/{id} [get]
func (c *DictController) GetData(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		http.Error(ctx, 400, "无效的ID")
		return
	}

	dictData, err := c.dictService.GetData(uint(id))
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, dictData)
}

// ListDataByType godoc
// @Summary 根据类型获取字典数据
// @Tags 字典数据
// @Produce json
// @Security Bearer
// @Param type_code query string false "字典类型编码"
// @Param type_id query int false "字典类型ID"
// @Param status query int false "状态"
// @Success 200 {object} http.Response{data=[]model.DictData}
// @Router /dict-data [get]
func (c *DictController) ListDataByType(ctx *gin.Context) {
	typeCode := ctx.Query("type_code")
	typeIDStr := ctx.Query("type_id")

	if typeCode != "" {
		dataList, err := c.dictService.ListDataByTypeCode(typeCode)
		if err != nil {
			http.Error(ctx, 500, err.Error())
			return
		}
		http.Success(ctx, dataList)
		return
	}

	if typeIDStr != "" {
		typeID, err := strconv.ParseUint(typeIDStr, 10, 32)
		if err != nil {
			http.Error(ctx, 400, "无效的类型ID")
			return
		}
		var status *int8
		if s := ctx.Query("status"); s != "" {
			if v, err := strconv.ParseInt(s, 10, 8); err == nil {
				sv := int8(v)
				status = &sv
			}
		}
		dataList, err := c.dictService.ListDataByTypeID(uint(typeID), status)
		if err != nil {
			http.Error(ctx, 500, err.Error())
			return
		}
		http.Success(ctx, dataList)
		return
	}

	http.Error(ctx, 400, "请提供 type_code 或 type_id")
}

// UpdateData godoc
// @Summary 更新字典数据
// @Tags 字典数据
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "字典数据ID"
// @Param dict body model.UpdateDictDataReq true "字典数据信息"
// @Success 200 {object} http.Response
// @Router /dict-data/{id} [put]
func (c *DictController) UpdateData(ctx *gin.Context) {
	if !middleware.IsAdmin(ctx) {
		http.Error(ctx, 403, "仅管理员可以更新字典数据")
		return
	}

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		http.Error(ctx, 400, "无效的ID")
		return
	}

	var req model.UpdateDictDataReq
	if !http.BindJSON(ctx, &req) {
		return
	}

	if err := c.dictService.UpdateData(uint(id), &req); err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, nil)
}

// DeleteData godoc
// @Summary 删除字典数据
// @Tags 字典数据
// @Produce json
// @Security Bearer
// @Param id path int true "字典数据ID"
// @Success 200 {object} http.Response
// @Router /dict-data/{id} [delete]
func (c *DictController) DeleteData(ctx *gin.Context) {
	if !middleware.IsAdmin(ctx) {
		http.Error(ctx, 403, "仅管理员可以删除字典数据")
		return
	}

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		http.Error(ctx, 400, "无效的ID")
		return
	}

	if err := c.dictService.DeleteData(uint(id)); err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, nil)
}

// GetAllDict godoc
// @Summary 获取所有字典
// @Description 获取所有字典数据，用于前端缓存
// @Tags 字典管理
// @Produce json
// @Security Bearer
// @Success 200 {object} http.Response{data=map[string][]model.DictData}
// @Router /dicts [get]
func (c *DictController) GetAllDict(ctx *gin.Context) {
	result, err := c.dictService.GetAllDict()
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, result)
}
