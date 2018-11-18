package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func listen() {
	r := gin.New()

	// Public
	public := r.Group("/")
	public.GET("/ping", pong)
	public.GET("/team", teamStats)

	// Private
	private := r.Group("/")
	private.Use(tokenAuthMiddleware())
	private.POST("/upload", upload)

	r.Run(":" + os.Getenv("PORT"))
}

func respondWithError(code int, message string, c *gin.Context) {
	resp := map[string]string{"error": message}

	c.JSON(code, resp)
	c.Abort()
}

func tokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			respondWithError(401, "API token required", c)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") ||
			(parts[1] != os.Getenv("API_TOKEN")) {
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
	defer c.Request.Body.Close()
	if err != nil {
		panic(err)
	}
	match, _ := parseSpreadsheet(bytes.NewReader(buff))
	err = createMatch(match)
	if err != nil {
		switch err.Error() {
		case DBErrDublicate:
			respondWithError(http.StatusConflict, "Match already exists", c)
		default:
			panic(err)
		}
		return
	}
	c.JSON(200, match)
}

func teams(c *gin.Context) {

}

func teamStats(c *gin.Context) {
	stats := getTeamStats("LP Kang 1")
	fmt.Println(stats)
}
