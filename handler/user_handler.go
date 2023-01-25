package handler

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kwtryo/go-sample/model"
)

//go:generate go run github.com/matryer/moq -out moq_test.go . UserService
type UserService interface {
	RegisterUser(ctx context.Context, form *model.FormRequest) (*model.User, error)
	GetUser(ctx context.Context, userName string) (*model.User, error)
}

type UserHandler struct {
	Service UserService
}

// POST /register
// ユーザーを登録し、登録したユーザーのIDをレスポンスとして返す
func (uh *UserHandler) RegisterUser(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")

	form := &model.FormRequest{}
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		c.Abort()
		return
	}

	u, err := uh.Service.RegisterUser(c.Request.Context(), form)
	if err != nil {
		log.Printf("err: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": u.Id})
}

// GET /user?user_name=user_name
// ユーザー名からユーザーを取得し、レスポンスとして返す。
func (uh *UserHandler) GetUser(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")

	userName := c.Query("user_name")
	u, err := uh.Service.GetUser(c.Request.Context(), userName)
	if err != nil {
		log.Printf("err: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, u)
}
