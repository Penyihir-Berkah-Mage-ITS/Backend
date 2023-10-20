package controller

import (
	"Backend/middleware"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
	"log"
	"net/http"
)

var (
	// WebSocket upgrader
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // In a real application, you would validate the request origin for security.
		},
	}

	// Secret key for JWT (replace with your secret)
	secretKey = []byte("your-secret-key")
)

func Chat(db *gorm.DB, q *gin.Engine) {
	r := q.Group("/api/v1/chat")

	r.GET("/api/v1/chat", middleware.Authorization(), func(c *gin.Context) {
		// Check and validate JWT token
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Println(err)
			return
		}

		defer conn.Close()

		// Handle WebSocket communication here
		for {
			// Read message
			messageType, p, err := conn.ReadMessage()
			if err != nil {
				return
			}

			message := string(p)

			fmt.Printf("Message type: %d\n", messageType)

			// Handle the received message, e.g., broadcast to other clients
			fmt.Printf("Received: %s\n", message)
		}
	})
}
