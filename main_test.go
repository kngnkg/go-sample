package main

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRun(t *testing.T) {
	r := gin.Default()
	// テスト用のエンドポイントを作成
	//
	srv := &http.Server{
		Addr:    ":0",
		Handler: r,
	}
	go run(srv)
	// ポート番号を確認
	//

	// レスポンスを確認 test.http
	//

	// サーバーの終了動作を確認
	//

	// 戻り値を確認
	//
}
