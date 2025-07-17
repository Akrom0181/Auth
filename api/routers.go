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

	// Custom auth middleware
	authMiddleware := h.AuthorizerMiddleware() // Casbin emas, faqat JWT tekshiradigan

	// Protected routes
	protected := r.Group("/")

	protected.POST("/auth/login", h.Login)
	protected.POST("/auth/send-otp", h.SendOTP)
	protected.POST("/auth/confirm-otp", h.ConfirmOTP)
	protected.POST("/auth/signup", h.Signup)
	protected.Use(authMiddleware)

	protected.POST("/role", h.CreateRole)
	protected.GET("/role/list", h.GetListRoles)
	protected.GET("/role/:id", h.GetSingleRole)
	protected.PUT("/role/:id", h.UpdateRole)

	protected.POST("/sysuser", h.CreateSysuser)

	return r
}
