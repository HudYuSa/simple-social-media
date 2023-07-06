package routes

import (
	"github.com/HudYuSa/comments/pkg/controllers"
	"github.com/HudYuSa/comments/pkg/middleware"
	"github.com/gin-gonic/gin"
)

type PostRoutes interface {
	SetupRoutes(rg *gin.RouterGroup)
}

type postRoutes struct {
	PostController controllers.PostController
}

func NewPostRoutes(postController controllers.PostController) PostRoutes {
	return &postRoutes{
		PostController: postController,
	}
}

func (pr postRoutes) SetupRoutes(rg *gin.RouterGroup) {
	router := rg.Group("posts")

	router.Use(middleware.DeserializeUser())
	router.POST("/", pr.PostController.CreatePost)
	router.GET("/", pr.PostController.GetAllPost)
	router.GET("/:post_id", pr.PostController.GetPostByID)
	router.PATCH("/:post_id", pr.PostController.UpdatePost)
	router.DELETE("/:post_id", pr.PostController.DeletePost)
}
