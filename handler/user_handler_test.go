package handler

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/kngnkg/go-sample/clock"
	"github.com/kngnkg/go-sample/model"
	"github.com/kngnkg/go-sample/store"
	"github.com/kngnkg/go-sample/testutil"
	"github.com/kngnkg/go-sample/testutil/fixture"
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

func TestUserHandler_GetAllUsers(t *testing.T) {
	// ランダムなユーザーを5人生成する
	cl := clock.FixedClocker{}
	users := model.Users{}
	for i := 1; i < 6; i++ {
		str := strconv.Itoa(i)
		user := fixture.User(&model.User{
			Id:       i,
			Name:     "fullName" + str,
			UserName: "userName" + str,
			Password: "testPassword",
			Created:  cl.Now(),
			Modified: cl.Now(),
		})
		user.Password = "testPassword"
		users = append(users, user)
	}
	moqErr := errors.New("err from mock")

	tests := []struct {
		name         string
		wantStatus   int    // ステータスコード
		wantRespFile string // レスポンス
	}{
		{
			"ok",
			http.StatusOK,
			"testdata/get_all_users/ok_response.json.golden",
		},
		// ユーザーが見つからなかった場合
		{
			"notFound",
			http.StatusInternalServerError,
			"testdata/get_all_users/not_found_response.json.golden",
		},
		// 内部エラー
		{
			"serverError",
			http.StatusInternalServerError,
			"testdata/get_all_users/server_err_response.json.golden",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			moqService := &UserServiceMock{}
			moqService.GetAllUsersFunc =
				func(ctx context.Context) (model.Users, error) {
					if tt.name == "serverError" {
						return nil, moqErr
					}
					if tt.name == "notFound" {
						return nil, store.ErrNotFound
					}
					return users, nil
				}
			uh := &UserHandler{
				Service: moqService,
			}

			testutil.CheckHandlerFunc(
				t,
				uh.GetAllUsers,
				"GET",
				"",
				nil,
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
		{
			"internalServerError",
			validUser.UserName,
			http.StatusInternalServerError,
			"testdata/get_user/server_err_response.json.golden",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			moqService := &UserServiceMock{}
			moqService.GetUserFunc = func(ctx context.Context, userName string) (*model.User, error) {
				if userName != UserName {
					return nil, store.ErrNotFound
				}
				if tt.name == "internalServerError" {
					return nil, errors.New("err from mock")
				}
				validUser.Password = "hashedPassword"
				return validUser, nil
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
