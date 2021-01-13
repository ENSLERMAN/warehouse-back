package service

import (
	"github.com/ENSLERMAN/warehouse-back/internal/handlers"
	"github.com/gin-gonic/gin"
)

var Router *gin.Engine

func StartServer() {
	Router = gin.New()
	Router.Use(gin.Recovery())
	Router.Use(cors())
	Router.GET("/ping", func(c *gin.Context) {
		c.String(200, "ping ok!")
	})
	v1 := Router.Group("/api/v1")
	{
		v1.GET("/users", handlers.GetAllUsers)
		v1.GET("/users:id", handlers.GetUserByID)
		v1.POST("/login", handlers.GetUserByID)
	}
}

func cors() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		ctx.Next()
	}
}
