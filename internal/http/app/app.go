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
package app

import (
	"log/slog"

	_ "github.com/abdulazizax/yelp/internal/http/app/docs"
	"github.com/abdulazizax/yelp/internal/http/handler"
	"github.com/abdulazizax/yelp/pkg/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Run(handler *handler.Handler, logger *slog.Logger, config *config.Config) error {
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

	return router.Run(config.Server.Port)
}
