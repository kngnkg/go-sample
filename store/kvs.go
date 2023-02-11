package store

import (
	"context"
	"fmt"
	"time"

	"github.com/kwtryo/go-sample/config"
	"github.com/redis/go-redis/v9"
)

type KVS struct {
	Cli *redis.Client
}

func NewKVS(ctx context.Context, cfg *config.Config) (*KVS, error) {
	cli := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort),
	})
	if err := cli.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	return &KVS{Cli: cli}, nil
}

// Keyをセットする
func (k *KVS) Save(ctx context.Context, key string, uid int) error {
	// TTLは30分に設定
	return k.Cli.Set(ctx, key, uid, 30*time.Minute).Err()
}

// Keyから値を取得する
func (k *KVS) Load(ctx context.Context, key string) (int, error) {
	id, err := k.Cli.Get(ctx, key).Int()
	if err != nil {
		return 0, fmt.Errorf("failed to get by %q: %w", key, ErrNotFound)
	}
	return id, nil
}

func (k *KVS) Delete(ctx context.Context, key string) error {
	return k.Cli.Del(ctx, key).Err()
}
