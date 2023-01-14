package handler

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kwtryo/go-sample/clock"
	"github.com/kwtryo/go-sample/model"
	"github.com/kwtryo/go-sample/store"
	"github.com/kwtryo/go-sample/testutil"
	"github.com/stretchr/testify/assert"
)

type userHandlerTest struct {
	c       *gin.Context
	router  *gin.Engine
	rec     *httptest.ResponseRecorder
	handler *UserHandler
}

func prepareTest(t *testing.T) *userHandlerTest {
	t.Helper()

	router := gin.Default()
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	db := testutil.OpenDbForTest(t)
	repo := &store.Repository{Clocker: clock.FixedClocker{}}

	uht := &userHandlerTest{
		c:      c,
		router: router,
		rec:    rec,
		handler: &UserHandler{
			DB:   db,
			Repo: repo,
		},
	}
	uht.handler.Repo.DeleteUserAll(c, db)
	return uht
}

func TestRegisterUserRoute(t *testing.T) {
	uht := prepareTest(t)

	uht.router.POST("/register", uht.handler.RegisterUser)

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

	uht.c.Request = httptest.NewRequest("POST", "/register", body)
	uht.c.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	uht.handler.RegisterUser(uht.c)

	assert.Equal(t, http.StatusOK, uht.rec.Code)
	// assert.Contains(t, w.Body.String(), "登録するユーザのID")
}
