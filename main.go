package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/kwtryo/go-sample/config"
)

// go run . {任意のポート番号}
func main() {
	if err := run(); err != nil {
		log.Printf("failed to terminate server: %v", err)
		os.Exit(1)
	}
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	return r
}

func run() error {
	cfg, err := config.New()
	if err != nil {
		return err
	}

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	return r.Run(fmt.Sprintf(":%d", cfg.Port))
}
