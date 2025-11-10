package controller

import (
	"strconv"

	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/service"
	"github.com/Kevin-Jii/tower-go/utils/http"

	"github.com/gin-gonic/gin"
)

// CreateRole 创建角色
// @Summary 创建角色
// @Description 新增一个角色，status：1=启用 0=禁用（默认为1）
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param data body model.Role true "角色信息"
// @Success 200 {object} http.Response{data=model.Role}
// @Failure 400 {object} map[string]string
// @Router /roles [post]
func CreateRole(c *gin.Context) {
	var req model.Role
	if err := c.ShouldBindJSON(&req); err != nil {
		http.Error(c, 400, err.Error())
		return
	}
	role, err := service.CreateRole(&req)
	if err != nil {
		http.Error(c, 500, err.Error())
		return
	}
	http.Success(c, role)
}

// UpdateRole 局部更新角色
// @Summary 更新角色
// @Description 根据ID更新角色信息，可局部更新，仅传需要修改的字段；status=1启用 0禁用
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param id path int true "角色ID"
// @Param data body model.UpdateRoleReq true "仅传需要修改的字段 (name, code, description, status)，未传字段保持原值"
// @Success 200 {object} http.Response{data=model.Role}
// @Failure 400 {object} map[string]string
// @Router /roles/{id} [put]
func UpdateRole(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(c, 400, "invalid id")
		return
	}
	var req model.UpdateRoleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		http.Error(c, 400, err.Error())
		return
	}
	role, err := service.UpdateRole(uint(id), &req)
	if err != nil {
		http.Error(c, 500, err.Error())
		return
	}
	http.Success(c, role)
}

// DeleteRole 删除角色
// @Summary 删除角色
// @Description 根据ID删除角色
// @Tags 角色管理
// @Produce json
// @Param id path int true "角色ID"
// @Success 200 {object} http.Response
// @Failure 400 {object} map[string]string
// @Router /roles/{id} [delete]
func DeleteRole(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(c, 400, "invalid id")
		return
	}
	if err := service.DeleteRole(uint(id)); err != nil {
		http.Error(c, 500, err.Error())
		return
	}
	http.Success(c, nil)
}

// GetRole 获取单个角色
// @Summary 获取角色详情
// @Description 根据ID获取角色信息
// @Tags 角色管理
// @Produce json
// @Param id path int true "角色ID"
// @Success 200 {object} http.Response{data=model.Role}
// @Failure 404 {object} map[string]string
// @Router /roles/{id} [get]
func GetRole(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(c, 400, "invalid id")
		return
	}
	role, err := service.GetRole(uint(id))
	if err != nil {
		http.Error(c, 404, err.Error())
		return
	}
	http.Success(c, role)
}

// ListRoles 获取角色列表
// @Summary 获取角色列表
// @Description 获取角色列表（排除admin），支持 keyword 模糊查询 name/code/description，status=0|1 过滤
// @Tags 角色管理
// @Produce json
// @Param keyword query string false "关键字(模糊匹配 name/code/description)"
// @Param status query int false "状态过滤 1=启用 0=禁用"
// @Success 200 {object} http.Response{data=[]model.Role}
// @Router /roles [get]
func ListRoles(c *gin.Context) {
	keyword := c.Query("keyword")
	statusStr := c.Query("status")
	var statusPtr *int8
	if statusStr != "" {
		if statusStr == "0" || statusStr == "1" {
			var v int8
			if statusStr == "1" {
				v = 1
			} else {
				v = 0
			}
			statusPtr = &v
		} else {
			http.Error(c, 400, "invalid status, must be 0 or 1")
			return
		}
	}
	roles, err := service.ListRolesFiltered(keyword, statusPtr)
	if err != nil {
		http.Error(c, 500, err.Error())
		return
	}
	http.Success(c, roles)
}
