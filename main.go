package main

import "github.com/ENSLERMAN/warehouse-back/internal/service"

func main() {
	service.StartServer()
	service.Router.Run(":8080")
}