package middleware

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/HudYuSa/comments/database/models"
	"github.com/HudYuSa/comments/internal/config"
	"github.com/HudYuSa/comments/internal/connection"
	"github.com/HudYuSa/comments/pkg/dtos"
	"github.com/HudYuSa/comments/pkg/utils"
	"github.com/gin-gonic/gin"
)

func DeserializeUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		accessToken := utils.GetToken(ctx, "access_token", "Authorization")

		// if there's no token from header or cookie
		if reflect.ValueOf(accessToken).IsZero() {
			dtos.RespondWithError(ctx, http.StatusUnauthorized, "you're not allowed to access this endpoint")
			return
		}

		// validate the token and get the user id from the sub/subject
		sub, err := utils.ValidateToken(accessToken, config.GlobalConfig.AccessTokenPublicKey)

		if err != nil {
			dtos.RespondWithError(ctx, http.StatusUnauthorized, err.Error())
			return
		}

		var user models.User
		result := connection.DB.First(&user, "id=?", fmt.Sprint(sub))

		if result.Error != nil {
			dtos.RespondWithError(ctx, http.StatusUnauthorized, "The user belonging to this token doesn't exist anymore")
			return
		}

		ctx.Set("currentUser", user)
		ctx.Next()
	}
}
