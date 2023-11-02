package controller

import (
	"Backend/middleware"
	"Backend/model"
	"Backend/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func Home(db *gorm.DB, q *gin.Engine) {
	r := q.Group("/api/v1/home")
	r.GET("/nearest", middleware.Authorization(), func(c *gin.Context) {
		latitudeStr := c.Query("lat")
		longitudeStr := c.Query("lng")

		var posts []model.Post
		if err := db.Order("created_at desc").Find(&posts).Error; err != nil {
			utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
			return
		}

		for i := range posts {
			distance := utils.LocationToKM(
				c,
				latitudeStr,
				longitudeStr,
				strconv.FormatFloat(posts[i].Latitude, 'f', -1, 64),
				strconv.FormatFloat(posts[i].Longitude, 'f', -1, 64),
			)
			posts[i].Distance = distance
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
