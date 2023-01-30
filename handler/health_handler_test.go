package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kwtryo/go-sample/testutil"
)

func TestHealthHandler_HealthCheck(t *testing.T) {
	tests := []struct {
		name         string
		wantStatus   int    // ステータスコード
		wantRespFile string // レスポンス

	}{
		{
			"ok",
			http.StatusOK,
			"testdata/health_check/ok_response.json.golden",
		},
		{
			"internalServerError",
			http.StatusInternalServerError,
			"testdata/health_check/server_err_response.json.golden",
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

			url := fmt.Sprintf(testServer.URL + "/health")
			t.Logf("try request to %q", url)
			// テストサーバーにリクエストを送信
			resp, err := http.Get(url)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			defer resp.Body.Close()

			testutil.AssertResponse(
				t,
				resp,
				tt.wantStatus,
				testutil.LoadFile(t, tt.wantRespFile),
			)
		})
	}
}
