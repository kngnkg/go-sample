package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kwtryo/go-sample/clock"
	"github.com/kwtryo/go-sample/config"
	"github.com/kwtryo/go-sample/handler"
	"github.com/kwtryo/go-sample/service"
	"github.com/kwtryo/go-sample/store"
)

func SetupRouter(cfg *config.Config) (*gin.Engine, func(), error) {
	// TODO: database.init()に分離して、mainからDIしたい
	db, cleanup, err := store.New(cfg)
	if err != nil {
		return nil, cleanup, err
	}
	r := &store.Repository{Clocker: clock.RealClocker{}}

	userHandler := &handler.UserHandler{
		Service: &service.UserService{DB: db, Repo: r},
	}

	router := gin.Default()
	// ヘルスチェック
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	router.POST("/register", userHandler.RegisterUser)

	return router, cleanup, nil
}
