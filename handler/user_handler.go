package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kwtryo/go-sample/entity"
)

// GET /users
func GetAllUser(c *gin.Context) {
	users, _ := entity.AllUser()
	c.JSON(http.StatusOK, users)
}

// POST /register
func RegisterUser(c *gin.Context) {
}
