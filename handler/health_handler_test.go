package handler

import (
	"context"
	"errors"
	"net/http"
	"testing"

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

			testutil.CheckHandlerFunc(
				t,
				hh.HealthCheck,
				"GET",
				"",
				nil,
				tt.wantStatus,
				tt.wantRespFile,
			)
		})
	}
}
