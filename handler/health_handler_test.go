package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHealthHandler_HealthCheck(t *testing.T) {
	tests := []struct {
		name       string
		wantStatus int // ステータスコード
	}{
		{
			"ok",
			http.StatusOK,
		},
		{
			"internalServerError",
			http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Serviceのモックを生成
			moqService := &HealthServiceMock{
				HealthCheckFunc: func(ctx context.Context) error {
					if tt.name == "ok" {
						return nil
					}
					return errors.New("error-test")
				},
			}
			hh := &HealthHandler{
				Service: moqService,
			}

			router := gin.Default()
			router.GET("/health", hh.HealthCheck)
			testServer := httptest.NewServer(router) // サーバを立てる
			t.Cleanup(func() {
				testServer.Close()
			})

			// テストサーバーにリクエストを送信
			resp, err := http.Get(testServer.URL + "/health")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			defer resp.Body.Close()

			// 期待するステータスコードと一致するか確認する
			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}
