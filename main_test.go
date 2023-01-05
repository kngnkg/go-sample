package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPingRoute(t *testing.T) {
	t.Skip("リファクタリング中")
	router := setupRouter()

	got := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/ping", nil)
	if err != nil {
		t.Fatalf("cannot create request: %v", err)
	}
	router.ServeHTTP(got, req)

	assert.Equal(t, 200, got.Code)
	assert.Equal(t, "pong", got.Body.String())
}

func TestRun(t *testing.T) {
	t.Skip("リファクタリング中")
}
