package controller

import (
	"Backend/middleware"
	"Backend/model"
	"Backend/utils"
	supabasestorageuploader "github.com/adityarizkyramadhan/supabase-storage-uploader"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"os"
	"time"
)

func Profile(db *gorm.DB, q *gin.Engine) {
	supClient := supabasestorageuploader.New(
		os.Getenv("SUPABASE_PROJECT_URL"),
		os.Getenv("SUPABASE_PROJECT_API_KEY"),
		os.Getenv("SUPABASE_PROJECT_STORAGE_NAME"),
	)

	r := q.Group("/api/v1/user")
	r.GET("/profile", middleware.Authorization(), func(c *gin.Context) {
		ID, _ := c.Get("id")

		var user model.User
		if err := db.Where("id = ?", ID).First(&user).Error; err != nil {
			utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
			return
		}

		utils.HttpRespSuccess(c, http.StatusOK, "Success get user profile", user)
	})

	r.POST("/edit-profile", middleware.Authorization(), func(c *gin.Context) {
		var userEditInput model.UserEditInput
		if err := c.BindJSON(&userEditInput); err != nil {
			utils.HttpRespFailed(c, http.StatusUnprocessableEntity, err.Error())
			return
		}

		ID, _ := c.Get("id")
		var user model.User
		if err := db.Where("id = ?", ID).First(&user).Error; err != nil {
			utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
			return
		}

		user.Username = userEditInput.Username
		user.Phone = userEditInput.Phone
		user.Gender = userEditInput.Gender

		user.UpdatedAt = time.Now()

		if err := db.Save(&user).Error; err != nil {
			utils.HttpRespFailed(c, http.StatusInternalServerError, err.Error())
			return
		}

		utils.HttpRespSuccess(c, http.StatusOK, "Success edit user profile", user)
	})

	r.POST("/edit-profile-picture", middleware.Authorization(), func(c *gin.Context) {
		picture, _ := c.FormFile("picture")
		uploaded, err := supClient.Upload(picture)
		if err != nil {
			utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
			return
		}

		id, _ := c.Get("id")
		var user model.User
		if err := db.Where("id = ?", id).First(&user).Error; err != nil {
			utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
			return
		}

		user.ProfilePicture = uploaded
		user.UpdatedAt = time.Now()

		if err := db.Save(&user).Error; err != nil {
			utils.HttpRespFailed(c, http.StatusInternalServerError, err.Error())
			return
		}

		utils.HttpRespSuccess(c, http.StatusOK, "Success edit user profile picture", user)
	})

	r.POST("/edit-password", middleware.Authorization(), func(c *gin.Context) {
		id, _ := c.Get("id")
		var user model.User
		if err := db.Where("id = ?", id).First(&user).Error; err != nil {
			utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
			return
		}

		var userEditPassword model.UserEditPassword
		if err := c.BindJSON(&userEditPassword); err != nil {
			utils.HttpRespFailed(c, http.StatusUnprocessableEntity, err.Error())
			return
		}

		if !utils.CompareHash(userEditPassword.OldPassword, user.Password) {
			utils.HttpRespFailed(c, http.StatusUnauthorized, "Password is wrong")
			return
		}

		if userEditPassword.NewPassword != userEditPassword.ConfirmNewPassword {
			utils.HttpRespFailed(c, http.StatusUnauthorized, "Confirm password is wrong")
			return
		}

		hashed, err := utils.Hash(userEditPassword.NewPassword)
		if err != nil {
			utils.HttpRespFailed(c, http.StatusInternalServerError, err.Error())
			return
		}

		user.Password = hashed

		user.UpdatedAt = time.Now()

		if err := db.Save(&user).Error; err != nil {
			utils.HttpRespFailed(c, http.StatusInternalServerError, err.Error())
			return
		}

		utils.HttpRespSuccess(c, http.StatusOK, "Success edit user password", user)
	})
}
