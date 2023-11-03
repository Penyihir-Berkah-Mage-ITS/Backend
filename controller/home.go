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
		if err := db.Preload("User").Order("created_at desc").Find(&posts).Error; err != nil {
			utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
			return
		}

		var postsResponse []model.PostResponse

		for _, post := range posts {
			distance := utils.LocationToKM(
				c,
				latitudeStr,
				longitudeStr,
				strconv.FormatFloat(post.Latitude, 'f', -1, 64),
				strconv.FormatFloat(post.Longitude, 'f', -1, 64),
			)

			postResponse := model.PostResponse{
				ID:         post.ID,
				User:       post.User,
				UserID:     post.UserID,
				Content:    post.Content,
				Attachment: post.Attachment,
				Likes:      post.Likes,
				Latitude:   post.Latitude,
				Longitude:  post.Longitude,
				Distance:   distance,
				CreatedAt:  post.CreatedAt,
				UpdatedAt:  post.UpdatedAt,
			}

			id, _ := c.Get("id")
			var like model.UserLikePost
			if err := db.Where("user_id = ? AND post_id = ?", id, post.ID).First(&like).Error; err != nil {
				postResponse.IsLiked = false
			} else {
				postResponse.IsLiked = true
			}

			var totalComment int64
			if err := db.Table("comments").Where("post_id = ?", post.ID).Count(&totalComment).Error; err != nil {
				utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
				return
			}

			postResponse.TotalComment = totalComment

			postsResponse = append(postsResponse, postResponse)
		}

		utils.HttpRespSuccess(c, http.StatusOK, "Success get all posts", postsResponse)
	})

	r.GET("/popular", middleware.Authorization(), func(c *gin.Context) {
		var posts []model.Post

		if err := db.Preload("User").Order("likes desc").Find(&posts).Error; err != nil {
			utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
			return
		}

		var postsResponse []model.PostResponse

		for _, post := range posts {
			postResponse := model.PostResponse{
				ID:         post.ID,
				User:       post.User,
				UserID:     post.UserID,
				Content:    post.Content,
				Attachment: post.Attachment,
				Likes:      post.Likes,
				Latitude:   post.Latitude,
				Longitude:  post.Longitude,
				Distance:   post.Distance,
				CreatedAt:  post.CreatedAt,
				UpdatedAt:  post.UpdatedAt,
			}

			id, _ := c.Get("id")
			var like model.UserLikePost
			if err := db.Where("user_id = ? AND post_id = ?", id, post.ID).First(&like).Error; err != nil {
				postResponse.IsLiked = false
			} else {
				postResponse.IsLiked = true
			}

			var totalComment int64
			if err := db.Table("comments").Where("post_id = ?", post.ID).Count(&totalComment).Error; err != nil {
				utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
				return
			}

			postResponse.TotalComment = totalComment

			postsResponse = append(postsResponse, postResponse)
		}

		utils.HttpRespSuccess(c, http.StatusOK, "Success get all posts", postsResponse)
	})

	r.GET("/latest", middleware.Authorization(), func(c *gin.Context) {
		var posts []model.Post

		if err := db.Order("created_at desc").Preload("User").Find(&posts).Error; err != nil {
			utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
			return
		}

		var postsResponse []model.PostResponse

		for _, post := range posts {
			postResponse := model.PostResponse{
				ID:         post.ID,
				User:       post.User,
				UserID:     post.UserID,
				Content:    post.Content,
				Attachment: post.Attachment,
				Likes:      post.Likes,
				Latitude:   post.Latitude,
				Longitude:  post.Longitude,
				Distance:   post.Distance,
				CreatedAt:  post.CreatedAt,
				UpdatedAt:  post.UpdatedAt,
			}

			id, _ := c.Get("id")
			var like model.UserLikePost
			if err := db.Where("user_id = ? AND post_id = ?", id, post.ID).First(&like).Error; err != nil {
				postResponse.IsLiked = false
			} else {
				postResponse.IsLiked = true
			}

			var totalComment int64
			if err := db.Table("comments").Where("post_id = ?", post.ID).Count(&totalComment).Error; err != nil {
				utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
				return
			}

			postResponse.TotalComment = totalComment

			postsResponse = append(postsResponse, postResponse)
		}

		utils.HttpRespSuccess(c, http.StatusOK, "Success get all posts", postsResponse)
	})
}
