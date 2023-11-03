package controller

import (
	"Backend/middleware"
	"Backend/model"
	"Backend/utils"
	supabasestorageuploader "github.com/adityarizkyramadhan/supabase-storage-uploader"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"net/http"
	"os"
	"strings"
	"time"
)

func Post(db *gorm.DB, q *gin.Engine) {
	supClient := supabasestorageuploader.New(
		os.Getenv("SUPABASE_PROJECT_URL"),
		os.Getenv("SUPABASE_PROJECT_API_KEY"),
		os.Getenv("SUPABASE_PROJECT_STORAGE_NAME"),
	)

	r := q.Group("/api/v1/post")
	r.GET("/:post_id", middleware.Authorization(), func(c *gin.Context) {
		postID := c.Param("post_id")

		var post model.Post
		if err := db.Where("id = ?", postID).First(&post).Error; err != nil {
			utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
			return
		}

		utils.HttpRespSuccess(c, http.StatusOK, "Success get post", post)
	})

	r.POST("/create", middleware.Authorization(), func(c *gin.Context) {
		id, _ := c.Get("id")
		latitudeStr := c.Query("lat")
		longitudeStr := c.Query("lng")
		latitude, err := utils.StringToFloat64(latitudeStr)
		if err != nil {
			utils.HttpRespFailed(c, http.StatusBadRequest, err.Error())
			return
		}
		longitude, err := utils.StringToFloat64(longitudeStr)
		if err != nil {
			utils.HttpRespFailed(c, http.StatusBadRequest, err.Error())
			return
		}

		randomID := utils.GenerateStringID()
		content := c.PostForm("content")

		attachment, _ := c.FormFile("attachment")

		if attachment == nil {
			newPost := model.Post{
				ID:         utils.GenerateStringID(),
				UserID:     id.(uuid.UUID),
				Content:    content,
				Attachment: "",
				Likes:      0,
				Latitude:   latitude,
				Longitude:  longitude,
				CreatedAt:  time.Now(),
			}

			if err := db.Create(&newPost).Error; err != nil {
				utils.HttpRespFailed(c, http.StatusInternalServerError, err.Error())
				return
			}

			utils.HttpRespSuccess(c, http.StatusOK, "Success create post", newPost)
			return
		}

		filename := strings.ReplaceAll(strings.TrimSpace(attachment.Filename), " ", "")
		newFilename := randomID + "_" + filename
		attachment.Filename = newFilename

		photoLink, err := supClient.Upload(attachment)
		if err != nil {
			utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
			return
		}

		newPost := model.Post{
			ID:         utils.GenerateStringID(),
			UserID:     id.(uuid.UUID),
			Content:    content,
			Attachment: photoLink,
			Likes:      0,
			Latitude:   latitude,
			Longitude:  longitude,
			CreatedAt:  time.Now(),
		}

		if err := db.Create(&newPost).Error; err != nil {
			utils.HttpRespFailed(c, http.StatusInternalServerError, err.Error())
			return
		}

		utils.HttpRespSuccess(c, http.StatusOK, "Success create post", newPost)
	})

	r.DELETE("/delete/:post_id", middleware.Authorization(), func(c *gin.Context) {
		id, _ := c.Get("id")
		postID := c.Param("post_id")

		var post model.Post
		if err := db.Where("id = ?", postID).First(&post).Error; err != nil {
			utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
			return
		}

		if post.UserID != id.(uuid.UUID) {
			utils.HttpRespFailed(c, http.StatusForbidden, "You are not authorized to delete this post")
			return
		}

		if err := db.Delete(&post).Error; err != nil {
			utils.HttpRespFailed(c, http.StatusInternalServerError, err.Error())
			return
		}

		utils.HttpRespSuccess(c, http.StatusOK, "Success delete post", nil)
	})

	// user like status
	r.GET("/:post_id/likestatus", middleware.Authorization(), func(c *gin.Context) {
		id, _ := c.Get("id")
		postID := c.Param("post_id")

		var isExist model.UserLikePost
		if err := db.Where("user_id = ? AND post_id = ?", id, postID).First(&isExist).Error; err != nil {
			utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
			return
		}

		utils.HttpRespSuccess(c, http.StatusOK, "Boolean", isExist)
	})

	// user like post
	r.POST("/:post_id/like", middleware.Authorization(), func(c *gin.Context) {
		id, _ := c.Get("id")
		postID := c.Param("post_id")

		var isExist model.UserLikePost
		if err := db.Where("user_id = ? AND post_id = ?", id, postID).First(&isExist).Error; err == nil {
			utils.HttpRespFailed(c, http.StatusConflict, "You've already liked this post")
			return
		}

		var post model.Post
		if err := db.Preload("User").Where("id = ?", postID).First(&post).Error; err != nil {
			utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
			return
		}

		userLikePost := model.UserLikePost{
			UserID: id.(uuid.UUID),
			PostID: post.ID,
		}

		post.Likes += 1

		if err := db.Create(&userLikePost).Error; err != nil {
			utils.HttpRespFailed(c, http.StatusInternalServerError, err.Error())
			return
		}

		if err := db.Save(&post).Error; err != nil {
			utils.HttpRespFailed(c, http.StatusInternalServerError, err.Error())
			return
		}

		var totalComment int64
		if err := db.Table("comments").Where("post_id = ?", post.ID).Count(&totalComment).Error; err != nil {
			utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
			return
		}

		postResponse := model.PostResponse{
			ID:           post.ID,
			User:         post.User,
			UserID:       post.UserID,
			Content:      post.Content,
			Attachment:   post.Attachment,
			Likes:        post.Likes,
			Latitude:     post.Latitude,
			Longitude:    post.Longitude,
			Distance:     post.Distance,
			IsLiked:      true,
			TotalComment: totalComment,
			CreatedAt:    time.Now(),
		}

		utils.HttpRespSuccess(c, http.StatusOK, "Success like post", postResponse)
	})

	// user unlike post
	r.DELETE("/:post_id/unlike", middleware.Authorization(), func(c *gin.Context) {
		id, _ := c.Get("id")
		postID := c.Param("post_id")

		var userLikePost model.UserLikePost
		if err := db.Where("user_id = ? AND post_id = ?", id, postID).First(&userLikePost).Error; err != nil {
			utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
			return
		}

		if err := db.Where("user_id = ? AND post_id = ?", id, postID).Delete(&userLikePost).Error; err != nil {
			utils.HttpRespFailed(c, http.StatusInternalServerError, err.Error())
			return
		}

		var post model.Post
		if err := db.Preload("User").Where("id = ?", postID).First(&post).Error; err != nil {
			utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
			return
		}

		post.Likes -= 1

		if err := db.Save(&post).Error; err != nil {
			utils.HttpRespFailed(c, http.StatusInternalServerError, err.Error())
			return
		}

		var totalComment int64
		if err := db.Table("comments").Where("post_id = ?", post.ID).Count(&totalComment).Error; err != nil {
			utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
			return
		}

		postResponse := model.PostResponse{
			ID:           post.ID,
			User:         post.User,
			UserID:       post.UserID,
			Content:      post.Content,
			Attachment:   post.Attachment,
			Likes:        post.Likes,
			Latitude:     post.Latitude,
			Longitude:    post.Longitude,
			Distance:     post.Distance,
			IsLiked:      false,
			TotalComment: totalComment,
			CreatedAt:    time.Now(),
		}

		utils.HttpRespSuccess(c, http.StatusOK, "Success unlike post", postResponse)
	})
}
