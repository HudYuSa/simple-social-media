package controllers

import (
	"net/http"

	"github.com/HudYuSa/comments/database/models"
	"github.com/HudYuSa/comments/pkg/dtos"
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

	dtos.RespondWithJson(ctx, http.StatusOK, gin.H{
		"user": dtos.UserToUserResponse(&currentUser),
	})
}
