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
	"strconv"
	"time"
)

func Comment(db *gorm.DB, q *gin.Engine) {
	supClient := supabasestorageuploader.New(
		os.Getenv("SUPABASE_PROJECT_URL"),
		os.Getenv("SUPABASE_PROJECT_API_KEY"),
		os.Getenv("SUPABASE_PROJECT_STORAGE_NAME"),
	)

	r := q.Group("/api/v1/post")
	r.GET("/:post_id/comment", middleware.Authorization(), func(c *gin.Context) {
		postID := c.Param("post_id")

		var comments []model.Comment
		if err := db.Where("post_id = ?", postID).Find(&comments).Error; err != nil {
			utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
			return
		}

		utils.HttpRespSuccess(c, http.StatusOK, "Success get comments", comments)
	})

	r.POST("/:post_id/comment", middleware.Authorization(), func(c *gin.Context) {
		content := c.PostForm("content")

		attachment, _ := c.FormFile("attachment")
		uploadedAttachment, err := supClient.Upload(attachment)
		if err != nil {
			utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
			return
		}

		userID, _ := c.Get("id")
		postID := c.Param("post_id")

		newComment := model.Comment{
			ID:         strconv.FormatInt(utils.GenerateID(), 10),
			PostID:     postID,
			UserID:     userID.(uuid.UUID),
			Content:    content,
			Attachment: uploadedAttachment,
			CreatedAt:  time.Now(),
		}

		if err := db.Create(&newComment).Error; err != nil {
			utils.HttpRespFailed(c, http.StatusInternalServerError, err.Error())
			return
		}

		utils.HttpRespSuccess(c, http.StatusOK, "Success created comment", newComment)
	})

	r.DELETE("/:post_id/delete/:comment_id", middleware.Authorization(), func(c *gin.Context) {
		id, _ := c.Get("id")
		commentID := c.Param("comment_id")
		postID := c.Param("post_id")

		var comment model.Comment
		if err := db.Where("id = ? AND post_id = ?", commentID, postID).First(&comment).Error; err != nil {
			utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
			return
		}

		if comment.UserID != id.(uuid.UUID) {
			utils.HttpRespFailed(c, http.StatusForbidden, "You are not authorized to delete this comment")
			return
		}

		if err := db.Delete(&comment).Error; err != nil {
			utils.HttpRespFailed(c, http.StatusInternalServerError, err.Error())
			return
		}

		utils.HttpRespSuccess(c, http.StatusOK, "Success delete comment", comment)
	})
}
