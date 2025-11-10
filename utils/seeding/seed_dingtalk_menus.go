package seeding

import (
	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/utils/logging"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// InitDingTalkMenuSeeds 初始化钉钉管理菜单种子数据
func InitDingTalkMenuSeeds(db *gorm.DB) error {
	// 检查钉钉菜单是否已存在 (静默模式,不记录 not found)
	var count int64
	if err := db.Model(&model.Menu{}).Where("id = ?", 50).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		logging.LogInfo("钉钉管理菜单已存在，跳过初始化", zap.Int64("count", count))
		return nil
	}

	// 钉钉管理菜单数据
	dingTalkMenus := []model.Menu{
		// 钉钉管理（目录）
		{
			ID:         50,
			ParentID:   0,
			Name:       "dingtalk",
			Title:      "钉钉管理",
			Icon:       "link",
			Type:       1, // 目录
			Sort:       50,
			Permission: "",
			Visible:    1,
			Status:     1,
		},
		// 机器人配置（菜单页面）
		{
			ID:         51,
			ParentID:   50,
			Name:       "dingtalk-robot",
			Title:      "机器人配置",
			Icon:       "robot",
			Path:       "/dingtalk/robot",
			Component:  "dingtalk/robot/index",
			Type:       2, // 菜单
			Sort:       1,
			Permission: "dingtalk:robot:list",
			Visible:    1,
			Status:     1,
		},
		// 操作按钮
		{
			ID:         52,
			ParentID:   51,
			Name:       "dingtalk-robot-add",
			Title:      "新增机器人",
			Type:       3, // 按钮
			Sort:       1,
			Permission: "dingtalk:robot:add",
			Visible:    1,
			Status:     1,
		},
		{
			ID:         53,
			ParentID:   51,
			Name:       "dingtalk-robot-edit",
			Title:      "编辑机器人",
			Type:       3,
			Sort:       2,
			Permission: "dingtalk:robot:edit",
			Visible:    1,
			Status:     1,
		},
		{
			ID:         54,
			ParentID:   51,
			Name:       "dingtalk-robot-delete",
			Title:      "删除机器人",
			Type:       3,
			Sort:       3,
			Permission: "dingtalk:robot:delete",
			Visible:    1,
			Status:     1,
		},
		{
			ID:         55,
			ParentID:   51,
			Name:       "dingtalk-robot-test",
			Title:      "测试推送",
			Type:       3,
			Sort:       4,
			Permission: "dingtalk:robot:test",
			Visible:    1,
			Status:     1,
		},
		{
			ID:         56,
			ParentID:   51,
			Name:       "dingtalk-robot-status",
			Title:      "启用/禁用",
			Type:       3,
			Sort:       5,
			Permission: "dingtalk:robot:status",
			Visible:    1,
			Status:     1,
		},
	}

	// 批量插入
	if err := db.Create(&dingTalkMenus).Error; err != nil {
		logging.LogError("创建钉钉管理菜单失败", zap.Error(err))
		return err
	}

	logging.LogInfo("✅ 钉钉管理菜单创建成功", zap.Int("count", len(dingTalkMenus)))

	// 为管理员角色分配权限
	if err := assignDingTalkMenusToAdmin(db); err != nil {
		logging.LogError("分配钉钉菜单权限失败", zap.Error(err))
		return err
	}

	return nil
}

// assignDingTalkMenusToAdmin 为管理员角色分配钉钉菜单权限
func assignDingTalkMenusToAdmin(db *gorm.DB) error {
	// 查找 admin 角色
	var adminRole model.Role
	if err := db.Where("code = ?", model.RoleCodeAdmin).First(&adminRole).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			logging.LogInfo("admin 角色不存在，跳过权限分配")
			return nil
		}
		return err
	}

	// 查找超级管理员角色 (ID: 999)
	var superAdminRole model.Role
	db.Where("id = ?", 999).First(&superAdminRole)

	// 钉钉菜单 ID 范围: 50-56
	menuIDs := []uint{50, 51, 52, 53, 54, 55, 56}

	// 为 admin 角色分配
	if err := assignMenusToRole(db, adminRole.ID, menuIDs); err != nil {
		return err
	}
	logging.LogInfo("✅ admin 角色钉钉权限分配成功", zap.Uint("role_id", adminRole.ID))

	// 为超级管理员分配 (如果存在)
	if superAdminRole.ID > 0 {
		if err := assignMenusToRole(db, superAdminRole.ID, menuIDs); err != nil {
			return err
		}
		logging.LogInfo("✅ 超级管理员钉钉权限分配成功", zap.Uint("role_id", superAdminRole.ID))
	}

	return nil
}

// assignMenusToRole 为角色分配菜单权限（避免重复）
func assignMenusToRole(db *gorm.DB, roleID uint, menuIDs []uint) error {
	for _, menuID := range menuIDs {
		// 检查是否已存在 (使用 Count 避免 not found 日志)
		var count int64
		if err := db.Model(&model.RoleMenu{}).
			Where("role_id = ? AND menu_id = ?", roleID, menuID).
			Count(&count).Error; err != nil {
			return err
		}

		if count > 0 {
			continue // 已存在，跳过
		}

		// 创建关联
		roleMenu := model.RoleMenu{
			RoleID: roleID,
			MenuID: menuID,
		}
		if err := db.Create(&roleMenu).Error; err != nil {
			return err
		}
	}
	return nil
}
