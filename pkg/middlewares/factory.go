package middlewares

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	AuthorizationHeaderKey  = "authorization"
	AuthorizationType       = "bearer"
	AuthorizationPayloadKey = "user"
)

func NewAuthMiddleware(tokenMaker Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		
		cookie, err := ctx.Cookie("session_cookie")
		 
		fmt.Println(err, "---")
		if err  != nil {
			err := errors.New("cookie is not provided")
			ctx.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		accessToken := cookie
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			err := errors.New("invalid token")
			ctx.AbortWithError(http.StatusUnauthorized, err)
			return
		}
		ctx.Set(AuthorizationPayloadKey, payload)
		ctx.Next()

		// authorizationHeader := ctx.GetHeader(AuthorizationHeaderKey)
		// if authorizationHeader == "" {
		// 	err := errors.New("authorization header is not provided")
		// 	ctx.AbortWithError(http.StatusUnauthorized, err)
		// 	return
		// }
		// fields := strings.Fields(authorizationHeader)
		// if len(fields) != 2 {
		// 	err := errors.New("invalid authorization header format")
		// 	ctx.AbortWithError(http.StatusUnauthorized, err)
		// 	return
		// }
		// authType := strings.ToLower(fields[0])
		// if authType != AuthorizationType {
		// 	err := fmt.Errorf("not support %v token type", authType)
		// 	ctx.AbortWithError(http.StatusUnauthorized, err)
		// 	return
		// }
		// accessToken := fields[1]
		// payload, err := tokenMaker.VerifyToken(accessToken)
		// if err != nil {
		// 	err := errors.New("invalid token")
		// 	ctx.AbortWithError(http.StatusUnauthorized, err)
		// 	return
		// }
		// ctx.Set(AuthorizationPayloadKey, payload)
		// ctx.Next()
	}
}