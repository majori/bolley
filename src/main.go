package main

import (
	"github.com/majori/bolley/src/db"
	"github.com/majori/bolley/src/server"
)

func main() {
	db.Init()
	server.Init()
}
