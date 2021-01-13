package handlers

import (
	"github.com/gin-gonic/gin"
)

func GetAllUsers(c *gin.Context) {
	c.String(200, "get all users")
}

func GetUserByID(c *gin.Context) {
	c.String(200, "get user by id")
}