package router

import (
	"context"
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

	db, cleanup, err := store.New(cfg)
	if err != nil {
		return nil, cleanup, err
	}
	kvs, err := store.NewKVS(context.TODO(), cfg)
	if err != nil {
		return nil, cleanup, err
	}

	cl := clock.RealClocker{}
	r := &store.Repository{Clocker: cl}

	j := &handler.JWTer{
		Service: &service.AuthService{DB: db, Repo: r, Store: kvs},
		Clocker: cl,
	}
	healthHandler := &handler.HealthHandler{
		Service: &service.HealthService{DB: db, Repo: r},
	}
	userHandler := &handler.UserHandler{
		Service: &service.UserService{DB: db, Repo: r},
	}

	authMiddleware, err := j.NewJWTMiddleware()
	if err != nil {
		return nil, nil, err
	}

	router.Use(GetCorsMiddleware())
	router.GET("/health", healthHandler.HealthCheck)
	router.POST("/register", userHandler.RegisterUser)
	router.POST("/login", authMiddleware.LoginHandler)
	auth := router.Group("/auth")
	auth.GET("/refresh_token", authMiddleware.RefreshHandler)
	auth.Use(authMiddleware.MiddlewareFunc())
	{
		auth.GET("/logout", authMiddleware.LogoutHandler)
		auth.GET("/user", userHandler.GetUser)
	}

	return router, cleanup, nil
}

// CORSの設定
func GetCorsMiddleware() gin.HandlerFunc {
	return cors.New(cors.Config{
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
	})
}
