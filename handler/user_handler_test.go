package handler

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kwtryo/go-sample/model"
	"github.com/kwtryo/go-sample/store"
	"github.com/kwtryo/go-sample/testutil"
)

func TestUserHandler_RegisterUser(t *testing.T) {
	tests := []struct {
		name         string
		req          io.Reader // リクエスト
		wantStatus   int       // ステータスコード
		wantRespFile string    // レスポンス
	}{
		{
			"ok",
			validBody(t),
			http.StatusOK,
			"testdata/register_user/ok_response.json.golden",
		},
		// リクエストが不正な場合
		{
			"badRequest",
			invalidBody(t),
			http.StatusBadRequest,
			"testdata/register_user/bad_req_response.json.golden",
		},
		// 内部エラー
		{
			"internalServerError",
			validBody(t),
			http.StatusInternalServerError,
			"testdata/register_user/server_err_response.json.golden",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			moqService := &UserServiceMock{}
			moqService.RegisterUserFunc = func(ctx context.Context, form *model.FormRequest) (*model.User, error) {
				if tt.name == "ok" {
					u := testutil.GetTestUser(t)
					u.Id = 1
					return u, nil
				}
				return nil, errors.New("error from mock")
			}
			uh := &UserHandler{
				Service: moqService,
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			// リクエストを作成
			req := httptest.NewRequest("POST", "/register", tt.req)
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			// リクエスト情報をコンテキストに入れる
			c.Request = req
			uh.RegisterUser(c)
			resp := w.Result()

			testutil.AssertResponse(
				t,
				resp,
				tt.wantStatus,
				testutil.LoadFile(t, tt.wantRespFile),
			)
		})
	}
}

func TestUserHandler_GetUser(t *testing.T) {
	tests := []struct {
		name         string
		queryParam   string // クエリパラメータ
		wantStatus   int    // ステータスコード
		wantRespFile string // レスポンス

	}{
		{
			"ok",
			testutil.VALID_USER_NAME,
			http.StatusOK,
			"testdata/get_user/ok_response.json.golden",
		},
		{
			"notFound",
			testutil.INVALID_USER_NAME,
			http.StatusNotFound,
			"testdata/get_user/not_found_response.json.golden",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			moqService := &UserServiceMock{}
			moqService.GetUserFunc = func(ctx context.Context, userName string) (*model.User, error) {
				if userName == testutil.VALID_USER_NAME {
					u := testutil.GetTestUser(t)
					u.Id = 1
					return u, nil
				}
				return nil, store.ErrNotFound
			}
			uh := &UserHandler{
				Service: moqService,
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			// リクエストを作成
			req := httptest.NewRequest("GET", "/user?user_name="+tt.queryParam, nil)
			// リクエスト情報をコンテキストに入れる
			c.Request = req
			uh.GetUser(c)
			resp := w.Result()

			testutil.AssertResponse(
				t,
				resp,
				tt.wantStatus,
				testutil.LoadFile(t, tt.wantRespFile),
			)
		})
	}
}

func validBody(t *testing.T) *strings.Reader {
	u := testutil.GetTestUser(t)

	// リクエストを作成
	form := url.Values{}
	form.Add("name", u.Name)
	form.Add("username", u.UserName)
	form.Add("password", u.Password)
	form.Add("role", u.Role)
	form.Add("email", u.Email)
	form.Add("address", u.Address)
	form.Add("phone", u.Phone)
	form.Add("website", u.Website)
	form.Add("company", u.Company)
	body := strings.NewReader(form.Encode())
	return body
}

func invalidBody(t *testing.T) *strings.Reader {
	u := testutil.GetTestUser(t)

	// リクエストを作成
	form := url.Values{}
	// nameを設定しない
	form.Add("username", u.UserName)
	form.Add("password", u.Password)
	form.Add("role", u.Role)
	form.Add("email", u.Email)
	form.Add("address", u.Address)
	form.Add("phone", u.Phone)
	form.Add("website", u.Website)
	form.Add("company", u.Company)
	body := strings.NewReader(form.Encode())
	return body
}
