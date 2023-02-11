package handler

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kwtryo/go-sample/model"
	"github.com/kwtryo/go-sample/store"
)

const (
	BadRequestMessage       = "不正なリクエストです。"
	ServerErrorMessage      = "サーバー内部でエラーが発生しました。"
	UserNotFoundMessage     = "ユーザーが存在しません。"
	UserAlreadyEntryMessage = "既に登録されたユーザーです。"
)

type UserService interface {
	RegisterUser(ctx context.Context, form *model.FormRequest) (*model.User, error)
	GetAllUsers(ctx context.Context) (model.Users, error)
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
		log.Printf("ERROR: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"msg": BadRequestMessage})
		return
	}

	u, err := uh.Service.RegisterUser(c.Request.Context(), form)
	if err != nil {
		log.Printf("ERROR: %v", err)
		if errors.Is(err, store.ErrAlreadyEntry) {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": UserAlreadyEntryMessage})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"msg": ServerErrorMessage})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": u.Id})
}

// GET /auth/users
// 全てのユーザーを取得し、レスポンスとして返す
func (uh *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := uh.Service.GetAllUsers(c.Request.Context())
	if err != nil {
		log.Printf("ERROR: %v", err)
		if errors.Is(err, store.ErrNotFound) {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": UserNotFoundMessage})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"msg": ServerErrorMessage})
		return
	}

	c.JSON(http.StatusOK, users)
}

// GET /auth/user?user_name=user_name
// ユーザー名からユーザーを取得し、レスポンスとして返す。
func (uh *UserHandler) GetUser(c *gin.Context) {
	userName := c.Query("user_name")
	u, err := uh.Service.GetUser(c.Request.Context(), userName)
	if err != nil {
		log.Printf("ERROR: %v", err)
		if errors.Is(err, store.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"msg": UserNotFoundMessage})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"msg": ServerErrorMessage})
		return
	}

	c.JSON(http.StatusOK, u)
}
