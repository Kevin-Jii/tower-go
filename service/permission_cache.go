package service

import (
	"fmt"

	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/module"
	"github.com/Kevin-Jii/tower-go/utils/cache"
	"github.com/Kevin-Jii/tower-go/utils/database"
)

const (
	userRoleCacheKeyFormat = "tower:user:role:%d"
	roleMenuCacheKeyFormat = "tower:role:menu:%d"
)

// userPermCacheKey 须包含门店与角色：仅 userId 会导致换店/改角色后仍命中旧权限缓存，接口 403 而菜单仍可见。
func userPermCacheKey(userID, storeID, roleID uint) string {
	// v2：与门店菜单闭包扩展对齐，避免沿用旧权限列表
	return fmt.Sprintf("tower:user:perm:v2:%d:%d:%d", userID, storeID, roleID)
}

type permissionBuildUser struct {
	ID       uint
	StoreID  uint
	RoleID   uint
	RoleCode string
}

func BuildUserPermissionCache(userID uint, storeID uint, roleID uint, roleCode string) ([]string, error) {
	if userID == 0 || roleID == 0 {
		return []string{}, nil
	}

	menuModule := NewMenuService(
		module.NewMenuModule(database.DB),
		module.NewRoleMenuModule(database.DB),
		module.NewStoreRoleMenuModule(database.DB),
	)
	var perms []string
	var err error

	if roleCode == model.RoleCodeAdmin || roleCode == model.RoleCodeSuperAdmin {
		perms, err = menuModule.GetAllPermissions()
	} else {
		perms, err = menuModule.GetUserPermissions(storeID, roleID)
	}
	if err != nil {
		return nil, err
	}

	_ = cache.CacheSet(fmt.Sprintf(userRoleCacheKeyFormat, userID), []uint{roleID}, cache.PermissionsTTL)
	_ = cache.CacheSet(fmt.Sprintf(roleMenuCacheKeyFormat, roleID), perms, cache.PermissionsTTL)
	_ = cache.CacheSet(userPermCacheKey(userID, storeID, roleID), perms, cache.PermissionsTTL)

	return perms, nil
}

func GetUserPermissionCodes(userID uint, storeID uint, roleID uint, roleCode string) ([]string, error) {
	if userID == 0 {
		return []string{}, nil
	}

	var perms []string
	err := cache.CacheGet(userPermCacheKey(userID, storeID, roleID), &perms)
	if err == nil {
		return perms, nil
	}

	return BuildUserPermissionCache(userID, storeID, roleID, roleCode)
}

func InvalidateUserPermissionCache(userID uint) {
	if userID == 0 {
		return
	}
	_ = cache.CacheDelete(fmt.Sprintf(userRoleCacheKeyFormat, userID))
	_ = cache.CacheDeleteByPattern(fmt.Sprintf("tower:user:perm:%d:*", userID))
	_ = cache.CacheDeleteByPattern(fmt.Sprintf("tower:user:perm:v2:%d:*", userID))
}

func InvalidateRolePermissionCache(roleID uint) {
	if roleID == 0 || database.DB == nil {
		return
	}

	_ = cache.CacheDelete(fmt.Sprintf(roleMenuCacheKeyFormat, roleID))

	var users []permissionBuildUser
	if err := database.DB.Table("users").
		Select("users.id, users.store_id, users.role_id, roles.code as role_code").
		Joins("left join roles on roles.id = users.role_id").
		Where("users.role_id = ?", roleID).
		Find(&users).Error; err != nil {
		return
	}

	for _, user := range users {
		InvalidateUserPermissionCache(user.ID)
	}
}

func InvalidateAllPermissionCache() {
	_ = cache.CacheDeleteByPattern("tower:user:perm:*")
	_ = cache.CacheDeleteByPattern("tower:user:role:*")
	_ = cache.CacheDeleteByPattern("tower:role:menu:*")
}
