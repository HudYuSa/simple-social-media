package main

import (
	"log"
	"net/http"
	"time"

	"github.com/HudYuSa/comments/internal/config"
	"github.com/HudYuSa/comments/internal/connection"
	"github.com/HudYuSa/comments/pkg/controllers"
	"github.com/HudYuSa/comments/pkg/middleware"
	"github.com/HudYuSa/comments/pkg/routes"
	"github.com/HudYuSa/comments/pkg/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	server     *gin.Engine
	corsConfig cors.Config

	authController    controllers.AuthController
	userController    controllers.UserController
	postController    controllers.PostController
	commentController controllers.CommentController
	replyController   controllers.ReplyController

	authRoutes    routes.AuthRoutes
	userRoutes    routes.UserRoutes
	postRoutes    routes.PostRoutes
	commentRoutes routes.CommentRoutes
	replyRoutes   routes.ReplyRoutes
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
	commentController = controllers.NewCommentController(connection.DB)
	replyController = controllers.NewReplyController(connection.DB)

	authRoutes = routes.NewAuthRoutes(authController)
	userRoutes = routes.NewUserRoutes(userController)
	postRoutes = routes.NewPostRoutes(postController)
	commentRoutes = routes.NewCommentRoutes(commentController)
	replyRoutes = routes.NewReplyRoutes(replyController)
}

func main() {
	//set allowed origins
	corsConfig.AllowOrigins = []string{"http://localhost:8000", config.GlobalConfig.ClientOrigin, "http://localhost:5173"}

	// middleware
	server.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8000", config.GlobalConfig.ClientOrigin, "http://localhost:5173"},
		AllowMethods:     []string{"PUT", "PATCH"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour}))
	server.Use(middleware.CORSMiddleware())

	router := server.Group("/api")
	router.GET("/test", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "welcome to this project"})
	})

	// Routes
	authRoutes.SetupRoutes(router)
	userRoutes.SetupRoutes(router)
	postRoutes.SetupRoutes(router)
	commentRoutes.SetupRoutes(router)
	replyRoutes.SetupRoutes(router)

	// firestore
	client, err := utils.InitializeFirestore()
	if err != nil {
		log.Fatal(err)
	}
	// ngrok
	tun, err := utils.RunNgrok()
	if err != nil {
		log.Fatal(err)
	}

	// update firestore url
	err = utils.UpdateUrl(client, tun)
	if err != nil {
		log.Fatal(err)
	}

	// run app
	log.Fatal(http.Serve(tun, server.Handler()))
	// log.Fatal(server.Run(":" + config.GlobalConfig.ServerPort))
}
