package controller

import (
	"net/http"
	"strings"
	"time"
	"tower-go/utils"

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
	auth := c.GetHeader("Authorization")
	if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
		return
	}
	token := strings.TrimPrefix(auth, "Bearer ")
	claims, err := utils.ValidateToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
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

	sm := utils.GetSessionManager()
	if sm == nil { // 未初始化
		conn.WriteJSON(utils.WebSocketMessage{Type: "error", Payload: gin.H{"message": "session manager not ready"}, Ts: time.Now().Unix()})
		conn.Close()
		return
	}

	// 创建 session，可能返回被挤下线的旧会话
	session, kicked := sm.CreateSession(claims.UserID, deviceID, token, claims.ExpiresAt.Time, conn)
	for _, ks := range kicked {
		_ = ks.Conn.WriteJSON(utils.WebSocketMessage{Type: "kick", Payload: gin.H{"reason": "replaced"}, Ts: time.Now().Unix()})
		ks.Conn.Close()
	}
	conn.WriteJSON(utils.WebSocketMessage{Type: "connected", Payload: gin.H{"session_id": session.ID}, Ts: time.Now().Unix()})

	// 读取循环，支持 ping/pong
	for {
		var incoming map[string]interface{}
		if err := conn.ReadJSON(&incoming); err != nil {
			break
		}
		msgType, _ := incoming["type"].(string)
		switch msgType {
		case "ping":
			conn.WriteJSON(utils.WebSocketMessage{Type: "pong", Ts: time.Now().Unix()})
		default:
			conn.WriteJSON(utils.WebSocketMessage{Type: "echo", Payload: incoming, Ts: time.Now().Unix()})
		}
	}
	sm.RemoveSession(session.ID)
}
