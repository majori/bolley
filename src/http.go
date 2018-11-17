package main

import (
	"bytes"
	"io/ioutil"
	"os"

	"github.com/gin-gonic/gin"
)

func listen() {
	r := gin.Default()
	r.POST("/upload", upload)
	r.Run("0.0.0.0:" + os.Getenv("PORT")) // listen and serve on 0.0.0.0:8080
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
