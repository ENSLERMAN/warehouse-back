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
	users := r.Group("/api/user", gin.BasicAuth(gin.Accounts{
		"admin": "develop",
	}))
	{
		users.GET("/users", handlers.GetAllUsers(db))
		users.GET("/users:id", handlers.GetUserByID)
		users.GET("/me", handlers.ShowMe(db))
		users.POST("/update_role", handlers.UpdateRole(db))
	}
	shipments := r.Group("/api/shipments", gin.BasicAuth(gin.Accounts{
		"admin": "develop",
	}))
	{
		shipments.POST("/new_shipment", handlers.AddNewShipment(db))
	}
	dispatch := r.Group("/api/dispatch", gin.BasicAuth(gin.Accounts{
		"admin": "develop",
	}))
	{
		dispatch.POST("/new_dispatch", handlers.AddNewDispatch(db))
		dispatch.POST("/close_dispatch", handlers.CloseDispatch(db))
	}
	return r
}

func cors() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		ctx.Next()
	}
}
