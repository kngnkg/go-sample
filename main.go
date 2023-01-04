package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

// go run . {任意のポート番号}
func main() {
	if len(os.Args) != 2 {
		log.Printf("need port number\n")
		os.Exit(1)
	}
	r := setupRouter()
	r.Run(":" + os.Args[1])
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	return r
}
