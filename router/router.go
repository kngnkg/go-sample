package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kwtryo/go-sample/clock"
	"github.com/kwtryo/go-sample/config"
	"github.com/kwtryo/go-sample/handler"
	"github.com/kwtryo/go-sample/store"
)

func SetupRouter(cfg *config.Config) (*gin.Engine, func(), error) {
	clocker := clock.RealClocker{}
	r := store.Repository{Clocker: clocker}
	db, cleanup, err := store.New(cfg)
	if err != nil {
		return nil, cleanup, err
	}

	router := gin.Default()
	// ヘルスチェック
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// テスト用
	router.GET("/users", handler.GetAllUser)

	// ru := &handler.RegisterUser{
	// 	DB:   db,
	// 	Repo: &r,
	// }
	// router.POST("/register", ru.ServeHTTP)

	return router, cleanup, nil
}
