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

var matchModel = new(models.Match)
var teamModel = new(models.Team)

// upload logic
func Upload(c *gin.Context) {
	buff, err := ioutil.ReadAll(c.Request.Body)
	defer c.Request.Body.Close()
	if err != nil {
		panic(err)
	}

	match, err := parser.Parse(bytes.NewReader(buff))
	if (err != nil) {
		panic(err)
	}

	err = matchModel.Create(match)
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
	teams := teamModel.GetAll()
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"Title": "Main website",
		"Teams": teams,
	})
}

func TeamStats(c *gin.Context) {
	stats := teamModel.GetCumulativeStats("LP Kang 1")
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
