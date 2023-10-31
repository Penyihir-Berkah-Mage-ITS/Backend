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
		content := c.PostForm("content")

		attachment, _ := c.FormFile("attachment")
		uploadedAttachment, err := supClient.Upload(attachment)
		if err != nil {
			utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
			return
		}

		id, _ := c.Get("id")

		newPost := model.Post{
			ID:         utils.GenerateStringID(),
			UserID:     id.(uuid.UUID),
			Content:    content,
			Attachment: uploadedAttachment,
			Like:       0,
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
		if err := db.Where("id = ?", postID).First(&post).Error; err != nil {
			utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
			return
		}

		userLikePost := model.UserLikePost{
			UserID: id.(uuid.UUID),
			PostID: post.ID,
		}

		post.Like += 1

		if err := db.Create(&userLikePost).Error; err != nil {
			utils.HttpRespFailed(c, http.StatusInternalServerError, err.Error())
			return
		}

		if err := db.Save(&post).Error; err != nil {
			utils.HttpRespFailed(c, http.StatusInternalServerError, err.Error())
			return
		}

		utils.HttpRespSuccess(c, http.StatusOK, "Success like post", post)
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
		if err := db.Where("id = ?", postID).First(&post).Error; err != nil {
			utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
			return
		}

		post.Like -= 1

		if err := db.Save(&post).Error; err != nil {
			utils.HttpRespFailed(c, http.StatusInternalServerError, err.Error())
			return
		}

		utils.HttpRespSuccess(c, http.StatusOK, "Success unlike post", post)
	})
}
