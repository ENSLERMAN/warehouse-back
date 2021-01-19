package service

import (
	"database/sql"
	"github.com/ENSLERMAN/warehouse-back/internal/handlers"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func StartServer() *gin.Engine {
	db := initDB()
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())
	r.Use(cors())

	accs := initBasicAuthLogins(db)

	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "ping ok!")
	})
	nonAuth := r.Group("/auth")
	{
		nonAuth.POST("/register", handlers.Register(db))
		nonAuth.POST("/login", handlers.Login(db))
	}
	users := r.Group("/api/user", basicAuth(accs))
	{
		users.GET("/users", handlers.GetAllUsers(db))
		users.GET("/users:id", handlers.GetUserByID)
		users.GET("/me", handlers.ShowMe(db))
		users.POST("/update_role", handlers.UpdateRole(db))
	}
	shipments := r.Group("/api/shipments", basicAuth(accs))
	{
		shipments.POST("/new_shipment", handlers.AddNewShipment(db))
	}
	dispatch := r.Group("/api/dispatch", basicAuth(accs))
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

func initBasicAuthLogins(db *sql.DB) gin.Accounts {
	accs := make(map[string]string)
	result, err := db.Query("select login, password from warehouse.users")
	if err != nil {
		logrus.Fatalf("cannot get users with err: %v", err.Error())
		return nil
	}
	var accounts struct {
		login    string
		password string
	}
	for result.Next() {
		err := result.Scan(&accounts.login, &accounts.password)
		if err != nil {
			logrus.Fatal("cannot get user " + err.Error())
			return nil
		}
		accs[accounts.login] = accounts.password
	}

	if err = result.Err(); err != nil {
		logrus.Error(err.Error())
		return nil
	}

	return accs
}
