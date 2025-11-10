package main

import (
	"fmt"
	"github.com/Kevin-Jii/tower-go/bootstrap"
	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/utils/database"
)

func main() {
	// 初始化日志
	bootstrap.InitLogger()

	// 加载配置
	bootstrap.LoadAppConfig()

	// 初始化数据库
	bootstrap.InitDatabase()

	fmt.Println("开始强制初始化钉钉菜单...")

	// 1. 删除旧的钉钉菜单 (如果存在)
	fmt.Println("1. 清理旧数据...")
	database.DB.Exec("DELETE FROM role_menus WHERE menu_id BETWEEN 50 AND 56")
	database.DB.Exec("DELETE FROM menus WHERE id BETWEEN 50 AND 56")
	fmt.Println("✅ 旧数据清理完成")

	// 2. 创建钉钉菜单
	fmt.Println("2. 创建钉钉菜单...")
	dingTalkMenus := []model.Menu{
		// 钉钉管理（目录）
		{
			ID:         50,
			ParentID:   0,
			Name:       "dingtalk",
			Title:      "钉钉管理",
			Icon:       "link",
			Type:       1,
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
			Type:       2,
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
			Type:       3,
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

	if err := database.DB.Create(&dingTalkMenus).Error; err != nil {
		fmt.Printf("❌ 创建菜单失败: %v\n", err)
		return
	}
	fmt.Printf("✅ 创建了 %d 个菜单\n", len(dingTalkMenus))

	// 3. 分配权限给 admin 角色
	fmt.Println("3. 分配权限...")
	var adminRole model.Role
	if err := database.DB.Where("code = ?", "admin").First(&adminRole).Error; err != nil {
		fmt.Printf("❌ 查找 admin 角色失败: %v\n", err)
		return
	}

	menuIDs := []uint{50, 51, 52, 53, 54, 55, 56}
	for _, menuID := range menuIDs {
		roleMenu := model.RoleMenu{
			RoleID: adminRole.ID,
			MenuID: menuID,
		}
		if err := database.DB.Create(&roleMenu).Error; err != nil {
			fmt.Printf("⚠️ 分配权限失败 (menu_id=%d): %v\n", menuID, err)
		}
	}
	fmt.Printf("✅ 为 admin 角色分配了 %d 个权限\n", len(menuIDs))

	// 4. 分配权限给超级管理员 (ID: 999)
	var superAdmin model.Role
	if err := database.DB.Where("id = ?", 999).First(&superAdmin).Error; err == nil {
		for _, menuID := range menuIDs {
			roleMenu := model.RoleMenu{
				RoleID: 999,
				MenuID: menuID,
			}
			database.DB.Create(&roleMenu)
		}
		fmt.Printf("✅ 为超级管理员分配了 %d 个权限\n", len(menuIDs))
	}

	// 5. 验证
	fmt.Println("\n4. 验证结果...")
	var count int64
	database.DB.Model(&model.Menu{}).Where("id BETWEEN ? AND ?", 50, 56).Count(&count)
	fmt.Printf("✅ 菜单数量: %d\n", count)

	database.DB.Model(&model.RoleMenu{}).
		Where("role_id = ? AND menu_id BETWEEN ? AND ?", adminRole.ID, 50, 56).
		Count(&count)
	fmt.Printf("✅ admin 权限数量: %d\n", count)

	fmt.Println("\n🎉 初始化完成! 请重启服务器并刷新前端页面。")
}
