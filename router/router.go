package router

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kwtryo/go-sample/clock"
	"github.com/kwtryo/go-sample/config"
	"github.com/kwtryo/go-sample/handler"
	"github.com/kwtryo/go-sample/service"
	"github.com/kwtryo/go-sample/store"
)

func SetupRouter(cfg *config.Config) (*gin.Engine, func(), error) {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		// 許可したいHTTPメソッドの一覧
		AllowMethods: []string{
			"POST",
			"GET",
			"OPTIONS",
			"PUT",
			"DELETE",
		},
		// 許可したいHTTPリクエストヘッダの一覧
		AllowHeaders: []string{
			"Access-Control-Allow-Headers",
			"Content-Type",
			"Content-Length",
			"Accept-Encoding",
			"X-CSRF-Token",
			"Authorization",
		},
		// 許可したいアクセス元の一覧
		AllowOrigins: []string{
			"http://localhost:3000",
		},
		// preflight requestで許可した後の接続可能時間
		MaxAge: 24 * time.Hour,
	}))

	db, cleanup, err := store.New(cfg)
	if err != nil {
		return nil, cleanup, err
	}
	r := &store.Repository{Clocker: clock.RealClocker{}}

	healthHandler := &handler.HealthHandler{
		Service: &service.HealthService{DB: db, Repo: r},
	}
	router.GET("/health", healthHandler.HealthCheck)

	userHandler := &handler.UserHandler{
		Service: &service.UserService{DB: db, Repo: r},
	}
	router.POST("/register", userHandler.RegisterUser)
	router.GET("/user", userHandler.GetUser)

	return router, cleanup, nil
}
