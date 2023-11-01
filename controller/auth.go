package controller

import (
	"Backend/model"
	"Backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"net/http"
	"os"
	"time"
)

func Register(db *gorm.DB, q *gin.Engine) {
	r := q.Group("/api/v1")
	r.POST("/register", func(c *gin.Context) {
		var userRegister model.UserRegisterInput

		if err := c.BindJSON(&userRegister); err != nil {
			utils.HttpRespFailed(c, http.StatusUnprocessableEntity, err.Error())
		}

		hashed, err := utils.Hash(userRegister.Password)
		if err != nil {
			utils.HttpRespFailed(c, http.StatusInternalServerError, err.Error())
			return
		}

		newUser := model.User{
			ID:             uuid.New(),
			Username:       userRegister.Username,
			Phone:          userRegister.Phone,
			Email:          userRegister.Email,
			Password:       hashed,
			ProfilePicture: userRegister.ProfilePicture,
			AccountType:    "ASA Peeps",
			Gender:         userRegister.Gender,
			IsVerified:     false,
			CreatedAt:      time.Now(),
		}

		if err := db.Create(&newUser).Error; err != nil {
			utils.HttpRespFailed(c, http.StatusInternalServerError, err.Error())
			return
		}

		utils.HttpRespSuccess(c, http.StatusOK, "Success create new user", newUser)
	})
}

func Login(db *gorm.DB, q *gin.Engine) {
	r := q.Group("/api/v1")
	r.POST("/login", func(c *gin.Context) {
		var userLogin model.UserLoginInput

		if err := c.BindJSON(&userLogin); err != nil {
			utils.HttpRespFailed(c, http.StatusUnprocessableEntity, err.Error())
		}

		var user model.User

		if err := db.Where("username = ?", userLogin.Username).First(&user).Error; err != nil {
			utils.HttpRespFailed(c, http.StatusInternalServerError, err.Error())
			return
		}

		if !utils.CompareHash(userLogin.Password, user.Password) {
			utils.HttpRespFailed(c, http.StatusUnauthorized, "Password is wrong")
			return
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
			"id":   user.ID,
			"type": "ASA Peeps",
			"exp":  time.Now().Add(time.Hour).Unix(),
		})

		strToken, err := token.SignedString([]byte(os.Getenv("TOKEN")))
		if err != nil {
			utils.HttpRespFailed(c, http.StatusInternalServerError, err.Error())
			return
		}

		utils.HttpRespSuccess(c, http.StatusOK, "Parsed token", gin.H{
			"name":  user.Username,
			"token": strToken,
			"type":  "ASA Peeps",
		})

		return
	})
}
