package routes

import (
	"github.com/HudYuSa/comments/pkg/controllers"
	"github.com/HudYuSa/comments/pkg/middleware"
	"github.com/gin-gonic/gin"
)

type CommentRoutes interface {
	SetupRoutes(rg *gin.RouterGroup)
}

type commentRoutes struct {
	CommentController controllers.CommentController
}

func NewCommentRoutes(commentController controllers.CommentController) CommentRoutes {
	return &commentRoutes{
		CommentController: commentController,
	}
}

func (cr *commentRoutes) SetupRoutes(rg *gin.RouterGroup) {
	router := rg.Group("comments")

	router.Use(middleware.DeserializeUser())
	router.POST("/", cr.CommentController.CreateComment)
	router.GET("/:post_id", cr.CommentController.GetComments)
	router.DELETE("/:comment_id", cr.CommentController.DeleteComment)
}
