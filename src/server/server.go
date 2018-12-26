package server

import (
	"os"
)

func Init() {
	r := NewRouter()
	var domain string
	if os.Getenv("APP_ENV") != "production" {
		domain = "localhost"
	}
	r.Run(domain + ":" + os.Getenv("PORT"))
}
