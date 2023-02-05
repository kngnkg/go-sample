package handler

import (
	"bytes"
	"context"
	_ "embed"
	"io"
	"net/http"
	"testing"

	"github.com/kwtryo/go-sample/clock"
	"github.com/kwtryo/go-sample/model"
	"github.com/kwtryo/go-sample/service"
	"github.com/kwtryo/go-sample/store"
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
	)
	type fields struct {
		Service AuthService
		Clocker clock.Clocker
	}
	cl := clock.FixedClocker{}
	db := testutil.OpenDbForTest(t)
	repo := store.Repository{Clocker: cl}
	as := &service.AuthService{DB: db, Repo: &repo}

	// テスト用ユーザーをDBに登録する
	// 不具合が起こりそうなので考える
	err := repo.DeleteUserAll(context.Background(), db)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	user := fixture.User(&model.User{
		UserName: UserName,
		Password: RawPassword,
	})
	_, err = repo.RegisterUser(context.Background(), db, user)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	tests := []struct {
		name         string
		fields       fields
		req          io.Reader // リクエスト
		wantStatus   int       // ステータスコード
		wantRespFile string    // レスポンス
	}{
		{
			"ok",
			fields{Service: as, Clocker: cl},
			fixture.LoginFormBody(&model.Login{
				Username: UserName,
				Password: RawPassword,
			}),
			http.StatusOK,
			"testdata/login/ok_response.json.golden",
		},
		// ユーザーネームが不正
		{
			"invalidUserName",
			fields{Service: as, Clocker: cl},
			fixture.LoginFormBody(&model.Login{
				Username: "invalidUserName",
				Password: RawPassword,
			}),
			http.StatusUnauthorized,
			"testdata/login/invalid_username_response.json.golden",
		},
		// パスワードが不正
		{
			"invalidPassword",
			fields{Service: as, Clocker: cl},
			fixture.LoginFormBody(&model.Login{
				Username: UserName,
				Password: "invalidPassword",
			}),
			http.StatusUnauthorized,
			"testdata/login/invalid_pass_response.json.golden",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &JWTer{
				Service: tt.fields.Service,
				Clocker: tt.fields.Clocker,
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
