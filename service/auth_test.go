package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kwtryo/go-sample/model"
	"github.com/kwtryo/go-sample/store"
	"github.com/kwtryo/go-sample/testutil/fixture"
	"github.com/stretchr/testify/assert"
)

// 正常値
const (
	UserName    = "testUserName"
	RawPassword = "testPassword"
	Role        = "admin"
)

func TestAuthService_Authenticator(t *testing.T) {
	type fields struct {
		DB   store.DBConnection
		Repo AuthRepository
	}

	validUser := fixture.User(&model.User{
		UserName: UserName,
		Password: RawPassword,
	})
	moqDb := &DBConnectionMock{}
	moqRepo := &AuthRepositoryMock{}
	moqRepo.GetUserByUserNameFunc = func(ctx context.Context, db store.DBConnection, userName string) (*model.User, error) {
		if userName == UserName {
			return validUser, nil
		}
		return nil, fmt.Errorf("error from mock")
	}

	tests := []struct {
		name    string
		fields  fields
		req     io.Reader // リクエスト
		want    interface{}
		wantErr bool
	}{
		{
			"ok",
			fields{DB: moqDb, Repo: moqRepo},
			fixture.LoginFormBody(&model.Login{
				Username: UserName,
				Password: RawPassword,
			}),
			validUser,
			false,
		},
		// フォームのUserNameが未入力
		{
			"userNameNotFilled",
			fields{DB: moqDb, Repo: moqRepo},
			fixture.LoginFormBody(&model.Login{
				// Usernameは未定義
				Password: RawPassword,
			}),
			nil,
			true,
		},
		// フォームのPasswordが未入力
		{
			"passwordNotFilled",
			fields{DB: moqDb, Repo: moqRepo},
			fixture.LoginFormBody(&model.Login{
				Username: UserName,
				// Passwordは未定義
			}),
			nil,
			true,
		},
		// ユーザーネームが不正
		{
			"invalidUserName",
			fields{DB: moqDb, Repo: moqRepo},
			fixture.LoginFormBody(&model.Login{
				Username: "invalidUserName",
				Password: RawPassword,
			}),
			nil,
			true,
		},
		// パスワードが不正
		{
			"invalidPassword",
			fields{DB: moqDb, Repo: moqRepo},
			fixture.LoginFormBody(&model.Login{
				Username: UserName,
				Password: "invalidPassword",
			}),
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			as := &AuthService{
				DB:   tt.fields.DB,
				Repo: tt.fields.Repo,
			}
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			// リクエストを作成
			req := httptest.NewRequest("POST", "/login", tt.req)
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			// リクエスト情報をコンテキストに入れる
			c.Request = req

			got, err := as.Authenticator(c)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthService.Authenticator() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AuthService.Authenticator() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthService_PayloadFunc(t *testing.T) {
	type args struct {
		data interface{}
	}
	validUser := fixture.User(&model.User{
		UserName: UserName,
		Role:     Role,
	})

	tests := []struct {
		name         string
		args         args
		wantUserName string
		wantRole     string
	}{
		{
			"ok",
			args{validUser},
			UserName,
			Role,
		},
		// 引数の型が*model.Userでない場合
		{
			"notUser",
			args{nil},
			"",
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			moqStore := &StoreMock{
				SaveFunc: func(ctx context.Context, key string, uid int) error {
					return nil
				},
			}
			as := &AuthService{
				DB:    &DBConnectionMock{},
				Repo:  &AuthRepositoryMock{},
				Store: moqStore,
			}

			got := as.PayloadFunc(tt.args.data)
			if tt.name == "ok" {
				assert.Equal(t, "github.com/kwtryo/go-sample", got[IssuerKey])
				assert.Equal(t, "access_token", got[SubjectKey])
				assert.Equal(t, "github.com/kwtryo/go-sample", got[AudienceKey])
				assert.Equal(t, tt.wantUserName, got[UserNameKey])
				assert.Equal(t, tt.wantRole, got[RoleKey])
			} else {
				// 未定義であることを確認
				assert.Equal(t, nil, got[IssuerKey])
				assert.Equal(t, nil, got[SubjectKey])
				assert.Equal(t, nil, got[AudienceKey])
				assert.Equal(t, nil, got[UserNameKey])
				assert.Equal(t, nil, got[RoleKey])
			}
		})
	}
}

func TestAuthService_IdentityHandler(t *testing.T) {
	type payloadFuncArgs struct {
		data interface{}
	}
	validUser := fixture.User(&model.User{
		UserName: UserName,
		Role:     Role,
	})

	tests := []struct {
		name            string
		payloadFuncArgs payloadFuncArgs
		want            interface{}
	}{
		{
			"ok",
			payloadFuncArgs{validUser},
			&model.User{
				UserName: UserName,
				Role:     Role,
			},
		},
		// ユーザーが見つからない場合
		{
			"notFound",
			payloadFuncArgs{validUser},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			moqStore := &StoreMock{
				SaveFunc: func(ctx context.Context, key string, uid int) error {
					return nil
				},
				LoadFunc: func(ctx context.Context, key string) (int, error) {
					if tt.name == "notFound" {
						return 0, errors.New("error from mock")
					}
					return 0, nil
				},
			}
			as := &AuthService{
				DB:    &DBConnectionMock{},
				Repo:  &AuthRepositoryMock{},
				Store: moqStore,
			}

			// リクエストの値がnullになっている
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			c.Set("JWT_PAYLOAD", as.PayloadFunc(tt.payloadFuncArgs.data))
			// リクエストを作成
			req := httptest.NewRequest("GET", "/test", nil)
			// リクエスト情報をコンテキストに入れる
			c.Request = req

			if got := as.IdentityHandler(c); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AuthService.IdentityHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthService_Authorizator(t *testing.T) {
	type fields struct {
		DB   store.DBConnection
		Repo AuthRepository
	}
	type args struct {
		data interface{}
	}
	validUser := fixture.User(&model.User{
		UserName: UserName,
	})
	moqDb := &DBConnectionMock{}
	moqRepo := &AuthRepositoryMock{}
	moqRepo.GetUserByUserNameFunc = func(ctx context.Context, db store.DBConnection, userName string) (*model.User, error) {
		if userName == UserName {
			return validUser, nil
		}
		return nil, fmt.Errorf("error from mock")
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			"ok",
			fields{DB: moqDb, Repo: moqRepo},
			args{validUser},
			true,
		},
		// 引数がユーザーでない場合
		{
			"notUser",
			fields{DB: moqDb, Repo: moqRepo},
			args{nil},
			false,
		},
		// ユーザーが存在しない場合
		{
			"invalidUser",
			fields{DB: moqDb, Repo: moqRepo},
			args{fixture.User(&model.User{
				UserName: "invalidUser",
			})},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			as := &AuthService{
				DB:   tt.fields.DB,
				Repo: tt.fields.Repo,
			}
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			// リクエストを作成
			req := httptest.NewRequest("GET", "/auth", nil)
			// リクエスト情報をコンテキストに入れる
			c.Request = req

			if got := as.Authorizator(tt.args.data, c); got != tt.want {
				t.Errorf("AuthService.Authorizator() = %v, want %v", got, tt.want)
			}
		})
	}
}
