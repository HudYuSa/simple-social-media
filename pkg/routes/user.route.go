package routes

import (
	"github.com/HudYuSa/mod-name/pkg/controllers"
	"github.com/HudYuSa/mod-name/pkg/middleware"
	"github.com/gin-gonic/gin"
)

type UserRoutes interface {
	SetupRoutes(rg *gin.RouterGroup)
}

type userRoutes struct {
	userController controllers.UserController
}

func NewUserRoutes(userController controllers.UserController) UserRoutes {
	return &userRoutes{
		userController: userController,
	}
}

func (ur *userRoutes) SetupRoutes(rg *gin.RouterGroup) {
	router := rg.Group("/users")

	router.GET("/me", middleware.DeserializeUser(), ur.userController.GetMe)
}
