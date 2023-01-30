package router

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kwtryo/go-sample/config"
	"github.com/kwtryo/go-sample/testutil"
)

func TestSetupRouter(t *testing.T) {
	type args struct {
		cfg *config.Config
	}
	tests := []struct {
		name         string
		args         args
		wantStatus   int
		wantRespFile string
		wantErr      bool
	}{
		{
			"ok",
			args{testutil.CreateConfigForTest(t)},
			http.StatusOK,
			"testdata/ok_response.json.golden",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotr, gotf, err := SetupRouter(tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetupRouter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			testServer := httptest.NewServer(gotr) // サーバを立てる
			t.Cleanup(func() {
				gotf()
				testServer.Close()
			})

			url := fmt.Sprintf(testServer.URL + "/health")
			t.Logf("try request to %q", url)
			// サーバーにリクエストを送信
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
