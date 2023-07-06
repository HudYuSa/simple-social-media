package main

import (
	"log"
	"net/http"

	"github.com/HudYuSa/comments/internal/config"
	"github.com/HudYuSa/comments/internal/connection"
	"github.com/HudYuSa/comments/pkg/controllers"
	"github.com/HudYuSa/comments/pkg/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	server     *gin.Engine
	corsConfig cors.Config

	authController controllers.AuthController
	userController controllers.UserController
	postController controllers.PostController

	authRoutes routes.AuthRoutes
	userRoutes routes.UserRoutes
	postRoutes routes.PostRoutes
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

	authController = controllers.NewAuthController(connection.DB)
	userController = controllers.NewUserController(connection.DB)
	postController = controllers.NewPostController(connection.DB)

	authRoutes = routes.NewAuthRoutes(authController)
	userRoutes = routes.NewUserRoutes(userController)
	postRoutes = routes.NewPostRoutes(postController)
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
	authRoutes.SetupRoutes(router)
	userRoutes.SetupRoutes(router)
	postRoutes.SetupRoutes(router)

	// run app
	log.Fatal(server.Run(":" + config.GlobalConfig.ServerPort))
}
