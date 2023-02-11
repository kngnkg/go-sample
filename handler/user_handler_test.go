package handler

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/kwtryo/go-sample/clock"
	"github.com/kwtryo/go-sample/model"
	"github.com/kwtryo/go-sample/store"
	"github.com/kwtryo/go-sample/testutil"
	"github.com/kwtryo/go-sample/testutil/fixture"
)

func TestUserHandler_RegisterUser(t *testing.T) {
	const (
		UserName = "testUserName"
	)
	cl := clock.FixedClocker{}
	validUser := fixture.User(&model.User{
		Name:     "testUserFullName",
		UserName: UserName,
		Password: "testPassword",
		Created:  cl.Now(),
		Modified: cl.Now(),
	})
	form := fixture.RegisterUserBody(validUser)
	invalidForm := fixture.RegisterUserBody(validUser)
	// 不正な値にするためにnameを削除する
	invalidForm.Del("name")

	tests := []struct {
		name         string
		req          io.Reader // リクエスト
		wantStatus   int       // ステータスコード
		wantRespFile string    // レスポンス
	}{
		{
			"ok",
			strings.NewReader(form.Encode()),
			http.StatusOK,
			"testdata/register_user/ok_response.json.golden",
		},
		// リクエストが不正な場合
		{
			"badRequest",
			strings.NewReader(invalidForm.Encode()),
			http.StatusBadRequest,
			"testdata/register_user/bad_req_response.json.golden",
		},
		// 内部エラー
		{
			"internalServerError",
			strings.NewReader(form.Encode()),
			http.StatusInternalServerError,
			"testdata/register_user/server_err_response.json.golden",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			moqService := &UserServiceMock{}
			moqService.RegisterUserFunc = func(ctx context.Context, form *model.FormRequest) (*model.User, error) {
				if tt.name == "internalServerError" {
					return nil, errors.New("error from mock")
				}
				validUser.Id = 1
				return validUser, nil
			}
			uh := &UserHandler{
				Service: moqService,
			}

			testutil.CheckHandlerFunc(
				t,
				uh.RegisterUser,
				"POST",
				"",
				tt.req,
				tt.wantStatus,
				tt.wantRespFile,
			)
		})
	}
}

func TestUserHandler_GetUser(t *testing.T) {
	const (
		UserName = "testUserName"
	)
	cl := clock.FixedClocker{}
	validUser := fixture.User(&model.User{
		Id:       1,
		Name:     "testUserFullName",
		UserName: UserName,
		Password: "testPassword",
		Created:  cl.Now(),
		Modified: cl.Now(),
	})

	tests := []struct {
		name         string
		queryParam   string // クエリパラメータ
		wantStatus   int    // ステータスコード
		wantRespFile string // レスポンス
	}{
		{
			"ok",
			validUser.UserName,
			http.StatusOK,
			"testdata/get_user/ok_response.json.golden",
		},
		{
			"notFound",
			"invalidUserName",
			http.StatusNotFound,
			"testdata/get_user/not_found_response.json.golden",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			moqService := &UserServiceMock{}
			moqService.GetUserFunc = func(ctx context.Context, userName string) (*model.User, error) {
				if userName == UserName {
					validUser.Password = "hashedPassword"
					return validUser, nil
				}
				return nil, store.ErrNotFound
			}
			uh := &UserHandler{
				Service: moqService,
			}

			testutil.CheckHandlerFunc(
				t,
				uh.GetUser,
				"GET",
				"?user_name="+tt.queryParam,
				nil,
				tt.wantStatus,
				tt.wantRespFile,
			)
		})
	}
}
