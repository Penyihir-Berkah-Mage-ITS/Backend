package main

import (
	"Backend/controller"
	"Backend/controller/websocket"
	"Backend/database"
	"Backend/middleware"
	"fmt"
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

	databaseConf, err := database.NewDatabase()
	if err != nil {
		panic(err.Error())
	}

	db, err := database.MakeSupaBaseConnectionDatabase(databaseConf)
	if err != nil {
		panic(err.Error())
	}

	hub := websocket.NewHub()
	wsHandler := websocket.NewHandler(hub)
	go hub.Run()

	r := gin.Default()

	r.Use(middleware.CORS())

	controller.Register(db, r)
	controller.Login(db, r)
	controller.Verify(db, r)
	controller.Profile(db, r)
	controller.Home(db, r)
	controller.Post(db, r)
	controller.Comment(db, r)
	controller.Report(db, r)
	websocket.WebSocketInit(r, wsHandler)
	controller.Chat(db, r)

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"env": os.Getenv("ENV"),
		})
	})

	if err := r.Run(fmt.Sprintf(":%s", os.Getenv("PORT"))); err != nil {
		panic(err.Error())
	}
}
