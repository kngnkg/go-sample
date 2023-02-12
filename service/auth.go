package service

import (
	"context"
	"fmt"
	"log"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kwtryo/go-sample/model"
	"github.com/kwtryo/go-sample/store"
	"golang.org/x/crypto/bcrypt"
)

type AuthRepository interface {
	GetUserByUserName(ctx context.Context, db store.DBConnection, userName string) (*model.User, error)
}

type Store interface {
	Save(ctx context.Context, key string, uid int) error
	Load(ctx context.Context, key string) (int, error)
	Delete(ctx context.Context, key string) error
}

type AuthService struct {
	DB    store.DBConnection
	Repo  AuthRepository
	Store Store
}

const (
	JwtIdKey    = "jti"
	IssuerKey   = "iss"
	SubjectKey  = "sub"
	AudienceKey = "aud"
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

// クレームを設定する
func (as *AuthService) PayloadFunc(data interface{}) jwt.MapClaims {
	v, ok := data.(*model.User)
	if !ok {
		return jwt.MapClaims{}
	}

	tok := jwt.MapClaims{
		JwtIdKey:    uuid.New().String(),           // ユニークID
		IssuerKey:   "github.com/kwtryo/go-sample", // 発行者
		SubjectKey:  "access_token",                // 用途
		AudienceKey: "github.com/kwtryo/go-sample", // 想定利用者
		// 以下独自クレーム
		UserNameKey: v.UserName,
		RoleKey:     v.Role,
	}

	if err := as.Store.Save(context.TODO(), tok[JwtIdKey].(string), v.Id); err != nil {
		log.Printf("cannot save jti: %v", err)
		return jwt.MapClaims{}
	}
	return tok
}

// クレームからユーザー情報を取得する
func (as *AuthService) IdentityHandler(c *gin.Context) interface{} {
	claims := jwt.ExtractClaims(c)
	// kvsに存在するか確認
	if _, err := as.Store.Load(c.Request.Context(), claims[JwtIdKey].(string)); err != nil {
		log.Printf("cannot load jti: %v", err)
		return nil
	}
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
	if err != nil {
		log.Printf("cannot get user: %v", err)
		return false
	}
	return true
}

// ログアウト処理
func (as *AuthService) Logout(c *gin.Context) error {
	// JWTトークンを無効化する
	claims := jwt.ExtractClaims(c)
	return as.Store.Delete(c.Request.Context(), claims[JwtIdKey].(string))
}

// ログインユーザーが管理者かどうか判定する
func (as *AuthService) IsAdmin(c *gin.Context) bool {
	claims := jwt.ExtractClaims(c)
	role, ok := claims[RoleKey].(string)
	if !ok {
		return false
	}
	return role == "admin"
}
