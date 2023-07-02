package controllers

import (
	"net/http"

	"github.com/HudYuSa/mod-name/database/models"
	"github.com/HudYuSa/mod-name/pkg/dtos"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserController interface {
	GetMe(ctx *gin.Context)
}

type userController struct {
	DB *gorm.DB
}

func NewUserController(db *gorm.DB) UserController {
	return &userController{
		DB: db,
	}
}

func (uc *userController) GetMe(ctx *gin.Context) {
	currentUser, ok := ctx.MustGet("currentUser").(models.User)

	if !ok {
		dtos.RespondWithError(ctx, http.StatusInternalServerError, "cannot get currentUser")
		return
	}

	userResponse := dtos.DBUserToUserResponse(&currentUser)

	dtos.RespondWithJson(ctx, http.StatusOK, dtos.WebResponse{
		Message: "Successfully get user data",
		Data: gin.H{
			"user": userResponse,
		},
	})
}
