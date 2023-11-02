package controller

import (
	"Backend/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Chat(db *gorm.DB, q *gin.Engine) {
	r := q.Group("/api/v1/chat")
	r.GET("/chat", middleware.Authorization(), func(c *gin.Context) {
		// mock up chat

	})

}
