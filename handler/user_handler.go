package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kwtryo/go-sample/model"
)

//go:generate go run github.com/matryer/moq -out moq_test.go . UserService
type UserService interface {
	RegisterUser(ctx context.Context, form *model.FormRequest) (*model.User, error)
}

type UserHandler struct {
	Service UserService
}

// POST /register
// ユーザーを登録し、登録したユーザーのIDをレスポンスとして返す
func (uh *UserHandler) RegisterUser(c *gin.Context) {
	form := &model.FormRequest{}
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		c.Abort()
		return
	}

	u, err := uh.Service.RegisterUser(c.Request.Context(), form)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": u.Id})
}
