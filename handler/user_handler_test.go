package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestUsersRoute(t *testing.T) {
	r := gin.Default()
	r.GET("/users", GetAllUser)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "Leanne Graham")
}

func TestRegisterUserRoute(t *testing.T) {
	r := gin.Default()
	r.POST("/register", RegisterUser)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/register", "リクエストボディ")
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "登録するユーザのID")
}
