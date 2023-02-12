package handler

import (
	_ "embed"
	"net/http"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/kwtryo/go-sample/clock"
)

//go:embed cert/secret.pem
var rawPrivKey []byte

//go:embed cert/public.pem
var rawPubKey []byte

const (
	// 署名に用いるアルゴリズム
	SigningAlgorithm = "RS256"
)

type AuthService interface {
	Authenticator(c *gin.Context) (interface{}, error)
	PayloadFunc(data interface{}) jwt.MapClaims
	IdentityHandler(c *gin.Context) interface{}
	Authorizator(data interface{}, c *gin.Context) bool
	Logout(c *gin.Context) error
	IsAdmin(c *gin.Context) bool
}

type JWTer struct {
	Service AuthService
	Clocker clock.Clocker
}

func (j *JWTer) NewJWTMiddleware() (*jwt.GinJWTMiddleware, error) {
	return jwt.New(&jwt.GinJWTMiddleware{
		Realm: "test zone",

		// Key: rawPrivKey,
		PrivKeyBytes:     rawPrivKey,
		PubKeyBytes:      rawPubKey,
		SigningAlgorithm: SigningAlgorithm,

		// トークンの生存期間
		Timeout: 30 * time.Minute,

		MaxRefresh: 30 * time.Minute,
		// IdentityKey: UserNameKey,

		// ペイロードの独自クレーム設定
		// LoginHandlerで呼ばれる
		PayloadFunc: j.Service.PayloadFunc,

		// クレームからユーザー情報を取得する
		IdentityHandler: j.Service.IdentityHandler,

		// ログイン認証のための関数
		// LoginHandlerで呼ばれる
		Authenticator: j.Service.Authenticator,

		//トークンのユーザ情報からの認証
		Authorizator: j.Service.Authorizator,

		LogoutResponse: j.logoutResponse,

		Unauthorized: unauthorized,

		// "<source>:<name>"形式の文字列
		// リクエストからトークンを抽出するために使用される
		TokenLookup: "header: Authorization, query: token, cookie: jwt",

		// ヘッダの文字列
		TokenHeadName: "Bearer",

		// "orig_iat"
		// 現在の時間(トークンが生成された時間)
		TimeFunc: j.Clocker.Now,
	})
}

// 管理者のみ
func (j *JWTer) AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !j.Service.IsAdmin(c) {
			unauthorized(c, http.StatusUnauthorized, "許可されていません。")
			c.Abort()
		}
		c.Next()
	}
}

func (j *JWTer) logoutResponse(c *gin.Context, code int) {
	if err := j.Service.Logout(c); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
	})
}

func unauthorized(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"code":    code,
		"message": message,
	})
}
