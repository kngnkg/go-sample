package handler

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kwtryo/go-sample/model"
	"github.com/kwtryo/go-sample/store"
)

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
	userName := c.Query("user_name")
	u, err := uh.Service.GetUser(c.Request.Context(), userName)
	if err != nil {
		log.Printf("err: %v", err)
		if err == store.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"msg": "ユーザーが見つかりません。"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		}
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, u)
}
