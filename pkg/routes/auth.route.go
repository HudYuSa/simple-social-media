package routes

import (
	"github.com/HudYuSa/mod-name/pkg/controllers"
	"github.com/gin-gonic/gin"
)

type AuthRoutes interface {
	SetupRoutes(rg *gin.RouterGroup)
}

type authRoutes struct {
	AuthController controllers.AuthController
}

func NewAuthRoutes(authController controllers.AuthController) AuthRoutes {
	return &authRoutes{
		AuthController: authController,
	}
}

func (ar *authRoutes) SetupRoutes(rg *gin.RouterGroup) {
	router := rg.Group("/auth")

	router.POST("/signup", ar.AuthController.SignUpUser)
	router.POST("/login", ar.AuthController.SignInUser)
}
