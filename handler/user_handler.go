package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/kwtryo/go-sample/model"
	"github.com/kwtryo/go-sample/store"
)

// GET /users
func GetAllUser(c *gin.Context) {
	users, _ := model.AllUser()
	c.JSON(http.StatusOK, users)
}

type RegisterUser struct {
	DB   *sqlx.DB
	Repo *store.Repository
}

// POST /register
// ユーザーを登録し、登録したユーザーのIDをレスポンスとして返す
func (ru *RegisterUser) ServeHTTP(c *gin.Context) {
	ctx := c.Request.Context()

	form := &model.FormRequest{}
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		c.Abort()
		return
	}

	u := &model.User{
		Name:     form.Name,
		UserName: form.UserName,
		Password: form.Password,
		Role:     form.Role,
		Email:    form.Email,
		Address:  form.Address,
		Phone:    form.Phone,
		Website:  form.Website,
		Company:  form.Company,
	}
	err := ru.Repo.RegisterUser(ctx, ru.DB, u)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": u.Id})
}
