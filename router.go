package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kwtryo/go-sample/handler"
)

func setupRouter(ctx context.Context) *gin.Engine {
	router := gin.Default()
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	router.GET("/users", handler.GetAllUser)
	return router
}
