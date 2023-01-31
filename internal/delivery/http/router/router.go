package router

import (
	"api/config"
	"api/docs"
	"api/internal/delivery/http/handler"
	"api/pkg/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"       // swagger embed files
	"github.com/swaggo/gin-swagger" // gin-swagger middleware
	"go.uber.org/zap"
	"time"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3011")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func SetupRoutes(router *gin.RouterGroup, handler handler.AppHandler) {

	docs.SwaggerInfo.Title = "Database Backup API"
	docs.SwaggerInfo.Description = "This is database backup API."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "x.swagger.io"
	docs.SwaggerInfo.BasePath = "/v1"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	tokenMaker, _ := middlewares.NewPasetoMaker(config.C.Token)
	authMiddleware := middlewares.NewAuthMiddleware(tokenMaker)

	m := router.Group("/api/v1")

	m.POST("/login", handler.AuthHandler.Login)
	m.POST("/register", handler.AuthHandler.Register)
	m.GET("/confirm/:token", handler.AuthHandler.Confirm)
	m.POST("/reset-link", handler.AuthHandler.ResetLink)
	m.POST("/reset-password", handler.AuthHandler.ResetPassword)
	m.GET("/verify", handler.AuthHandler.VerifyToken)
	m.POST("/logout", handler.AuthHandler.Logout)
	m.POST("/verify/tf", handler.AuthHandler.VerifyTF)


	//USER API
	u := router.Group("/api/v1/users").Use(authMiddleware)
	u.GET("", handler.UserHandler.GetAllUsers)
	u.PUT("/:id", handler.UserHandler.UpdateUser)
	u.PUT("/current-user", handler.UserHandler.UpdateCurrentUser)
	u.PUT("/current-user/password", handler.UserHandler.UpdateCurrentUserPassword)
	u.GET("/current-user", handler.UserHandler.GetCurrentUser)
	u.DELETE("/collection", handler.UserHandler.DeleteUsers)
	u.GET("/current-user/qr", handler.UserHandler.GetQr)
	u.POST("/current-user/qr/verify", handler.UserHandler.VerifyTF)
	u.PUT("/current-user/qr/disable", handler.UserHandler.DisableTF)


	

}

func NewRouter(logger *zap.Logger, handler handler.AppHandler) *gin.Engine {

	docs.SwaggerInfo.Title = "Database Backup API"
	docs.SwaggerInfo.Description = "This is database backup API."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "x.swagger.io"
	docs.SwaggerInfo.BasePath = "/v1"

	router := gin.New()

	// Logs all requests, like a combined access and error log.
	// Logs to stdout.
	// RFC3339 with UTC time format.
	router.Use(middlewares.Ginzap(logger, time.RFC3339, true))

	// Logs all panic to error log
	router.Use(middlewares.RecoveryWithZap(logger, true))

	router.Use(CORSMiddleware())

	SetupRoutes(&router.RouterGroup, handler)

	return router
}
