package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kwtryo/go-sample/model"
	"github.com/stretchr/testify/assert"
)

type userHandlerTest struct {
	c      *gin.Context
	router *gin.Engine
	rec    *httptest.ResponseRecorder
	// handler *UserHandler
}

// userHandlerTest構造体を初期化する
func prepareTest(t *testing.T) *userHandlerTest {
	t.Helper()

	router := gin.Default()
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	// db := testutil.OpenDbForTest(t)
	// repo := &store.Repository{Clocker: clock.FixedClocker{}}

	uht := &userHandlerTest{
		c:      c,
		router: router,
		rec:    rec,
		// handler: &UserHandler{
		// 	DB:   db,
		// 	Repo: repo,
		// },
	}
	// uht.handler.Repo.DeleteUserAll(c, db)
	return uht
}

func TestRegisterUserRoute(t *testing.T) {
	type want struct {
		status int
		body   *strings.Reader
	}

	type test struct {
		body *strings.Reader
		want want
	}

	tests := map[string]test{
		// 正常系
		"ok": {
			body: validBody(),
			want: want{
				status: http.StatusOK,
				body:   validBody(),
			},
		},
		// 異常系
		"badRequest": {
			body: invalidBody(),
			want: want{
				status: http.StatusBadRequest,
				body:   invalidBody(),
			},
		},
	}

	for n, tst := range tests {
		t.Run(n, func(t *testing.T) {
			tst := tst
			t.Parallel()

			uht := prepareTest(t)
			mockedUserService := &UserServiceMock{}
			mockedUserService.RegisterUserFunc = func(ctx context.Context) (int, error) {
				if tst.want.status == http.StatusOK {
					return 1, nil
				}
				return 0, errors.New("error from mock")
			}
			handler := &UserHandler{
				Service: mockedUserService,
			}
			uht.router.POST("/register", handler.RegisterUser)

			uht.c.Request = httptest.NewRequest("POST", "/register", tst.body)
			uht.c.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			handler.RegisterUser(uht.c)

			// 期待するステータスコードと一致するか確認する
			assert.Equal(t, tst.want.status, uht.rec.Code)
		})
	}
}

func validBody() *strings.Reader {
	// テスト用ユーザー
	u := &model.User{
		Name:     "testUserFullName",
		UserName: "testUser",
		Password: "testPassword",
		Role:     "admin",
		Email:    "test@example.com",
		Address:  "testAddress",
		Phone:    "000-0000-0000",
		Website:  "ttp://test.com",
		Company:  "testCompany",
	}
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
func invalidBody() *strings.Reader {
	// テスト用ユーザー
	u := &model.User{
		// Name:     "testUserFullName",
		UserName: "testUser",
		Password: "testPassword",
		Role:     "admin",
		Email:    "test@example.com",
		Address:  "testAddress",
		Phone:    "000-0000-0000",
		Website:  "ttp://test.com",
		Company:  "testCompany",
	}
	// リクエストを作成
	form := url.Values{}
	// form.Add("name", u.Name)
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
