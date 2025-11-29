package main

import (
	"task02/internal/config"
)

func main() {
	config.InitConfig("./etc/config.yaml")
}
