package controllers

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	
	"github.com/majori/bolley/parser"
	"github.com/majori/bolley/models"
	"github.com/gin-gonic/gin"
)

// upload logic
func Upload(c *gin.Context) {
	buff, err := ioutil.ReadAll(c.Request.Body)
	defer c.Request.Body.Close()
	if err != nil {
		panic(err)
	}

	match, _ := parser.Parse(bytes.NewReader(buff))
	err = models.CreateMatch(match)
	if err != nil {
		switch err.Error() {
		case models.DBErrDublicate:
			RespondWithError(http.StatusConflict, "Match already exists", c)
		default:
			panic(err)
		}
		return
	}
	c.JSON(200, match)
}

func Teams(c *gin.Context) {

}

func TeamStats(c *gin.Context) {
	stats := models.GetTeamStats("LP Kang 1")
	fmt.Println(stats)
}

func Pong(c *gin.Context) {
	c.String(200, "pong")
}

func RespondWithError(code int, message string, c *gin.Context) {
	resp := map[string]string{"error": message}

	c.JSON(code, resp)
	c.Abort()
}
