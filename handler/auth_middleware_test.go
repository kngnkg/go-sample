package handler

import (
	"bytes"
	_ "embed"
	"errors"
	"io"
	"net/http"
	"testing"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/kwtryo/go-sample/clock"
	"github.com/kwtryo/go-sample/model"
	"github.com/kwtryo/go-sample/service"
	"github.com/kwtryo/go-sample/testutil"
	"github.com/kwtryo/go-sample/testutil/fixture"
)

func TestEmbed(t *testing.T) {
	want := []byte("-----BEGIN PUBLIC KEY-----")
	if !bytes.Contains(rawPubKey, want) {
		t.Errorf("want %s but got %s", want, rawPubKey)
	}

	want = []byte("-----BEGIN RSA PRIVATE KEY-----")
	if !bytes.Contains(rawPrivKey, want) {
		t.Errorf("want %s but got %s", want, rawPrivKey)
	}
}

// LoginHandlerで呼ばれる関数
// Authenticator
// PayloadFunc
func TestLoginRoute(t *testing.T) {
	// 正常値
	const (
		UserName    = "testUserName"
		RawPassword = "testPassword"
		uuid        = "8c2b2ad4-9598-44ed-b1b9-9ec355332f60"
	)
	type mock struct {
		user  *model.User
		token jwt.MapClaims
		err   error
	}
	user := fixture.User(&model.User{
		UserName: UserName,
		Password: RawPassword,
	})

	tok := jwt.MapClaims{
		service.JwtIdKey:    uuid,                          // ユニークID
		service.IssuerKey:   "github.com/kwtryo/go-sample", // 発行者
		service.SubjectKey:  "access_token",                // 用途
		service.AudienceKey: "github.com/kwtryo/go-sample", // 想定利用者
		// 以下独自クレーム
		service.UserNameKey: user.UserName,
		service.RoleKey:     user.Role,
	}

	tests := []struct {
		name         string
		mock         mock
		req          io.Reader // リクエスト
		wantStatus   int       // ステータスコード
		wantRespFile string    // レスポンス
	}{
		{
			"ok",
			mock{user, tok, nil},
			fixture.LoginFormBody(&model.Login{
				Username: UserName,
				Password: RawPassword,
			}),
			http.StatusOK,
			"testdata/login/ok_response.json.golden",
		},
		// リクエストが不正
		{
			"badRequest",
			mock{nil, jwt.MapClaims{}, errors.New("err from mock")},
			fixture.LoginFormBody(&model.Login{
				Username: "invalidUserName",
				Password: RawPassword,
			}),
			http.StatusUnauthorized,
			"testdata/login/bad_req_response.json.golden",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			serviceMoq := &AuthServiceMock{}
			serviceMoq.AuthenticatorFunc = func(c *gin.Context) (interface{}, error) {
				return user, tt.mock.err
			}
			serviceMoq.PayloadFuncFunc = func(data interface{}) jwt.MapClaims {
				return tt.mock.token
			}
			j := &JWTer{
				Service: serviceMoq,
				Clocker: clock.FixedClocker{},
			}
			got, err := j.NewJWTMiddleware()
			if err != nil {
				t.Errorf("JWTer.NewJWTMiddleware() error = %v", err)
				return
			}

			testutil.CheckHandlerFunc(
				t,
				got.LoginHandler,
				"POST",
				"",
				tt.req,
				tt.wantStatus,
				tt.wantRespFile,
			)
		})
	}
}

func TestLogoutRoute(t *testing.T) {
	// 正常値
	const (
		UserName = "testUserName"
	)
	serviceMoq := &AuthServiceMock{}
	serviceMoq.LogoutFunc = func(c *gin.Context) error {
		return nil
	}
	tests := []struct {
		name         string
		user         *model.User
		wantStatus   int    // ステータスコード
		wantRespFile string // レスポンス
	}{
		{
			"ok",
			fixture.User(&model.User{
				UserName: UserName,
			}),
			http.StatusOK,
			"testdata/logout/ok_response.json.golden",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &JWTer{
				Service: serviceMoq,
				Clocker: clock.FixedClocker{},
			}
			authMiddleware, err := j.NewJWTMiddleware()
			if err != nil {
				t.Errorf("JWTer.NewJWTMiddleware() error = %v", err)
				return
			}

			testutil.CheckHandlerFunc(
				t,
				authMiddleware.LogoutHandler,
				"GET",
				"",
				nil,
				tt.wantStatus,
				tt.wantRespFile,
			)
		})
	}
}

func TestAdminMiddleware(t *testing.T) {
	tests := []struct {
		name         string
		wantStatus   int    // ステータスコード
		wantRespFile string // レスポンス
	}{
		{
			"ok",
			http.StatusOK,
			"testdata/admin/ok_response.json.golden",
		},
		{
			"unAuthorized",
			http.StatusUnauthorized,
			"testdata/admin/unauthorized_response.json.golden",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			serviceMoq := &AuthServiceMock{}
			serviceMoq.IsAdminFunc = func(c *gin.Context) bool {
				return tt.name == "ok"
			}
			j := &JWTer{
				Service: serviceMoq,
				Clocker: clock.FixedClocker{},
			}
			got := j.AdminMiddleware()

			testutil.CheckMiddleware(
				t,
				got,
				tt.wantStatus,
				tt.wantRespFile,
			)
		})
	}
}
