package api

import (
	_ "github.com/Akrom0181/Auth/api/docs"
	"github.com/Akrom0181/Auth/api/handler"
	"github.com/Akrom0181/Auth/config"
	"github.com/Akrom0181/Auth/pkg/logger"
	"github.com/Akrom0181/Auth/storage"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// NewRouter -.
// Swagger spec:
// @title       GO Auth API
// @description This is a sample server GO Auth API server.
// @version     1.0
// @host        localhost:8080
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func New(storage *storage.IStorage, log logger.ILogger, config config.Config) *gin.Engine {
	h := handler.NewStrg(*storage, log, config)
	r := gin.Default()

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Auth middleware
	authMiddleware := h.AuthorizerMiddleware()

	// Public routes
	r.POST("/auth/login", h.Login)
	r.POST("/auth/send-otp", h.SendOTP)
	r.POST("/auth/confirm-otp", h.ConfirmOTP)
	r.POST("/auth/signup", h.Signup)

	// Protected routes
	protected := r.Group("/")
	protected.Use(authMiddleware)

	// Role routes
	protected.POST("/role", h.CreateRole)
	protected.GET("/role/list", h.GetListRoles)
	protected.GET("/role/:id", h.GetSingleRole)
	protected.PUT("/role/:id", h.UpdateRole)

	// Sysuser route
	protected.POST("/sysuser", h.CreateSysuser)

	return r
}
