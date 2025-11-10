package seeding

import (
	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/utils"
	"github.com/Kevin-Jii/tower-go/utils/auth"
	"github.com/Kevin-Jii/tower-go/utils/logging"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// InitMenuSeeds 初始化菜单种子数据
func InitMenuSeeds(db *gorm.DB) error {
	// 检查是否已有菜单数据
	var count int64
	if err := db.Model(&model.Menu{}).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		// 已有数据，跳过初始化
		return nil
	}

	// 基础菜单数据（使用 TDesign 图标）
	menus := []model.Menu{
		// 系统管理（目录）
		{ID: 1, ParentID: 0, Name: "system", Title: "系统管理", Icon: "setting", Type: 1, Sort: 1, Visible: 1, Status: 1},
		{ID: 2, ParentID: 1, Name: "user", Title: "用户管理", Icon: "user", Path: "/system/user", Component: "system/user/index", Type: 2, Sort: 1, Permission: "system:user:list", Visible: 1, Status: 1},
		{ID: 3, ParentID: 2, Name: "user-add", Title: "新增用户", Type: 3, Permission: "system:user:add", Visible: 1, Status: 1},
		{ID: 4, ParentID: 2, Name: "user-edit", Title: "编辑用户", Type: 3, Permission: "system:user:edit", Visible: 1, Status: 1},
		{ID: 5, ParentID: 2, Name: "user-delete", Title: "删除用户", Type: 3, Permission: "system:user:delete", Visible: 1, Status: 1},

		{ID: 6, ParentID: 1, Name: "role", Title: "角色管理", Icon: "usergroup", Path: "/system/role", Component: "system/role/index", Type: 2, Sort: 2, Permission: "system:role:list", Visible: 1, Status: 1},
		{ID: 7, ParentID: 6, Name: "role-add", Title: "新增角色", Type: 3, Permission: "system:role:add", Visible: 1, Status: 1},
		{ID: 8, ParentID: 6, Name: "role-edit", Title: "编辑角色", Type: 3, Permission: "system:role:edit", Visible: 1, Status: 1},
		{ID: 9, ParentID: 6, Name: "role-delete", Title: "删除角色", Type: 3, Permission: "system:role:delete", Visible: 1, Status: 1},
		{ID: 10, ParentID: 6, Name: "role-menu", Title: "分配菜单", Type: 3, Permission: "system:role:menu", Visible: 1, Status: 1},

		{ID: 11, ParentID: 1, Name: "menu", Title: "菜单管理", Icon: "menu-fold", Path: "/system/menu", Component: "system/menu/index", Type: 2, Sort: 3, Permission: "system:menu:list", Visible: 1, Status: 1},
		{ID: 12, ParentID: 11, Name: "menu-add", Title: "新增菜单", Type: 3, Permission: "system:menu:add", Visible: 1, Status: 1},
		{ID: 13, ParentID: 11, Name: "menu-edit", Title: "编辑菜单", Type: 3, Permission: "system:menu:edit", Visible: 1, Status: 1},
		{ID: 14, ParentID: 11, Name: "menu-delete", Title: "删除菜单", Type: 3, Permission: "system:menu:delete", Visible: 1, Status: 1},

		// 门店管理（目录）
		{ID: 20, ParentID: 0, Name: "store", Title: "门店管理", Icon: "shop", Type: 1, Sort: 2, Visible: 1, Status: 1},
		{ID: 21, ParentID: 20, Name: "store-list", Title: "门店列表", Icon: "view-list", Path: "/store/list", Component: "store/list/index", Type: 2, Sort: 1, Permission: "store:list", Visible: 1, Status: 1},
		{ID: 22, ParentID: 21, Name: "store-add", Title: "新增门店", Type: 3, Permission: "store:add", Visible: 1, Status: 1},
		{ID: 23, ParentID: 21, Name: "store-edit", Title: "编辑门店", Type: 3, Permission: "store:edit", Visible: 1, Status: 1},
		{ID: 24, ParentID: 21, Name: "store-delete", Title: "删除门店", Type: 3, Permission: "store:delete", Visible: 1, Status: 1},
		{ID: 25, ParentID: 21, Name: "store-menu", Title: "配置权限", Type: 3, Permission: "store:menu", Visible: 1, Status: 1},

		// 菜品管理（目录）
		{ID: 30, ParentID: 0, Name: "dish", Title: "菜品管理", Icon: "food", Type: 1, Sort: 3, Visible: 1, Status: 1},
		{ID: 31, ParentID: 30, Name: "dish-list", Title: "菜品列表", Icon: "view-list", Path: "/dish/list", Component: "dish/list/index", Type: 2, Sort: 1, Permission: "dish:list", Visible: 1, Status: 1},
		{ID: 32, ParentID: 31, Name: "dish-add", Title: "新增菜品", Type: 3, Permission: "dish:add", Visible: 1, Status: 1},
		{ID: 33, ParentID: 31, Name: "dish-edit", Title: "编辑菜品", Type: 3, Permission: "dish:edit", Visible: 1, Status: 1},
		{ID: 34, ParentID: 31, Name: "dish-delete", Title: "删除菜品", Type: 3, Permission: "dish:delete", Visible: 1, Status: 1},
		{ID: 35, ParentID: 31, Name: "dish-status", Title: "上下架", Type: 3, Permission: "dish:status", Visible: 1, Status: 1},

		// 报菜管理（目录）
		{ID: 40, ParentID: 0, Name: "report", Title: "报菜管理", Icon: "file-paste", Type: 1, Sort: 4, Visible: 1, Status: 1},
		{ID: 41, ParentID: 40, Name: "report-list", Title: "报菜记录", Icon: "view-list", Path: "/report/list", Component: "report/list/index", Type: 2, Sort: 1, Permission: "report:list", Visible: 1, Status: 1},
		{ID: 42, ParentID: 41, Name: "report-add", Title: "创建报菜", Type: 3, Permission: "report:add", Visible: 1, Status: 1},
		{ID: 43, ParentID: 41, Name: "report-edit", Title: "编辑报菜", Type: 3, Permission: "report:edit", Visible: 1, Status: 1},
		{ID: 44, ParentID: 41, Name: "report-delete", Title: "删除报菜", Type: 3, Permission: "report:delete", Visible: 1, Status: 1},

		{ID: 45, ParentID: 40, Name: "report-stats", Title: "数据统计", Icon: "chart-bar", Path: "/report/statistics", Component: "report/statistics/index", Type: 2, Sort: 2, Permission: "report:statistics", Visible: 1, Status: 1},
	}

	// 批量插入
	if err := db.Create(&menus).Error; err != nil {
		return err
	}

	return nil
}

// InitRoleSeeds 初始化角色种子数据（确保 admin / store_admin / staff 存在）
func InitRoleSeeds(db *gorm.DB) error {
	// 针对每个角色单独检查，避免已有其它角色但缺失目标角色时无法补全
	ensure := func(code, name, desc string) error {
		var role model.Role
		if err := db.Where("code = ?", code).First(&role).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				// 创建缺失角色
				role = model.Role{Name: name, Code: code, Description: desc}
				if err2 := db.Create(&role).Error; err2 != nil {
					return err2
				}
				return nil
			}
			return err
		}
		return nil
	}

	if err := ensure(model.RoleCodeAdmin, "总部管理员", "系统最高权限角色"); err != nil {
		return err
	}
	if err := ensure(model.RoleCodeStoreAdmin, "门店管理员", "门店维度管理权限角色"); err != nil {
		return err
	}
	if err := ensure(model.RoleCodeStaff, "普通员工", "基础操作权限角色"); err != nil {
		return err
	}
	return nil
}

// InitSuperAdmin 初始化超级管理员
func InitSuperAdmin(db *gorm.DB) error {
	// 检查超级管理员是否已存在
	var existingUser model.User
	if err := db.Where("id = ?", 999).First(&existingUser).Error; err == nil {
		// 超级管理员已存在
		logging.LogInfo("超级管理员已存在，跳过创建", zap.Uint("user_id", existingUser.ID))
		return nil
	}

	// 1. 创建/检查超级门店（ID=999）
	var superStore model.Store
	if err := db.Where("id = ?", 999).First(&superStore).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			superStore = model.Store{
				ID:            999,
				Name:          "总部",
				Address:       "系统默认总部",
				Phone:         "13082848180",
				BusinessHours: "全天",
				Status:        1,
				ContactPerson: "超级管理员",
				Remark:        "系统默认总部门店",
			}
			if err := db.Create(&superStore).Error; err != nil {
				logging.LogError("创建超级门店失败", zap.Error(err))
				return err
			}
			logging.LogInfo("超级门店创建成功", zap.Uint("store_id", superStore.ID))
		} else {
			return err
		}
	}

	// 2. 确保超级管理员角色存在（ID=999）
	var superRole model.Role
	if err := db.Where("id = ?", 999).First(&superRole).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// 查找 admin 角色
			var adminRole model.Role
			if err := db.Where("code = ?", model.RoleCodeAdmin).First(&adminRole).Error; err != nil {
				logging.LogError("查找admin角色失败", zap.Error(err))
				return err
			}

			// 创建超级管理员角色
			superRole = model.Role{
				ID:          999,
				Name:        "超级管理员",
				Code:        "super_admin",
				Description: "系统最高权限，不可删除",
			}
			if err := db.Create(&superRole).Error; err != nil {
				logging.LogError("创建超级管理员角色失败", zap.Error(err))
				return err
			}
			logging.LogInfo("超级管理员角色创建成功", zap.Uint("role_id", superRole.ID))
		} else {
			return err
		}
	}

	// 3. 创建超级管理员用户
	// 生成强随机密码作为默认密码
	defaultPassword, err := auth.GenerateStrongPassword(16)
	if err != nil {
		logging.LogError("生成默认密码失败", zap.Error(err))
		return err
	}

	hashedPassword, err := auth.HashPassword(defaultPassword)
	if err != nil {
		logging.LogError("密码加密失败", zap.Error(err))
		return err
	}

	// 使用原始 SQL 插入，避免 LastLoginAt 零值问题
	// 超级管理员固定工号为 999999
	sql := `INSERT INTO users (id, username, phone, password, nickname, email, employee_no, store_id, role_id, status, gender, created_at, updated_at) 
	        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())`

	if err := db.Exec(sql, 999, "admin", "13082848180", hashedPassword, "超级管理员", "admin@tower.com", "999999", 999, 999, 1, 1).Error; err != nil {
		logging.LogError("创建超级管理员失败", zap.Error(err))
		return err
	}

	logging.LogInfo("✅ 超级管理员创建成功",
		zap.Uint("user_id", 999),
		zap.String("username", "admin"),
		zap.String("phone", "13082848180"),
		zap.String("employee_no", "999999"),
		zap.String("default_password", defaultPassword),
	)
	logging.LogWarn("⚠️  重要提示：请立即修改默认密码，首次登录后请修改为自定义密码！")

	return nil
} // InitRoleMenuSeeds 初始化角色菜单关联（默认权限）
func InitRoleMenuSeeds(db *gorm.DB) error {
	// 检查是否已有数据
	var count int64
	if err := db.Model(&model.RoleMenu{}).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	// 查询角色
	var adminRole, storeAdminRole, staffRole model.Role
	db.Where("code = ?", model.RoleCodeAdmin).First(&adminRole)
	db.Where("code = ?", model.RoleCodeStoreAdmin).First(&storeAdminRole)
	db.Where("code = ?", model.RoleCodeStaff).First(&staffRole)

	if adminRole.ID == 0 || storeAdminRole.ID == 0 || staffRole.ID == 0 {
		// 角色未初始化，跳过
		return nil
	}

	// 总部管理员：所有权限（1-45）
	adminMenus := make([]model.RoleMenu, 0)
	for i := uint(1); i <= 45; i++ {
		adminMenus = append(adminMenus, model.RoleMenu{RoleID: adminRole.ID, MenuID: i})
	}

	// 门店管理员：门店、菜品、报菜管理权限（20-45）
	storeAdminMenus := make([]model.RoleMenu, 0)
	for i := uint(20); i <= 45; i++ {
		storeAdminMenus = append(storeAdminMenus, model.RoleMenu{RoleID: storeAdminRole.ID, MenuID: i})
	}

	// 普通员工：菜品和报菜管理权限（30-45，但不含删除）
	staffMenuIDs := []uint{30, 31, 32, 33, 35, 40, 41, 42, 43, 45}
	staffMenus := make([]model.RoleMenu, 0)
	for _, menuID := range staffMenuIDs {
		staffMenus = append(staffMenus, model.RoleMenu{RoleID: staffRole.ID, MenuID: menuID})
	}

	// 批量插入
	allMenus := append(adminMenus, storeAdminMenus...)
	allMenus = append(allMenus, staffMenus...)

	if err := db.Create(&allMenus).Error; err != nil {
		return err
	}

	return nil
}

// EnsureStoreCodes 为缺失编码的门店生成编码
func EnsureStoreCodes(db *gorm.DB) error {
	var stores []model.Store
	if err := db.Where("store_code IS NULL OR store_code = ''").Find(&stores).Error; err != nil {
		logging.LogError("查询缺失门店编码失败", zap.Error(err))
		return err
	}

	if len(stores) == 0 {
		logging.LogInfo("所有门店均已具备编码")
		return nil
	}

	for _, store := range stores {
		code, err := utils.GenerateStoreCode(db)
		if err != nil {
			return err
		}
		if err := db.Model(&model.Store{}).Where("id = ?", store.ID).Update("store_code", code).Error; err != nil {
			logging.LogDatabaseError("更新门店编码失败", err, zap.Uint("store_id", store.ID))
			return err
		}
		logging.LogInfo("门店编码已补全", zap.Uint("store_id", store.ID))
	}

	return nil
}
