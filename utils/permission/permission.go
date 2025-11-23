package permission

import "github.com/Kevin-Jii/tower-go/model"

// CheckMenuPermission 检查用户对菜单是否有指定权限
func CheckMenuPermission(menuPerms uint8, requiredPerm uint8) bool {
	return model.HasPermission(menuPerms, requiredPerm)
}

// HasViewPermission 检查是否有查看权限
func HasViewPermission(perms uint8) bool {
	return model.HasPermission(perms, model.PermView)
}

// HasCreatePermission 检查是否有新增权限
func HasCreatePermission(perms uint8) bool {
	return model.HasPermission(perms, model.PermCreate)
}

// HasUpdatePermission 检查是否有修改权限
func HasUpdatePermission(perms uint8) bool {
	return model.HasPermission(perms, model.PermUpdate)
}

// HasDeletePermission 检查是否有删除权限
func HasDeletePermission(perms uint8) bool {
	return model.HasPermission(perms, model.PermDelete)
}

// ParsePermissionString 将权限字符串转换为权限位
// 例如: "1111" -> 15, "1000" -> 8, "1110" -> 14
func ParsePermissionString(permStr string) uint8 {
	var perm uint8 = 0
	if len(permStr) >= 4 {
		if permStr[0] == '1' {
			perm |= model.PermView
		}
		if permStr[1] == '1' {
			perm |= model.PermCreate
		}
		if permStr[2] == '1' {
			perm |= model.PermUpdate
		}
		if permStr[3] == '1' {
			perm |= model.PermDelete
		}
	}
	return perm
}

// FormatPermissionBits 将权限位转换为二进制字符串
// 例如: 15 -> "1111", 8 -> "1000", 14 -> "1110"
func FormatPermissionBits(perm uint8) string {
	result := ""
	if perm&model.PermView != 0 {
		result += "1"
	} else {
		result += "0"
	}
	if perm&model.PermCreate != 0 {
		result += "1"
	} else {
		result += "0"
	}
	if perm&model.PermUpdate != 0 {
		result += "1"
	} else {
		result += "0"
	}
	if perm&model.PermDelete != 0 {
		result += "1"
	} else {
		result += "0"
	}
	return result
}

// GetPermissionDescription 获取权限描述
func GetPermissionDescription(perm uint8) map[string]bool {
	return map[string]bool{
		"view":   HasViewPermission(perm),
		"create": HasCreatePermission(perm),
		"update": HasUpdatePermission(perm),
		"delete": HasDeletePermission(perm),
	}
}
