package controller

import (
	"Backend/middleware"
	"Backend/model"
	"Backend/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func Home(db *gorm.DB, q *gin.Engine) {
	r := q.Group("/api/v1/home")
	r.GET("/", middleware.Authorization(), func(c *gin.Context) {
		var posts []model.Post
		if err := db.Find(&posts).Error; err != nil {
			utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
			return
		}

		utils.HttpRespSuccess(c, http.StatusOK, "Success get all posts", posts)
	})

	r.GET("/popular", middleware.Authorization(), func(c *gin.Context) {
		var posts []model.Post
		if err := db.Order("\"like\" DESC").Find(&posts).Error; err != nil {
			utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
			return
		}

		utils.HttpRespSuccess(c, http.StatusOK, "Success get all posts", posts)
	})

	r.GET("/latest", middleware.Authorization(), func(c *gin.Context) {
		var posts []model.Post
		if err := db.Order("created_at desc").Find(&posts).Error; err != nil {
			utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
			return
		}

		utils.HttpRespSuccess(c, http.StatusOK, "Success get all posts", posts)
	})
}
