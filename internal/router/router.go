package router

import (
	"start/internal/config"
	"start/internal/router/api"
	"start/internal/router/middleware"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(cfg *config.Config) *gin.Engine {
	router := gin.Default()
	router.LoadHTMLGlob("static/*.html")
	router.Static("/static", "./static")
	registerMiddlewares(router, cfg)
	registerRoutes(router, cfg)

	return router
}

func registerMiddlewares(router *gin.Engine, cfg *config.Config) {
	router.Use(middleware.BasicAuthMiddleware(cfg))
}

func registerRoutes(router *gin.Engine, cfg *config.Config) {

	router.POST("/upload", api.UploadFile)

	router.POST("/deploy", func(c *gin.Context) {
		api.DeployStack(c, &cfg.TerminalConfig)
	})

	router.POST("/stop", func(c *gin.Context) {
		api.StopStack(c, &cfg.TerminalConfig)
	})

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
