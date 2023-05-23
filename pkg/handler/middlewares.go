package handler

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// logFormatter is a custom log formatter
func (h *Handler) logFormatter(param gin.LogFormatterParams) string {
	return fmt.Sprintf("%s From(%s) at:[%s] %s %d %s \" %s\"\n",
		param.Method,
		param.ClientIP,
		param.TimeStamp.Format(time.RFC1123),
		param.Path,
		param.StatusCode,
		param.Latency,
		param.ErrorMessage,
	)
}

// Logger is a custom logger
func (h *Handler) Logger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(h.logFormatter)
}

// AuthMiddleware is a custom auth middleware
func (h *Handler) AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//recieves an Authorization header from the request
		header := ctx.GetHeader("Authorization")
		if header == "" {
			ctx.AbortWithError(401, errors.New("authorization header is empty"))
			return
		}
		//splits the header into parts
		headerParts := strings.Split(header, " ")

		//checks if the parts are of the correct type
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			ctx.AbortWithError(401, errors.New("authorization header did not provide a token"))
			return
		}
		username, err := h.services.ParseToken(headerParts[1])

		if err != nil {
			ctx.AbortWithError(401, err)
			return
		}
		//sets the user in the gin Engine context
		user, _ := h.services.Api.GetUserByUsername(username)
		ctx.Set("user", user)

		ctx.Next()

	}
}
