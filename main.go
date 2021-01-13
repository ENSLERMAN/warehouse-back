package main

import (
	"github.com/ENSLERMAN/warehouse-back/internal/service"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"os"
)

func init() {
	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("No .env file found")
	}
}

func main() {
	port, exist := os.LookupEnv("PORT")
	if !exist {
		logrus.Fatalf("PORT not exist in .env")
	}
	service.StartServer()
	if err := service.Router.Run(port); err != nil {
		logrus.Fatalf("cannot start server: %v", err)
	}
}
