package server

import (
	"github.com/gin-gonic/gin"
	"github.com/majori/bolley/src/controllers"
	"github.com/majori/bolley/src/middlewares"
)

func NewRouter() *gin.Engine {
	r := gin.New()
	r.LoadHTMLGlob("src/templates/*")
	r.Static("/assets", "assets")

	publicRoutes(r.Group("/"))
	apiRoutes(r.Group("/api/v1"))

	return r
}

func publicRoutes(group *gin.RouterGroup) {
	group.GET("/", controllers.Teams)
	group.GET("/teams/:team-name", controllers.TeamStats)
}

func apiRoutes(group *gin.RouterGroup) {
	group.Use(middlewares.TokenAuth())
	group.POST("/upload", controllers.Upload)
}
