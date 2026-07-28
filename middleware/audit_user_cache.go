package middleware

import (
	"sync"
	"time"

	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/utils/database"
)

const (
	auditUserCacheTTL   = 5 * time.Minute
	auditUserCacheLimit = 2048
)

type auditUserSnapshot struct {
	username  string
	nickname  string
	phone     string
	roleName  string
	roleCode  string
	storeID   uint
	storeName string
	expiresAt time.Time
}

var auditUserSnapshots = struct {
	sync.RWMutex
	values map[uint]auditUserSnapshot
}{values: make(map[uint]auditUserSnapshot)}

func getAuditUserSnapshot(userID uint) (auditUserSnapshot, bool) {
	if userID == 0 || database.DB == nil {
		return auditUserSnapshot{}, false
	}
	now := time.Now()
	auditUserSnapshots.RLock()
	cached, ok := auditUserSnapshots.values[userID]
	auditUserSnapshots.RUnlock()
	if ok && now.Before(cached.expiresAt) {
		return cached, true
	}

	var user model.User
	if err := database.DB.Preload("Role").Preload("Store").First(&user, userID).Error; err != nil {
		return auditUserSnapshot{}, false
	}
	snapshot := auditUserSnapshot{
		username:  user.Username,
		nickname:  user.Nickname,
		phone:     user.Phone,
		storeID:   user.StoreID,
		expiresAt: now.Add(auditUserCacheTTL),
	}
	if user.Role != nil {
		snapshot.roleName = user.Role.Name
		snapshot.roleCode = user.Role.Code
	}
	if user.Store != nil {
		snapshot.storeName = user.Store.Name
	}

	auditUserSnapshots.Lock()
	if len(auditUserSnapshots.values) >= auditUserCacheLimit {
		for id, item := range auditUserSnapshots.values {
			if !now.Before(item.expiresAt) {
				delete(auditUserSnapshots.values, id)
			}
		}
		if len(auditUserSnapshots.values) >= auditUserCacheLimit {
			for id := range auditUserSnapshots.values {
				delete(auditUserSnapshots.values, id)
				break
			}
		}
	}
	auditUserSnapshots.values[userID] = snapshot
	auditUserSnapshots.Unlock()
	return snapshot, true
}
