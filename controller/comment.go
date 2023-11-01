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

	r.GET("/:post_id/:comment_id", middleware.Authorization(), func(c *gin.Context) {
		commentID := c.Param("comment_id")

		var comment model.Comment
		if err := db.Where("id = ?", commentID).First(&comment).Error; err != nil {
			utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
			return
		}

		utils.HttpRespSuccess(c, http.StatusOK, "Success get comment", comment)
	})

	r.POST("/:post_id/comment", middleware.Authorization(), func(c *gin.Context) {
		randomID := utils.GenerateStringID()

		content := c.PostForm("content")

		attachment, _ := c.FormFile("attachment")
		filename := strings.ReplaceAll(strings.TrimSpace(attachment.Filename), " ", "")
		newFilename := randomID + "_" + filename
		attachment.Filename = newFilename

		_, err := supClient.Upload(attachment)
		if err != nil {
			utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
			return
		}

		userID, _ := c.Get("id")
		postID := c.Param("post_id")

		newComment := model.Comment{
			ID:         utils.GenerateStringID(),
			PostID:     postID,
			UserID:     userID.(uuid.UUID),
			Content:    content,
			Attachment: newFilename,
			Like:       0,
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

	// user like comment
	r.POST("/:post_id/:comment_id/like", middleware.Authorization(), func(c *gin.Context) {
		id, _ := c.Get("id")
		commentID := c.Param("comment_id")

		var isExist model.UserLikeComment
		if err := db.Where("user_id = ? AND comment_id = ?", id, commentID).First(&isExist).Error; err == nil {
			utils.HttpRespFailed(c, http.StatusConflict, "You've already liked this comment")
			return
		}

		var comment model.Comment
		if err := db.Where("id = ?", commentID).First(&comment).Error; err != nil {
			utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
			return
		}

		comment.Like += 1

		userLikeComment := model.UserLikeComment{
			UserID:    id.(uuid.UUID),
			CommentID: comment.ID,
		}

		if err := db.Create(&userLikeComment).Error; err != nil {
			utils.HttpRespFailed(c, http.StatusInternalServerError, err.Error())
			return
		}

		if err := db.Save(&comment).Error; err != nil {
			utils.HttpRespFailed(c, http.StatusInternalServerError, err.Error())
			return
		}

		utils.HttpRespSuccess(c, http.StatusOK, "Success like comment", comment)
	})

	r.DELETE("/:post_id/:comment_id/unlike", middleware.Authorization(), func(c *gin.Context) {
		id, _ := c.Get("id")
		commentID := c.Param("comment_id")

		var isExist model.UserLikeComment
		if err := db.Where("user_id = ? AND comment_id = ?", id, commentID).First(&isExist).Error; err != nil {
			utils.HttpRespFailed(c, http.StatusConflict, "You've already unliked this comment")
			return
		}

		var comment model.Comment
		if err := db.Where("id = ?", commentID).First(&comment).Error; err != nil {
			utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
			return
		}

		comment.Like -= 1

		if err := db.Where("user_id = ? AND comment_id = ?", id, commentID).Delete(&model.UserLikeComment{}).Error; err != nil {
			utils.HttpRespFailed(c, http.StatusInternalServerError, err.Error())
			return
		}

		if err := db.Save(&comment).Error; err != nil {
			utils.HttpRespFailed(c, http.StatusInternalServerError, err.Error())
			return
		}

		utils.HttpRespSuccess(c, http.StatusOK, "Success unlike comment", comment)
	})
}
