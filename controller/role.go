package controller

import (
	"net/http"
	"strconv"

	"tower-go/model"
	"tower-go/service"

	"github.com/gin-gonic/gin"
)

// CreateRole 创建角色
// @Summary 创建角色
// @Description 新增一个角色
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param data body model.Role true "角色信息"
// @Success 200 {object} model.Role
// @Failure 400 {object} map[string]string
// @Router /roles [post]
func CreateRole(c *gin.Context) {
	var req model.Role
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	role, err := service.CreateRole(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, role)
}

// UpdateRole 更新角色
// @Summary 更新角色
// @Description 根据ID更新角色信息
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param id path int true "角色ID"
// @Param data body model.Role true "角色信息"
// @Success 200 {object} model.Role
// @Failure 400 {object} map[string]string
// @Router /roles/{id} [put]
func UpdateRole(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req model.Role
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	role, err := service.UpdateRole(uint(id), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, role)
}

// DeleteRole 删除角色
// @Summary 删除角色
// @Description 根据ID删除角色
// @Tags 角色管理
// @Produce json
// @Param id path int true "角色ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /roles/{id} [delete]
func DeleteRole(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := service.DeleteRole(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// GetRole 获取单个角色
// @Summary 获取角色详情
// @Description 根据ID获取角色信息
// @Tags 角色管理
// @Produce json
// @Param id path int true "角色ID"
// @Success 200 {object} model.Role
// @Failure 404 {object} map[string]string
// @Router /roles/{id} [get]
func GetRole(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	role, err := service.GetRole(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, role)
}

// ListRoles 获取角色列表
// @Summary 获取角色列表
// @Description 获取所有角色（不含admin）
// @Tags 角色管理
// @Produce json
// @Success 200 {array} model.Role
// @Router /roles [get]
func ListRoles(c *gin.Context) {
	roles, err := service.ListRoles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// 过滤掉admin角色
	filtered := make([]model.Role, 0, len(roles))
	for _, r := range roles {
		if r.Code != model.RoleCodeAdmin {
			filtered = append(filtered, r)
		}
	}
	c.JSON(http.StatusOK, filtered)
}
