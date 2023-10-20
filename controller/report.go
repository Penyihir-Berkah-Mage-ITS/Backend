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
)

func Report(db *gorm.DB, q *gin.Engine) {
	supClient := supabasestorageuploader.New(
		os.Getenv("SUPABASE_PROJECT_URL"),
		os.Getenv("SUPABASE_PROJECT_API_KEY"),
		os.Getenv("SUPABASE_PROJECT_STORAGE_NAME"),
	)

	r := q.Group("/api/v1/report")
	r.POST("/create", middleware.Authorization(), func(c *gin.Context) {
		id, _ := c.Get("id")

		name := c.PostForm("name")
		address := c.PostForm("address")
		province := c.PostForm("province")
		city := c.PostForm("city")
		phone := c.PostForm("phone")
		detailReport := c.PostForm("detail_report")
		proof, _ := c.FormFile("proof")
		uploadedProof, err := supClient.Upload(proof)
		if err != nil {
			utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
			return
		}

		newReport := model.Report{
			ID:           uuid.New(),
			UserID:       id.(uuid.UUID),
			Name:         name,
			Address:      address,
			Province:     province,
			City:         city,
			Phone:        phone,
			DetailReport: detailReport,
			Proof:        uploadedProof,
		}

		if err := db.Create(&newReport).Error; err != nil {
			utils.HttpRespFailed(c, http.StatusInternalServerError, err.Error())
			return
		}

		utils.HttpRespSuccess(c, http.StatusOK, "Success create new report", newReport)
	})
}
