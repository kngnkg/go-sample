package testutil

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// ハンドラ関数を検証する。
//
// handler: 検証したいハンドラ
// reqMethod: リクエストメソッド
// queryParam: クエリパラメータ
// reqBody: リクエストボディ
// wantStatus: 期待するステータスコード
// wantRespFile: 期待するレスポンスボディのファイルパス
func CheckHandlerFunc(
	t *testing.T,
	handler gin.HandlerFunc,
	reqMethod string,
	queryParam string,
	reqBody io.Reader,
	wantStatus int,
	wantRespFile string) {
	t.Helper()

	url := RunTestServer(t, handler, reqMethod)

	if queryParam != "" {
		url = url + queryParam
	}
	resp := SendRequest(
		t,
		url,
		reqMethod,
		queryParam,
		reqBody,
	)

	AssertResponse(
		t,
		resp,
		wantStatus,
		LoadFile(t, wantRespFile),
	)
}

// ミドルウェアを検証する。
//
// middleware: 検証したいミドルウェア
// wantStatus: 期待するステータスコード
// wantRespFile: 期待するレスポンスボディのファイルパス
func CheckMiddleware(
	t *testing.T,
	middleware gin.HandlerFunc,
	wantStatus int,
	wantRespFile string) {
	t.Helper()

	url := RunTestServerWithMiddleware(t, middleware)

	resp := SendRequest(
		t,
		url,
		"GET",
		"",
		nil,
	)

	AssertResponse(
		t,
		resp,
		wantStatus,
		LoadFile(t, wantRespFile),
	)
}

// テスト用のサーバーを起動し、URLを返す。
// handler: 検証したいハンドラ
// reqMethod: リクエストメソッド
func RunTestServer(t *testing.T, handler gin.HandlerFunc, reqMethod string) string {
	t.Helper()
	router := gin.Default()
	end := "/test"

	// ハンドラの登録
	switch reqMethod {
	case "GET":
		router.GET(end, handler)
	case "POST":
		router.POST(end, handler)
	}

	testServer := httptest.NewServer(router) // サーバを立てる
	t.Cleanup(func() { testServer.Close() })

	return fmt.Sprintf(testServer.URL + end)
}

// テスト用のサーバーを起動し、URLを返す。
// middleware: 検証したいミドルウェア
func RunTestServerWithMiddleware(t *testing.T, middleware gin.HandlerFunc) string {
	t.Helper()

	handler := func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "hello",
		})
	}

	router := gin.Default()
	group := router.Group("/group")
	group.Use(middleware)
	{
		group.GET("/test", handler)
	}
	testServer := httptest.NewServer(router) // サーバを立てる
	t.Cleanup(func() { testServer.Close() })

	return fmt.Sprintf(testServer.URL + "/group/test")
}

// リクエストを送信し、レスポンスを返す。
// url: リクエストを送信する対象のURL
// reqMethod: リクエストメソッド
// reqBody: リクエストボディ
func SendRequest(t *testing.T, url string, reqMethod string, queryParam string, reqBody io.Reader) *http.Response {
	t.Helper()
	t.Logf("try request to %q", url)

	var (
		resp *http.Response
		err  error
	)
	// テストサーバーにリクエストを送信
	switch reqMethod {
	case "GET":
		resp, err = http.Get(url)
	case "POST":
		resp, err = http.Post(url, "application/x-www-form-urlencoded", reqBody)
	}
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	t.Cleanup(func() { _ = resp.Body.Close() })
	return resp
}

// レスポンスを検証する
func AssertResponse(t *testing.T, got *http.Response, status int, body []byte) {
	t.Helper()

	gb, err := io.ReadAll(got.Body)
	if err != nil {
		t.Fatal(err)
	}
	// ステータスコードの確認
	assert.Equal(t, status, got.StatusCode)

	if len(gb) == 0 && len(body) == 0 {
		// レスポンスボディが無い場合は確認不要
		return
	}

	// レスポンスボディの確認
	var jw, jg any
	if err := json.Unmarshal(body, &jw); err != nil {
		t.Fatalf("cannot unmarshal want %q: %v", body, err)
	}
	if err := json.Unmarshal(gb, &jg); err != nil {
		t.Fatalf("cannot unmarshal got %v: %v", got, err)
	}
	assert.Equal(t, jw, jg)
}

func LoadFile(t *testing.T, path string) []byte {
	t.Helper()

	bt, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("cannot read from %q: %v", path, err)
	}
	return bt
}
