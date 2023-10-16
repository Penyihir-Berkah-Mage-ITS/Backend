package main

import (
	"Backend/database"
	"Backend/middleware"
	"Backend/model"
	"fmt"
	supabasestorageuploader "github.com/adityarizkyramadhan/supabase-storage-uploader"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err.Error())
	}

	supClient := supabasestorageuploader.New(
		os.Getenv("SUPABASE_PROJECT_URL"),
		os.Getenv("SUPABASE_PROJECT_API_KEY"),
		os.Getenv("SUPABASE_PROJECT_STORAGE_NAME"),
	)

	databaseConf, err := database.NewDatabase()
	if err != nil {
		panic(err.Error())
	}

	db, err := database.MakeSupaBaseConnectionDatabase(databaseConf)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println(db)

	r := gin.Default()

	r.Use(middleware.CORS())

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"env": os.Getenv("ENV"),
		})
	})

	r.GET("/users", func(c *gin.Context) {
		log.Println("Start")
		var users []model.User
		if err := db.Find(&users).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": users,
		})
		log.Println("End")
	})

	r.POST("/upload", func(c *gin.Context) {
		log.Println("Start")
		file, err := c.FormFile("avatar")
		if err != nil {
			c.JSON(400, gin.H{"data": err.Error()})
			return
		}
		link, err := supClient.Upload(file)
		if err != nil {
			c.JSON(500, gin.H{"data": err.Error()})
			return
		}
		c.JSON(200, gin.H{"data": link})
		log.Println("End")
	})

	if err := r.Run(fmt.Sprintf(":%s", os.Getenv("PORT"))); err != nil {
		panic(err.Error())
	}
}
