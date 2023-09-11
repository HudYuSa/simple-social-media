package main

import (
	"log"
	"net/http"
	"time"

	"github.com/HudYuSa/comments/internal/config"
	"github.com/HudYuSa/comments/internal/connection"
	"github.com/HudYuSa/comments/pkg/controllers"
	"github.com/HudYuSa/comments/pkg/routes"
	"github.com/HudYuSa/comments/pkg/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	server *gin.Engine

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

	// // middleware
	server.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8000", config.GlobalConfig.ClientOrigin, "http://localhost:5173"},
		AllowMethods:     []string{"POST", "OPTIONS", "GET", "PUT", "PATCH"},
		AllowHeaders:     []string{"Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", " Authorization", " accept", "origin", "Cache-Control", " X-Requested-With", "ngrok-skip-browser-warning"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour}))

	// Handle CORS preflight requests
	server.OPTIONS("/*any", func(ctx *gin.Context) {
		ctx.Status(http.StatusOK)
	})

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
