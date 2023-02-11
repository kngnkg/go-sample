package testutil

import (
	"context"
	"fmt"
	"testing"

	"github.com/redis/go-redis/v9"
)

// テスト用のRedisクライアントオブジェクトを返す
func OpenRedisForTest(t *testing.T) *redis.Client {
	t.Helper()

	cfg := CreateConfigForTest(t)
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort),
		Password: "",
		DB:       0,
	})
	if err := client.Ping(context.Background()).Err(); err != nil {
		t.Fatalf("failed to connect redis: %v", err)
	}
	return client
}
