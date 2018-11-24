package server

import (
	"os"
)

func Init() {
	r := NewRouter()
	r.Run(":" + os.Getenv("PORT"))
}
