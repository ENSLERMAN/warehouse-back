package service

import (
	"database/sql"
	"github.com/ENSLERMAN/warehouse-back/internal/handlers"
	"github.com/gin-gonic/gin"
)

var (
	Router *gin.Engine
	db     *sql.DB
)

func StartServer() {
	db = initDB()

	Router = gin.New()
	Router.Use(gin.Recovery())
	Router.Use(gin.Logger())
	Router.Use(cors())
	Router.GET("/ping", func(c *gin.Context) {
		c.String(200, "ping ok!")
	})
	nonAuth := Router.Group("/static")
	{
		nonAuth.POST("/register", handlers.Register(db))
		nonAuth.POST("/login", handlers.Login(db))
	}
	v1 := Router.Group("/api/v1", gin.BasicAuth(gin.Accounts{
		"admin": "develop",
	}))
	{
		v1.GET("/users", handlers.GetAllUsers(db))
		v1.GET("/users:id", handlers.GetUserByID)
	}
}

func cors() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		ctx.Next()
	}
}
