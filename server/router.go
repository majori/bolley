package server

import (
	"github.com/gin-gonic/gin"
	"github.com/majori/bolley/controllers"
	"github.com/majori/bolley/middlewares"
)

func NewRouter() *gin.Engine {
	r := gin.New()
	r.LoadHTMLGlob("templates/*")

	publicRoutes(r.Group("/"))
	apiRoutes(r.Group("/api/v1"))

	return r
}

func publicRoutes(group *gin.RouterGroup) {
	group.GET("/", controllers.Teams)
	group.GET("/ping", controllers.Pong)
	group.GET("/team", controllers.TeamStats)
}

func apiRoutes(group *gin.RouterGroup) {
	group.Use(middlewares.TokenAuth())
	group.POST("/upload", controllers.Upload)
}
