package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/HudYuSa/mod-name/database/models"
	"github.com/HudYuSa/mod-name/internal/config"
	"github.com/HudYuSa/mod-name/internal/connection"
	"github.com/HudYuSa/mod-name/pkg/dtos"
	"github.com/HudYuSa/mod-name/pkg/utils"
	"github.com/gin-gonic/gin"
)

func DeserializeUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// get token from cookie or header
		var accessToken string
		cookie, err := ctx.Cookie("access_token")

		authorizationHeader := ctx.Request.Header.Get("Authorization")

		fields := strings.Fields(authorizationHeader)
		// if auth header is not empty
		if len(fields) != 0 && fields[0] == "Bearer" {
			accessToken = fields[1]
			// else if cookie is not error empty
		} else if err == nil {
			accessToken = cookie
		}

		// send error: if there's no token from either cookie or auth header then send error response
		if accessToken == "" {
			ctx.AbortWithStatusJSON(http.StatusForbidden, dtos.WebResponse{
				Error:   true,
				Message: err.Error(),
			})
			return
		}

		// validate the token and get the user id from the sub/subject
		sub, err := utils.ValidateToken(accessToken, config.GlobalConfig.AccessTokenPublicKey)

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, dtos.WebResponse{
				Error:   true,
				Message: err.Error(),
			})
			return
		}

		var user models.User
		result := connection.DB.First(&user, "id=?", fmt.Sprint(sub))

		if result.Error != nil {
			ctx.AbortWithStatusJSON(http.StatusForbidden, dtos.WebResponse{
				Error:   true,
				Message: "The user belonging to this token doesn't exist anymore",
			})
			return
		}

		ctx.Set("currentUser", user)
		ctx.Next()

	}
}
