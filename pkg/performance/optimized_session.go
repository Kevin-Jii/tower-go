package performance

import (
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// OptimizedSession represents a user session with WebSocket connection
type OptimizedSession struct {
	ID        string
	UserID    uint
	DeviceID  string
	Token     string
	ExpiresAt time.Time
	Conn      *websocket.Conn
	CreatedAt time.Time
}

// OptimizedSessionManager manages WebSocket sessions using sync.Map for better concurrent performance
// sync.Map is optimized for read-heavy workloads which is typical for session lookups
type OptimizedSessionManager struct {
	// userID -> map[sessionID]*OptimizedSession
	userSessions sync.Map
	// sessionID -> userID (reverse index for fast deletion)
	sessionIndex sync.Map
	// Strategy: "single" (single sign-on) or "multi" (allow multiple devices, limit max)
	strategy string
	// Maximum sessions for multi strategy
	maxSessions int
	// Mutex for operations that need to be atomic across multiple maps
	mu sync.Mutex
}

// NewOptimizedSessionManager creates a new OptimizedSessionManager instance
func NewOptimizedSessionManager(strategy string, maxSessions int) *OptimizedSessionManager {
	if maxSessions <= 0 {
		maxSessions = 3
	}
	return &OptimizedSessionManager{
		strategy:    strategy,
		maxSessions: maxSessions,
	}
}

// CreateSession creates and registers a new session (caller has already authenticated JWT)
func (osm *OptimizedSessionManager) CreateSession(userID uint, deviceID, token string, expiresAt time.Time, conn *websocket.Conn) (*OptimizedSession, []*OptimizedSession) {
	osm.mu.Lock()
	defer osm.mu.Unlock()

	// Get or create user sessions map
	var userSessionsMap *sync.Map
	if val, ok := osm.userSessions.Load(userID); ok {
		userSessionsMap = val.(*sync.Map)
	} else {
		userSessionsMap = &sync.Map{}
		osm.userSessions.Store(userID, userSessionsMap)
	}

	// Handle old sessions based on strategy
	kicked := []*OptimizedSession{}
	if osm.strategy == "single" {
		// Kick all existing sessions
		userSessionsMap.Range(func(key, value interface{}) bool {
			session := value.(*OptimizedSession)
			kicked = append(kicked, session)
			return true
		})
		// Clear all sessions
		userSessionsMap.Range(func(key, value interface{}) bool {
			userSessionsMap.Delete(key)
			osm.sessionIndex.Delete(key)
			return true
		})
	} else if osm.strategy == "multi" {
		// Count existing sessions
		count := 0
		userSessionsMap.Range(func(key, value interface{}) bool {
			count++
			return true
		})

		if count >= osm.maxSessions {
			// Find and kick the oldest session
			var oldestID string
			var oldestTime time.Time
			first := true
			userSessionsMap.Range(func(key, value interface{}) bool {
				session := value.(*OptimizedSession)
				if first || session.CreatedAt.Before(oldestTime) {
					oldestID = session.ID
					oldestTime = session.CreatedAt
					first = false
				}
				return true
			})

			if oldestID != "" {
				if val, ok := userSessionsMap.Load(oldestID); ok {
					kicked = append(kicked, val.(*OptimizedSession))
					userSessionsMap.Delete(oldestID)
					osm.sessionIndex.Delete(oldestID)
				}
			}
		}
	}

	// Create new session
	session := &OptimizedSession{
		ID:        uuid.NewString(),
		UserID:    userID,
		DeviceID:  deviceID,
		Token:     token,
		ExpiresAt: expiresAt,
		Conn:      conn,
		CreatedAt: time.Now(),
	}

	userSessionsMap.Store(session.ID, session)
	osm.sessionIndex.Store(session.ID, userID)

	return session, kicked
}

// RemoveSession removes a session
func (osm *OptimizedSessionManager) RemoveSession(sessionID string) {
	userIDVal, ok := osm.sessionIndex.Load(sessionID)
	if !ok {
		return
	}

	userID := userIDVal.(uint)
	osm.sessionIndex.Delete(sessionID)

	if val, ok := osm.userSessions.Load(userID); ok {
		userSessionsMap := val.(*sync.Map)
		userSessionsMap.Delete(sessionID)

		// Check if user has no more sessions
		hasSession := false
		userSessionsMap.Range(func(key, value interface{}) bool {
			hasSession = true
			return false // Stop iteration
		})

		if !hasSession {
			osm.userSessions.Delete(userID)
		}
	}
}

// KickUser kicks all sessions for a user
func (osm *OptimizedSessionManager) KickUser(userID uint, reason string) int {
	val, ok := osm.userSessions.Load(userID)
	if !ok {
		return 0
	}

	userSessionsMap := val.(*sync.Map)
	sessions := []*OptimizedSession{}

	// Collect all sessions
	userSessionsMap.Range(func(key, value interface{}) bool {
		session := value.(*OptimizedSession)
		sessions = append(sessions, session)
		return true
	})

	if len(sessions) == 0 {
		return 0
	}

	// Remove from maps
	osm.mu.Lock()
	osm.userSessions.Delete(userID)
	for _, s := range sessions {
		osm.sessionIndex.Delete(s.ID)
	}
	osm.mu.Unlock()

	// Close connections
	for _, s := range sessions {
		_ = s.Conn.WriteJSON(map[string]interface{}{
			"type":    "kick",
			"payload": map[string]string{"reason": reason},
			"ts":      time.Now().Unix(),
		})
		s.Conn.Close()
	}

	return len(sessions)
}

// KickSession kicks a single session
func (osm *OptimizedSessionManager) KickSession(sessionID string, reason string) bool {
	userIDVal, ok := osm.sessionIndex.Load(sessionID)
	if !ok {
		return false
	}

	userID := userIDVal.(uint)
	val, ok := osm.userSessions.Load(userID)
	if !ok {
		return false
	}

	userSessionsMap := val.(*sync.Map)
	sessionVal, ok := userSessionsMap.Load(sessionID)
	if !ok {
		return false
	}

	session := sessionVal.(*OptimizedSession)

	// Remove from maps
	osm.mu.Lock()
	userSessionsMap.Delete(sessionID)
	osm.sessionIndex.Delete(sessionID)

	// Check if user has no more sessions
	hasSession := false
	userSessionsMap.Range(func(key, value interface{}) bool {
		hasSession = true
		return false
	})
	if !hasSession {
		osm.userSessions.Delete(userID)
	}
	osm.mu.Unlock()

	// Close connection
	_ = session.Conn.WriteJSON(map[string]interface{}{
		"type":    "kick",
		"payload": map[string]string{"reason": reason},
		"ts":      time.Now().Unix(),
	})
	session.Conn.Close()

	return true
}

// ListUserSessions lists all sessions for a user
func (osm *OptimizedSessionManager) ListUserSessions(userID uint) []*OptimizedSession {
	val, ok := osm.userSessions.Load(userID)
	if !ok {
		return []*OptimizedSession{}
	}

	userSessionsMap := val.(*sync.Map)
	result := []*OptimizedSession{}

	userSessionsMap.Range(func(key, value interface{}) bool {
		session := value.(*OptimizedSession)
		result = append(result, session)
		return true
	})

	return result
}

// Broadcast broadcasts a message to all sessions of a user
func (osm *OptimizedSessionManager) Broadcast(userID uint, msg interface{}) int {
	val, ok := osm.userSessions.Load(userID)
	if !ok {
		return 0
	}

	userSessionsMap := val.(*sync.Map)
	count := 0

	userSessionsMap.Range(func(key, value interface{}) bool {
		session := value.(*OptimizedSession)
		if err := session.Conn.WriteJSON(msg); err == nil {
			count++
		}
		return true
	})

	return count
}

// ValidateStrategy validates configuration
func (osm *OptimizedSessionManager) ValidateStrategy() error {
	if osm.strategy != "single" && osm.strategy != "multi" {
		return errors.New("invalid session strategy")
	}
	if osm.strategy == "multi" && osm.maxSessions <= 0 {
		return errors.New("maxSessions must be > 0 for multi strategy")
	}
	return nil
}

// GetSessionCount returns the total number of sessions
func (osm *OptimizedSessionManager) GetSessionCount() int {
	count := 0
	osm.userSessions.Range(func(key, value interface{}) bool {
		userSessionsMap := value.(*sync.Map)
		userSessionsMap.Range(func(k, v interface{}) bool {
			count++
			return true
		})
		return true
	})
	return count
}

// GetUserCount returns the number of users with active sessions
func (osm *OptimizedSessionManager) GetUserCount() int {
	count := 0
	osm.userSessions.Range(func(key, value interface{}) bool {
		count++
		return true
	})
	return count
}
