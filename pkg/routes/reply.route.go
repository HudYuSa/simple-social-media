package routes

import (
	"github.com/HudYuSa/comments/pkg/controllers"
	"github.com/HudYuSa/comments/pkg/middleware"
	"github.com/gin-gonic/gin"
)

type ReplyRoutes interface {
	SetupRoutes(rg *gin.RouterGroup)
}

type replyRoutes struct {
	ReplyController controllers.ReplyController
}

func NewReplyRoutes(replyController controllers.ReplyController) ReplyRoutes {
	return &replyRoutes{
		ReplyController: replyController,
	}
}

func (rr *replyRoutes) SetupRoutes(rg *gin.RouterGroup) {
	router := rg.Group("replies")

	router.Use(middleware.DeserializeUser())
	router.POST("", rr.ReplyController.CreateReply)
	router.GET("/:comment_id", rr.ReplyController.GetReplies)
	router.DELETE("/:reply_id", rr.ReplyController.DeleteReply)
}
