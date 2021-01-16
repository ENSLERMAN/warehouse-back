package service

import (
	"github.com/ENSLERMAN/warehouse-back/internal/handlers"
	"github.com/gin-gonic/gin"
)

func StartServer() *gin.Engine {
	db := initDB()
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())
	r.Use(cors())
	r.Use()

	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "ping ok!")
	})
	nonAuth := r.Group("/static")
	{
		nonAuth.POST("/register", handlers.Register(db))
		nonAuth.POST("/login", handlers.Login(db))
	}
	v1 := r.Group("/api", gin.BasicAuth(gin.Accounts{
		"admin": "develop",
	}))
	{
		v1.GET("/users", handlers.GetAllUsers(db))
		v1.GET("/users:id", handlers.GetUserByID)
		v1.GET("/me", handlers.ShowMe(db))
		v1.POST("/update_role", handlers.UpdateRole(db))
	}
	return r
}

func cors() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		ctx.Next()
	}
}
