package controller

import (
	"net/http"
	"strings"
	"time"
	authPkg "github.com/Kevin-Jii/tower-go/utils/auth"
	"github.com/Kevin-Jii/tower-go/utils/session"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// WebSocketHandler 处理用户 WebSocket 连接
// 客户端应在 Header: Authorization: Bearer <token>
// 可选查询参数 device_id 指定设备ID
func WebSocketHandler(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		c.JSON(401, gin.H{"error": "missing bearer token"})
		return
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := authPkg.ValidateToken(token)
	if err != nil {
		c.JSON(401, gin.H{"error": "invalid token"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	deviceID := c.Query("device_id")
	if deviceID == "" {
		deviceID = "unknown"
	}

	sm := session.GetSessionManager()
	if sm == nil { // 未初始化
		msg := session.WebSocketMessage{
			Type:    "error",
			Payload: gin.H{"message": "session manager not ready"},
			Ts:      time.Now().Unix(),
		}
		conn.WriteJSON(msg)
		conn.Close()
		return
	}

	// 创建 session，可能返回被挤下线的旧会话
	newSession, kicked := sm.CreateSession(claims.UserID, deviceID, token, claims.ExpiresAt.Time, conn)
	for _, ks := range kicked {
		kickMsg := session.WebSocketMessage{
			Type:    "kick",
			Payload: gin.H{"reason": "replaced"},
			Ts:      time.Now().Unix(),
		}
		_ = ks.Conn.WriteJSON(kickMsg)
		ks.Conn.Close()
	}
	connMsg := session.WebSocketMessage{
		Type:    "connected",
		Payload: gin.H{"session_id": newSession.ID},
		Ts:      time.Now().Unix(),
	}
	conn.WriteJSON(connMsg)

	// 读取循环，支持 ping/pong
	for {
		var incoming map[string]interface{}
		if err := conn.ReadJSON(&incoming); err != nil {
			break
		}
		msgType, _ := incoming["type"].(string)
		switch msgType {
		case "ping":
			pongMsg := session.WebSocketMessage{
				Type: "pong",
				Ts:   time.Now().Unix(),
			}
			conn.WriteJSON(pongMsg)
		default:
			echoMsg := session.WebSocketMessage{
				Type:    "echo",
				Payload: incoming,
				Ts:      time.Now().Unix(),
			}
			conn.WriteJSON(echoMsg)
		}
	}
	sm.RemoveSession(newSession.ID)
}
