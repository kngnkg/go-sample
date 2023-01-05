package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		// // graceful shutdownの確認
		// time.Sleep(5 * time.Second)
		c.String(http.StatusOK, "pong")
	})
	return router
}
