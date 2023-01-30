package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kwtryo/go-sample/model"
)

type HealthService interface {
	HealthCheck(ctx context.Context) error
}

type HealthHandler struct {
	Service HealthService
}

// GET /health
// ヘルスチェック
func (hh *HealthHandler) HealthCheck(c *gin.Context) {
	if err := hh.Service.HealthCheck(c.Request.Context()); err != nil {
		c.JSON(http.StatusInternalServerError, model.Health{
			Health:   model.StatusOrange,
			Database: model.StatusRed,
		})
		return
	}
	c.JSON(http.StatusOK, model.Health{
		Health:   model.StatusGreen,
		Database: model.StatusGreen,
	})
}
