package controllers

import (
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/HudYuSa/comments/database/models"
	"github.com/HudYuSa/comments/internal/config"
	"github.com/HudYuSa/comments/pkg/dtos"
	"github.com/HudYuSa/comments/pkg/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthController interface {
	SignUpUser(ctx *gin.Context)
	SignInUser(ctx *gin.Context)
	RefreshAccessToken(ctx *gin.Context)
	LogOutUser(ctx *gin.Context)
	DeleteUser(ctx *gin.Context)
	// VerifyEmail(ctx *gin.Context)
	// ForgotPassword(ctx *gin.Context)
	// ResetPassword(ctx *gin.Context)
}

type authController struct {
	DB *gorm.DB
}

func NewAuthController(db *gorm.DB) AuthController {
	return &authController{DB: db}
}

func (ac *authController) SignUpUser(ctx *gin.Context) {
	// payload is like request body
	var payload *dtos.SignUpInput

	// try to bind the request body to the payload struct
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		dtos.RespondWithError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	// match the password with confirm password
	if payload.Password != payload.PasswordConfirm {
		dtos.RespondWithError(ctx, http.StatusBadRequest, "email or password is invalid")
		return
	}

	// hash the user password
	hashedPassword, err := utils.HashPassword(payload.Password)
	if err != nil {
		dtos.RespondWithError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	now := time.Now()

	newUser := models.User{
		Name:      payload.Name,
		Email:     strings.ToLower(payload.Email),
		Password:  hashedPassword,
		Photo:     payload.Photo,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// save the data to database using gorm
	result := ac.DB.Create(&newUser)

	// check for any possible error
	if result.Error != nil && strings.Contains(result.Error.Error(), "duplicate key value violates unique") {
		dtos.RespondWithError(ctx, http.StatusConflict, "user with that email already exist")
		return
	} else if result.Error != nil {
		dtos.RespondWithError(ctx, http.StatusBadGateway, "something bad just happen")
		return
	}

	// send the response
	dtos.RespondWithJson(ctx, http.StatusCreated, dtos.UserToUserResponse(&newUser))
}

func (ac *authController) SignInUser(ctx *gin.Context) {
	var payload *dtos.SignInInput

	// try to bind the request body to the payload struct
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		dtos.RespondWithError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	var dbUser models.User

	// find the user in db and bind it to user variable if exist
	result := ac.DB.First(&dbUser, "email = ?", strings.ToLower(payload.Email))

	// if there's an error when fetching from db
	if result.Error != nil {
		dtos.RespondWithError(ctx, http.StatusBadGateway, "invalid email or password")
		return
	}

	// verify user password
	if err := utils.VerifyPassword(dbUser.Password, payload.Password); err != nil {
		dtos.RespondWithError(ctx, http.StatusBadGateway, "invalid email or password")
		return
	}

	// Generate Tokens
	accessToken, err := utils.CreateToken(config.GlobalConfig.AccessTokenExpiresIn, dbUser.ID, config.GlobalConfig.AccessTokenPrivateKey)
	if err != nil {
		dtos.RespondWithError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	refreshToken, err := utils.CreateToken(config.GlobalConfig.RefreshTokenExpiresIn, dbUser.ID, config.GlobalConfig.RefreshTokenPrivateKey)
	if err != nil {
		dtos.RespondWithError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	// set accesstoken and refresh token to client cookie
	// max age time 60 so it become minute
	ctx.SetCookie("access_token", accessToken, config.GlobalConfig.AccessTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", refreshToken, config.GlobalConfig.RefreshTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "true", config.GlobalConfig.AccessTokenMaxAge, "/", "localhost", false, false)

	dtos.RespondWithJson(ctx, http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (ac *authController) RefreshAccessToken(ctx *gin.Context) {

	refreshToken := utils.GetToken(ctx, "refresh_token", "x-refresh-token")

	// if there's no token from header or cookie
	if reflect.ValueOf(refreshToken).IsZero() {
		dtos.RespondWithError(ctx, http.StatusUnauthorized, "you're not allowed to access this endpoint")
		return
	}

	// validate the token
	sub, err := utils.ValidateToken(refreshToken, config.GlobalConfig.RefreshTokenPublicKey)

	if err != nil {
		dtos.RespondWithError(ctx, http.StatusUnauthorized, err.Error())
		return
	}

	// find the user that has the refresh token
	var user models.User

	result := ac.DB.First(&user, "id = ?", sub)
	if result.Error != nil {
		dtos.RespondWithError(ctx, http.StatusNotFound, "the user belonging to this token doesn't exist anymore")
		return
	}

	// reissue new accesstoken
	accessToken, err := utils.CreateToken(config.GlobalConfig.AccessTokenExpiresIn, user.ID, config.GlobalConfig.AccessTokenPrivateKey)
	if err != nil {
		dtos.RespondWithError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	// set new accesstoken cookie
	ctx.SetCookie("access_token", accessToken, config.GlobalConfig.AccessTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "true", config.GlobalConfig.AccessTokenMaxAge*60, "/", "localhost", false, false)

	dtos.RespondWithJson(ctx, http.StatusOK, gin.H{
		"access_token": accessToken,
	})
}

func (ac *authController) LogOutUser(ctx *gin.Context) {
	ctx.SetCookie("access_token", "", -1, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", "", -1, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "", -1, "/", "localhost", false, false)

	dtos.RespondWithJson(ctx, http.StatusOK, "successfully logout user")
}

func (ac *authController) DeleteUser(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)

	// find user
	var user models.User
	result := ac.DB.Where("id = ?", currentUser.ID).Delete(&user)
	if result.Error != nil {
		switch result.Error.Error() {
		case "record not found":
			dtos.RespondWithError(ctx, http.StatusNotFound, "cannot delete user")
		default:
			dtos.RespondWithError(ctx, http.StatusBadGateway, result.Error.Error())
		}
		return
	}

	dtos.RespondWithJson(ctx, http.StatusNoContent, "successfully delete user")
}

// func (ac *authController) ForgotPassword(ctx *gin.Context) {
// }

// func (ac *authController) ResetPassword(ctx *gin.Context) {
// 	panic("not implemented") // TODO: Implement
// }
