package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kwtryo/go-sample/handler"
)

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	router.GET("/users", handler.GetAllUser)
	return router
}
