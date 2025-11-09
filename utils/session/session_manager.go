package session

import (
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// WebSocketMessage 基础消息协议
type WebSocketMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload,omitempty"`
	Ts      int64       `json:"ts"` // unix 秒
}

// Session 表示一个用户在某终端的一次在线会话
type Session struct {
	ID        string
	UserID    uint
	DeviceID  string
	Token     string
	ExpiresAt time.Time
	Conn      *websocket.Conn
	CreatedAt time.Time
}

// SessionManager 管理所有 WebSocket 会话
type SessionManager struct {
	mu sync.RWMutex
	// userID -> sessionID -> *Session
	userSessions map[uint]map[string]*Session
	// sessionID -> userID (反向索引方便快速删除)
	sessionIndex map[string]uint
	// 策略：single(单点登录) 或 multi(允许多端, 限制最大数)
	strategy string
	// multi 策略下最大会话数
	maxSessions int
}

var globalSessionManager *SessionManager

// InitSessionManager 初始化全局会话管理器
func InitSessionManager(strategy string, maxSessions int) {
	if maxSessions <= 0 {
		maxSessions = 3
	}
	globalSessionManager = &SessionManager{
		userSessions: make(map[uint]map[string]*Session),
		sessionIndex: make(map[string]uint),
		strategy:     strategy,
		maxSessions:  maxSessions,
	}
}

func GetSessionManager() *SessionManager { return globalSessionManager }

// CreateSession 创建并注册一个新会话（调用方已认证 JWT）
func (sm *SessionManager) CreateSession(userID uint, deviceID, token string, expiresAt time.Time, conn *websocket.Conn) (*Session, []*Session) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if sm.userSessions[userID] == nil {
		sm.userSessions[userID] = make(map[string]*Session)
	}

	// 根据策略处理旧会话
	kicked := []*Session{}
	if sm.strategy == "single" {
		for _, s := range sm.userSessions[userID] {
			kicked = append(kicked, s)
		}
		sm.userSessions[userID] = make(map[string]*Session)
	} else if sm.strategy == "multi" && len(sm.userSessions[userID]) >= sm.maxSessions {
		// 淘汰最早的
		var oldestID string
		var oldestTime time.Time
		first := true
		for id, s := range sm.userSessions[userID] {
			if first || s.CreatedAt.Before(oldestTime) {
				oldestID = id
				oldestTime = s.CreatedAt
				first = false
			}
		}
		if oldestID != "" {
			kicked = append(kicked, sm.userSessions[userID][oldestID])
			delete(sm.sessionIndex, oldestID)
			delete(sm.userSessions[userID], oldestID)
		}
	}

	session := &Session{
		ID:        uuid.NewString(),
		UserID:    userID,
		DeviceID:  deviceID,
		Token:     token,
		ExpiresAt: expiresAt,
		Conn:      conn,
		CreatedAt: time.Now(),
	}
	sm.userSessions[userID][session.ID] = session
	sm.sessionIndex[session.ID] = userID
	return session, kicked
}

// RemoveSession 删除会话
func (sm *SessionManager) RemoveSession(sessionID string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	userID, ok := sm.sessionIndex[sessionID]
	if !ok {
		return
	}
	delete(sm.sessionIndex, sessionID)
	if sessions, ok2 := sm.userSessions[userID]; ok2 {
		delete(sessions, sessionID)
		if len(sessions) == 0 {
			delete(sm.userSessions, userID)
		}
	}
}

// KickUser 踢出某用户全部会话
func (sm *SessionManager) KickUser(userID uint, reason string) int {
	sm.mu.Lock()
	sessionsMap := sm.userSessions[userID]
	if len(sessionsMap) == 0 {
		sm.mu.Unlock()
		return 0
	}
	sessions := []*Session{}
	for _, s := range sessionsMap {
		sessions = append(sessions, s)
	}
	delete(sm.userSessions, userID)
	for _, s := range sessions {
		delete(sm.sessionIndex, s.ID)
	}
	sm.mu.Unlock()

	for _, s := range sessions {
		_ = s.Conn.WriteJSON(WebSocketMessage{Type: "kick", Payload: map[string]string{"reason": reason}, Ts: time.Now().Unix()})
		s.Conn.Close()
	}
	return len(sessions)
}

// KickSession 踢出单个会话
func (sm *SessionManager) KickSession(sessionID string, reason string) bool {
	sm.mu.Lock()
	userID, ok := sm.sessionIndex[sessionID]
	if !ok {
		sm.mu.Unlock()
		return false
	}
	session := sm.userSessions[userID][sessionID]
	delete(sm.userSessions[userID], sessionID)
	delete(sm.sessionIndex, sessionID)
	if len(sm.userSessions[userID]) == 0 {
		delete(sm.userSessions, userID)
	}
	sm.mu.Unlock()
	_ = session.Conn.WriteJSON(WebSocketMessage{Type: "kick", Payload: map[string]string{"reason": reason}, Ts: time.Now().Unix()})
	session.Conn.Close()
	return true
}

// ListUserSessions 列出用户会话ID
func (sm *SessionManager) ListUserSessions(userID uint) []*Session {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	result := []*Session{}
	for _, s := range sm.userSessions[userID] {
		result = append(result, s)
	}
	return result
}

// Broadcast 向某用户广播
func (sm *SessionManager) Broadcast(userID uint, msg WebSocketMessage) int {
	sm.mu.RLock()
	sessions := sm.userSessions[userID]
	sm.mu.RUnlock()
	count := 0
	for _, s := range sessions {
		if err := s.Conn.WriteJSON(msg); err == nil {
			count++
		}
	}
	return count
}

// ValidateStrategy 配置合理性检查
func (sm *SessionManager) ValidateStrategy() error {
	if sm.strategy != "single" && sm.strategy != "multi" {
		return errors.New("invalid session strategy")
	}
	if sm.strategy == "multi" && sm.maxSessions <= 0 {
		return errors.New("maxSessions must be > 0 for multi strategy")
	}
	return nil
}
