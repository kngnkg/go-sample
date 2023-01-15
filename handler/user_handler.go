package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kwtryo/go-sample/model"
)

// type RegisterUserService interface {
// 	// ユーザーを登録し、登録したユーザーのIDを返す
// 	RegisterUser(ctx context.Context) (int, error)
// }

//go:generate go run github.com/matryer/moq -out moq_test.go . UserService
type UserService interface {
	// ユーザーを登録し、登録したユーザーのIDを返す
	RegisterUser(ctx context.Context) (int, error)
}

type UserHandler struct {
	// DB   *sqlx.DB
	// Repo *store.Repository
	// RegisterUserService RegisterUserService
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

	// u := &model.User{
	// 	Name:     form.Name,
	// 	UserName: form.UserName,
	// 	Password: form.Password,
	// 	Role:     form.Role,
	// 	Email:    form.Email,
	// 	Address:  form.Address,
	// 	Phone:    form.Phone,
	// 	Website:  form.Website,
	// 	Company:  form.Company,
	// }
	// err := uh.Repo.RegisterUser(c.Request.Context(), uh.DB, u)
	id, err := uh.Service.RegisterUser(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}
