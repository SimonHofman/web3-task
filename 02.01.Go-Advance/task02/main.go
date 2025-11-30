package main

import (
	"task02/internal/config"
	"task02/pkg/client"
	"task02/pkg/db"
)

func main() {
	config.InitConfig("./etc/config.yaml")
	db.InitDB()
	client.InitClient()
}
