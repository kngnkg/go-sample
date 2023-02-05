package service

import (
	"context"
	"fmt"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/kwtryo/go-sample/model"
	"github.com/kwtryo/go-sample/store"
	"golang.org/x/crypto/bcrypt"
)

type AuthRepository interface {
	GetUserByUserName(ctx context.Context, db store.DBConnection, userName string) (*model.User, error)
}

type AuthService struct {
	DB   store.DBConnection
	Repo AuthRepository
}

const (
	UserNameKey = "user_name"
	RoleKey     = "role"
)

// ログイン認証
func (as *AuthService) Authenticator(c *gin.Context) (interface{}, error) {
	var loginVals model.Login
	if err := c.ShouldBind(&loginVals); err != nil {
		return nil, jwt.ErrMissingLoginValues
	}
	userName := loginVals.Username
	password := loginVals.Password

	user, err := as.Repo.GetUserByUserName(c.Request.Context(), as.DB, userName)
	if err != nil {
		return nil, fmt.Errorf("%v: %v", jwt.ErrFailedAuthentication, err)
	}

	if userName != user.UserName {
		return nil, fmt.Errorf("%v: %v", jwt.ErrFailedAuthentication, err)
	}
	// パスワードの確認
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, fmt.Errorf("%v: %v", jwt.ErrFailedAuthentication, err)
	}

	return user, nil
}

// 独自クレームを設定する
func (as *AuthService) PayloadFunc(data interface{}) jwt.MapClaims {
	if v, ok := data.(*model.User); ok {
		// ペイロードにロールとユーザー名をセット
		return jwt.MapClaims{
			UserNameKey: v.UserName,
			RoleKey:     v.Role,
		}
	}
	return jwt.MapClaims{}
}

// クレームからユーザー情報を取得する
func (as *AuthService) IdentityHandler(c *gin.Context) interface{} {
	claims := jwt.ExtractClaims(c)
	return &model.User{
		// クレームからログインIDを取得する
		UserName: claims[UserNameKey].(string),
		Role:     claims[RoleKey].(string),
	}
}

// トークンのユーザ情報からの認証
func (as *AuthService) Authorizator(data interface{}, c *gin.Context) bool {
	v, ok := data.(*model.User)
	if !ok {
		return false
	}
	_, err := as.Repo.GetUserByUserName(c.Request.Context(), as.DB, v.UserName)
	// ユーザーが取得できたらtrueを返す
	return err == nil
}
