package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	router := gin.New()
	router.Use(gin.Recovery())
	router.GET("/ping", func(ctx *gin.Context) {
		ctx.String(200, "ok")
	})
	if err := router.Run(":8080"); err != nil {
		logrus.Fatalf("Cannot start server with err: %v", err)
	}
}