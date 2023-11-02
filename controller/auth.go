package controller

import (
	"Backend/model"
	"Backend/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"log"
	"net/http"
	"net/smtp"
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
			ProfilePicture: utils.DefaultAvatar(userRegister.ProfilePicture),
			AccountType:    "ASA Peeps",
			Gender:         userRegister.Gender,
			IsVerified:     false,
			CreatedAt:      time.Now(),
		}

		if err := db.Create(&newUser).Error; err != nil {
			utils.HttpRespFailed(c, http.StatusInternalServerError, err.Error())
			return
		}

		sendVerificationEmail(newUser.ID.String(), newUser.Email)

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

func Verify(db *gorm.DB, q *gin.Engine) {
	r := q.Group("/api/v1")
	r.GET(":user_id/verify", func(c *gin.Context) {
		userID := c.Param("user_id")
		var user model.User
		if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
			utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
			return
		}

		if user.IsVerified {
			utils.HttpRespFailed(c, http.StatusBadRequest, "User already verified")
			return
		}

		user.IsVerified = true

		verified := model.UserVerify{
			UserID:    user.ID,
			Verify:    true,
			CreatedAt: time.Now(),
		}

		if err := db.Create(&verified).Error; err != nil {
			utils.HttpRespFailed(c, http.StatusInternalServerError, err.Error())
			return
		}

		if err := db.Save(&user).Error; err != nil {
			utils.HttpRespFailed(c, http.StatusInternalServerError, err.Error())
			return
		}

		utils.HttpRespSuccess(c, http.StatusOK, "Success verify user", user)
	})
}

func sendVerificationEmail(userID string, userEmail string) {
	auth := smtp.PlainAuth("", os.Getenv("SMTP_EMAIL"), os.Getenv("SMTP_EMAIL_PASSWORD"), "smtp.gmail.com")

	verifyLink := os.Getenv("HOST") + "/api/v1/" + userID + "/verify"

	message := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: Account Verification\r\n\r\n"+
			"Dear User,\r\n\r\n"+
			"Thank you for creating an account. To verify your account, please click the link below:\r\n"+
			"<a href='%s'>Verify Account</a>\r\n\r\n"+
			"If you didn't sign up for this account, please ignore this email.\r\n\r\n"+
			"Best regards,\r\n"+
			"%s",
		os.Getenv("EMAIL"), userEmail, verifyLink, "ASA",
	)

	err := smtp.SendMail("smtp.gmail.com:587", auth, os.Getenv("EMAIL"), []string{userEmail}, []byte(message))
	if err != nil {
		log.Println("Error sending verification email:", err)
		utils.HttpRespFailed(nil, http.StatusInternalServerError, err.Error())
	}
}
