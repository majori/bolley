package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/majori/bolley/controllers"
	"github.com/majori/bolley/middleware"
)

func NewRouter() *gin.Engine {
	r := gin.New()
	r.LoadHTMLGlob("templates/*")

	publicRoutes(r.Group("/"))
	apiRoutes(r.Group("/api/v1"))

	return r
}

func publicRoutes(group *gin.RouterGroup) {
	group.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Main website",
		})
	})
	group.GET("/ping", controllers.Pong)
	group.GET("/team", controllers.TeamStats)
}

func apiRoutes(group *gin.RouterGroup) {
	group.Use(middleware.TokenAuthMiddleware())
	group.POST("/upload", controllers.Upload)
}
