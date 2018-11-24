package main

import (
	"github.com/majori/bolley/db"
	"github.com/majori/bolley/server"
)

func main() {
	db.Init()
	server.Init()
}
