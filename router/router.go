package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kwtryo/go-sample/clock"
	"github.com/kwtryo/go-sample/config"
	"github.com/kwtryo/go-sample/handler"
	"github.com/kwtryo/go-sample/middleware"
	"github.com/kwtryo/go-sample/service"
	"github.com/kwtryo/go-sample/store"
)

func SetupRouter(cfg *config.Config) (*gin.Engine, func(), error) {
	db, cleanup, err := store.New(cfg)
	if err != nil {
		return nil, cleanup, err
	}
	r := &store.Repository{Clocker: clock.RealClocker{}}

	userHandler := &handler.UserHandler{
		Service: &service.UserService{DB: db, Repo: r},
	}

	router := gin.Default()
	router.Use(middleware.DBTransactionMiddleware(db))

	// ヘルスチェック
	router.GET("/health", func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		// c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		// c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	router.POST("/register", userHandler.RegisterUser)

	return router, cleanup, nil
}
