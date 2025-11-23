package examples

import (
	"fmt"
	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/utils/permission"
)

// PermissionExample 权限系统使用示例
func PermissionExample() {
	fmt.Println("=== 菜单权限系统示例 ===\n")

	// 示例1: 设置不同的权限组合
	fmt.Println("1. 权限组合示例:")
	
	// 所有权限
	allPerms := model.PermAll
	fmt.Printf("   所有权限: %d (%s)\n", allPerms, permission.FormatPermissionBits(allPerms))
	
	// 只读权限
	readOnly := model.PermView
	fmt.Printf("   只读权限: %d (%s)\n", readOnly, permission.FormatPermissionBits(readOnly))
	
	// 查看 + 新增
	viewCreate := model.PermView | model.PermCreate
	fmt.Printf("   查看+新增: %d (%s)\n", viewCreate, permission.FormatPermissionBits(viewCreate))
	
	// 查看 + 新增 + 修改
	viewCreateUpdate := model.PermView | model.PermCreate | model.PermUpdate
	fmt.Printf("   查看+新增+修改: %d (%s)\n\n", viewCreateUpdate, permission.FormatPermissionBits(viewCreateUpdate))

	// 示例2: 权限检查
	fmt.Println("2. 权限检查示例:")
	userPerms := uint8(14) // 1110: 查看+新增+修改
	fmt.Printf("   用户权限: %d (%s)\n", userPerms, permission.FormatPermissionBits(userPerms))
	fmt.Printf("   - 可以查看: %v\n", permission.HasViewPermission(userPerms))
	fmt.Printf("   - 可以新增: %v\n", permission.HasCreatePermission(userPerms))
	fmt.Printf("   - 可以修改: %v\n", permission.HasUpdatePermission(userPerms))
	fmt.Printf("   - 可以删除: %v\n\n", permission.HasDeletePermission(userPerms))

	// 示例3: 从二进制字符串解析
	fmt.Println("3. 从二进制字符串解析:")
	permStr := "1010"
	parsed := permission.ParsePermissionString(permStr)
	fmt.Printf("   输入: %s\n", permStr)
	fmt.Printf("   解析结果: %d\n", parsed)
	fmt.Printf("   权限详情: %+v\n\n", permission.GetPermissionDescription(parsed))

	// 示例4: 实际应用场景
	fmt.Println("4. 实际应用场景:")
	
	// 模拟菜单权限
	menus := []struct {
		Name        string
		Permissions uint8
	}{
		{"用户管理", 15}, // 所有权限
		{"报表查看", 8},  // 只读
		{"订单管理", 14}, // 查看+新增+修改
		{"系统设置", 12}, // 查看+新增
	}

	for _, menu := range menus {
		fmt.Printf("   菜单: %s (权限值: %d, 二进制: %s)\n", 
			menu.Name, menu.Permissions, permission.FormatPermissionBits(menu.Permissions))
		desc := permission.GetPermissionDescription(menu.Permissions)
		fmt.Printf("     查看: %v, 新增: %v, 修改: %v, 删除: %v\n",
			desc["view"], desc["create"], desc["update"], desc["delete"])
	}

	fmt.Println("\n=== 示例结束 ===")
}

// AssignRoleMenuExample 角色菜单分配示例
func AssignRoleMenuExample() {
	fmt.Println("=== 角色菜单分配示例 ===\n")

	// 构建权限映射
	perms := map[uint]uint8{
		1: model.PermAll,                                    // 菜单1: 所有权限
		2: model.PermView,                                   // 菜单2: 只读
		3: model.PermView | model.PermCreate,                // 菜单3: 查看+新增
		4: model.PermView | model.PermCreate | model.PermUpdate, // 菜单4: 查看+新增+修改
	}

	fmt.Println("权限分配:")
	for menuID, perm := range perms {
		fmt.Printf("  菜单ID %d: 权限值=%d, 二进制=%s\n", 
			menuID, perm, permission.FormatPermissionBits(perm))
	}

	fmt.Println("\nAPI 请求示例:")
	fmt.Println(`{
  "role_id": 2,
  "menu_ids": [1, 2, 3, 4],
  "perms": {
    "1": 15,
    "2": 8,
    "3": 12,
    "4": 14
  }
}`)

	fmt.Println("\n=== 示例结束 ===")
}
