package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Create creates the gin engine with all routes.
func Create() *gin.Engine {
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		// // graceful shutdownの確認
		// time.Sleep(5 * time.Second)
		c.String(http.StatusOK, "pong")
	})
	return router
}
