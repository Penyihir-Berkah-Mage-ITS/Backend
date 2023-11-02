package websocket

import "github.com/gin-gonic/gin"

func WebSocketInit(q *gin.Engine, wsHandler *Handler) {
	r := q.Group("/api/v1/ws")
	r.POST("/ws/createRoom", wsHandler.CreateRoom)
	r.GET("/ws/joinRoom/:roomId", wsHandler.JoinRoom)
	r.GET("/ws/getRooms", wsHandler.GetRooms)
	r.GET("/ws/getClients/:roomId", wsHandler.GetClients)
}
