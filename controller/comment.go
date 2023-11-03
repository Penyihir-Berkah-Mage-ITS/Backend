package controller

import (
	"Backend/middleware"
	"Backend/model"
	"Backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"net/http"
	"time"
)

func Comment(db *gorm.DB, q *gin.Engine) {
	r := q.Group("/api/v1/post")
	r.GET("/:post_id/comment", middleware.Authorization(), func(c *gin.Context) {
		postID := c.Param("post_id")

		var comments []model.Comment
		if err := db.Preload("User").Where("post_id = ?", postID).Find(&comments).Error; err != nil {
			utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
			return
		}

		var commentsResponse []model.CommentResponse
		for _, comment := range comments {
			commentResponse := model.CommentResponse{
				ID:        comment.ID,
				UserID:    comment.UserID,
				User:      comment.User,
				PostID:    comment.PostID,
				Content:   comment.Content,
				Like:      comment.Like,
				CreatedAt: comment.CreatedAt,
				UpdatedAt: comment.UpdatedAt,
			}

			id, _ := c.Get("id")

			var isExist model.UserLikeComment
			if err := db.Where("user_id = ? AND comment_id = ?", id, comment.ID).First(&isExist).Error; err != nil {
				commentResponse.IsLiked = false
			} else {
				commentResponse.IsLiked = true
			}

			commentsResponse = append(commentsResponse, commentResponse)
		}

		utils.HttpRespSuccess(c, http.StatusOK, "Success get comments", commentsResponse)
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
		var commentInput model.CommentInput
		if err := c.BindJSON(&commentInput); err != nil {
			utils.HttpRespFailed(c, http.StatusUnprocessableEntity, err.Error())
			return
		}

		userID, _ := c.Get("id")
		postID := c.Param("post_id")

		newComment := model.Comment{
			ID:        utils.GenerateStringID(),
			PostID:    postID,
			UserID:    userID.(uuid.UUID),
			Content:   commentInput.Content,
			Like:      0,
			CreatedAt: time.Now(),
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

	// user like status
	r.GET("/:post_id/:comment_id/likestatus", middleware.Authorization(), func(c *gin.Context) {
		id, _ := c.Get("id")
		commentID := c.Param("comment_id")

		var isExist model.UserLikeComment
		if err := db.Where("user_id = ? AND comment_id = ?", id, commentID).First(&isExist).Error; err != nil {
			utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
			return
		}

		utils.HttpRespSuccess(c, http.StatusOK, "Boolean", isExist)
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

	// get total comment from post
	r.GET("/:post_id/totalcomment", middleware.Authorization(), func(c *gin.Context) {
		postID := c.Param("post_id")
		var commentsCount int64

		if err := db.Table("comments").Where("post_id = ?", postID).Count(&commentsCount).Error; err != nil {
			utils.HttpRespFailed(c, http.StatusInternalServerError, err.Error())
			return
		}

		utils.HttpRespSuccess(c, http.StatusOK, "Success get total comments", commentsCount)
	})
}
