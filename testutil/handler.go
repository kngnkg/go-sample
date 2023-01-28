package testutil

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func LoadFile(t *testing.T, path string) []byte {
	t.Helper()

	bt, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("cannot read from %q: %v", path, err)
	}
	return bt
}

// レスポンスを検証する
func AssertResponse(t *testing.T, got *http.Response, status int, body []byte) {
	t.Helper()
	t.Cleanup(func() { _ = got.Body.Close() })
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
