package module

import (
	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/utils/database"
	updatesPkg "github.com/Kevin-Jii/tower-go/utils/updates"
	"github.com/Kevin-Jii/tower-go/utils/search"

	"gorm.io/gorm"
)

var db *gorm.DB

func SetDB(database *gorm.DB) {
	db = database
}

// CreateRole 创建角色
func CreateRole(req *model.Role) (*model.Role, error) {
	if err := db.Create(req).Error; err != nil {
		return nil, err
	}
	return req, nil
}

// UpdateRole 更新角色
func UpdateRole(id uint, req *model.UpdateRoleReq) (*model.Role, error) {
	var role model.Role
	if err := db.First(&role, id).Error; err != nil {
		return nil, err
	}
	updateMap := updatesPkg.BuildUpdatesFromReq(req)
	if len(updateMap) == 0 { // 无字段需要更新直接返回原对象
		return &role, nil
	}
	if err := db.Model(&role).Updates(updateMap).Error; err != nil {
		return nil, err
	}
	// 重新查询确保返回最新（处理可能的数据库触发器或默认值）
	if err := db.First(&role, id).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

// DeleteRole 删除角色
func DeleteRole(id uint) error {
	return db.Delete(&model.Role{}, id).Error
}

// GetRole 获取单个角色
func GetRole(id uint) (*model.Role, error) {
	var role model.Role
	if err := db.First(&role, id).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

// ListRoles 获取角色列表
func ListRoles() ([]model.Role, error) {
	var roles []model.Role
	if err := db.Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

// ListRolesFiltered 带关键字与状态过滤（排除admin）
// keyword 模糊匹配 name/code/description；status 为 0/1 可选
func ListRolesFiltered(keyword string, status *int8) ([]model.Role, error) {
	var roles []model.Role

	qb := database.NewQueryBuilder(db).
		Where("code <> ?", model.RoleCodeAdmin).
		WhereIf(status != nil, "status = ?", func() interface{} {
			if status != nil {
				return *status
			}
			return nil
		}())

	if keyword != "" {
		search.ApplyMultiTermFuzzy(qb, []string{"name", "code", "description"}, keyword, "id")
	}

	qb.OrderByDesc("id")
	if err := qb.Find(&roles); err != nil {
		return nil, err
	}
	return roles, nil
}
