package main

import (
	"bytes"
	"io/ioutil"
	"os"

	"github.com/gin-gonic/gin"
)

func listen() {
	r := gin.New()

	// Public
	public := r.Group("/")
	public.GET("/ping", pong)

	// Private
	private := r.Group("/")
	private.Use(TokenAuthMiddleware())
	private.POST("/upload", upload)
	r.Run(":" + os.Getenv("PORT")) // listen and serve on 0.0.0.0:8080
}

func respondWithError(code int, message string, c *gin.Context) {
	resp := map[string]string{"error": message}

	c.JSON(code, resp)
	c.Abort()
}

func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.FormValue("api_token")

		if token == "" {
			respondWithError(401, "API token required", c)
			return
		}

		if token != os.Getenv("API_TOKEN") {
			respondWithError(401, "Invalid API token", c)
			return
		}

		c.Next()
	}
}

func pong(c *gin.Context) {
	c.String(200, "pong")
}

// upload logic
func upload(c *gin.Context) {
	buff, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		panic(err)
	}
	defer c.Request.Body.Close()
	match, _ := parseSpreadsheet(bytes.NewReader(buff))
	err = saveMatch(match)
	if err != nil {
		c.String(400, err.Error())
		return
	}
	c.String(200, "OK")
}
