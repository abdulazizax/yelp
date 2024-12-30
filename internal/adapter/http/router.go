// Package api API.
//
// @title # UdevsLab Homework3
// @version 1.03.67.83.145
//
// @description API Endpoints for MiniTwitter
// @termsOfService http://swagger.io/terms/
//
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
//
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
//
// @host localhost:8080
// @BasePath /
// @schemes http https
//
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package http

import (
	_ "github.com/abdulazizax/yelp/api/docs"
	"github.com/abdulazizax/yelp/config"
	handler "github.com/abdulazizax/yelp/internal/adapter/http/handlers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Roter(handler handler.Handlers, config *config.Config) error {
	router := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowCredentials = true
	corsConfig.AllowHeaders = []string{"*"}
	corsConfig.AllowBrowserExtensions = true
	corsConfig.AllowMethods = []string{"*"}
	router.Use(cors.New(corsConfig))

	url := ginSwagger.URL("/swagger/doc.json")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url, ginSwagger.PersistAuthorization(true)))

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Handlers
	router.POST("/sign-up", handler.UserHandler().SignUp)
	router.POST("/sign-in", handler.UserHandler().SignIn)
	router.POST("send-verification-code", handler.UserHandler().SendVerificationCode)
	router.POST("/update-password", handler.UserHandler().UpdateUserPassword)

	return router.Run(config.Server.Port)
}
