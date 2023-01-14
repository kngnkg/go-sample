package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kwtryo/go-sample/config"
	"github.com/stretchr/testify/assert"
)

func TestSetupRouter(t *testing.T) {
	cfg, err := config.CreateForTest()
	if err != nil {
		t.Fatalf("cannot get config: %v", err)
	}
	r, cleanup, err := SetupRouter(cfg)
	if err != nil {
		t.Fatalf("cannot setup router: %v", err)
	}
	defer cleanup()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "ok")
}
