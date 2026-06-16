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
	return fmt.Sprintf("tower:user:perm:v3:%d:%d:%d", userID, storeID, roleID)
}

// SyncUserPermissionCache 将已计算好的权限码写入 Redis，与 Permission 中间件使用的键一致。
// user-permissions 接口若不走本函数，而中间件仍读陈旧缓存，会出现「前端权限列表含 supplier:list、GET /suppliers 仍 403」。
func SyncUserPermissionCache(userID, storeID, roleID uint, perms []string) {
	if userID == 0 || roleID == 0 {
		return
	}
	_ = cache.CacheSet(userPermCacheKey(userID, storeID, roleID), perms, cache.PermissionsTTL)
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

	// 与 HQUnboundAdmin 一致：仅 Token 未绑店的总部账号用全量权限码；绑店 admin/super 走门店角色菜单
	if model.HQUnboundAdminRole(roleCode, storeID) {
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
	// Redis 未启用时 CacheGet 恒为 nil 且不填充 dest，若此处直接 return 会得到空切片，门店账号所有 Permission 路由都会 403。
	// 与 module/menu 一致：仅当缓存命中且非空时才短路；否则走库重建。
	if err == nil && len(perms) > 0 {
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
	_ = cache.CacheDeleteByPattern(fmt.Sprintf("tower:user:perm:v3:%d:*", userID))
}

func InvalidateRolePermissionCache(roleID uint) {
	if roleID == 0 || database.DB == nil {
		return
	}

	_ = cache.CacheDelete(fmt.Sprintf(roleMenuCacheKeyFormat, roleID))
	_ = cache.CacheDeleteByPattern(fmt.Sprintf("tower:store:role:menus:*:%d", roleID))

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

func InvalidateStoreRolePermissionCache(storeID uint, roleID uint) {
	if storeID == 0 || roleID == 0 || database.DB == nil {
		return
	}

	_ = cache.CacheDelete(fmt.Sprintf("%s:%d:%d", cache.CacheKeyStoreRoleMenus, storeID, roleID))

	var users []permissionBuildUser
	if err := database.DB.Table("users").
		Select("users.id, users.store_id, users.role_id, roles.code as role_code").
		Joins("left join roles on roles.id = users.role_id").
		Where("users.store_id = ? AND users.role_id = ?", storeID, roleID).
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
	_ = cache.CacheDeleteByPattern("tower:store:role:menus:*")
}
