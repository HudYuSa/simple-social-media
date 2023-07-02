package main

import (
	"log"
	"net/http"

	"github.com/HudYuSa/mod-name/internal/config"
	"github.com/HudYuSa/mod-name/internal/connection"
	"github.com/HudYuSa/mod-name/pkg/controllers"
	"github.com/HudYuSa/mod-name/pkg/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	server     *gin.Engine
	corsConfig cors.Config

	UserController controllers.UserController
	UserRoutes     routes.UserRoutes

	AuthController controllers.AuthController
	AuthRoutes     routes.AuthRoutes
)

func init() {
	server = gin.Default()
	err := config.LoadConfig(".env")
	if err != nil {
		log.Fatal("? Could not load environment variables ", err)
	}
	corsConfig = cors.DefaultConfig()

	// connect ke database
	connection.ConnectDB(&config.GlobalConfig)

	UserController = controllers.NewUserController(connection.DB)
	UserRoutes = routes.NewUserRoutes(UserController)

	AuthController = controllers.NewAuthController(connection.DB)
	AuthRoutes = routes.NewAuthRoutes(AuthController)
}

func main() {
	//set allowed origins
	corsConfig.AllowOrigins = []string{"http://localhost:8000", config.GlobalConfig.ClientOrigin}

	// middleware
	server.Use(cors.New(corsConfig))

	router := server.Group("/api")
	router.GET("/test", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "welcome to this project"})
	})

	// Routes
	UserRoutes.SetupRoutes(router)
	AuthRoutes.SetupRoutes(router)

	// run app
	log.Fatal(server.Run(":" + config.GlobalConfig.ServerPort))
}
