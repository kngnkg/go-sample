package store

import (
	"context"
	"testing"
	"time"

	"github.com/kngnkg/go-sample/testutil"
	"github.com/redis/go-redis/v9"
)

func TestKVS_Save(t *testing.T) {
	type fields struct {
		Cli *redis.Client
	}
	type args struct {
		ctx context.Context
		key string
		uid int
	}

	ctx := context.Background()
	cli := testutil.OpenRedisForTest(t)
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"ok",
			fields{Cli: cli},
			args{
				ctx: ctx,
				key: "TestKVS_Save",
				uid: 1234,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := &KVS{
				Cli: tt.fields.Cli,
			}
			if err := k.Save(tt.args.ctx, tt.args.key, tt.args.uid); (err != nil) != tt.wantErr {
				t.Errorf("KVS.Save() error = %v, wantErr %v", err, tt.wantErr)
			}
			t.Cleanup(func() {
				k.Cli.Del(tt.args.ctx, tt.args.key)
			})
		})
	}
}

func TestKVS_Load(t *testing.T) {
	type fields struct {
		Cli *redis.Client
	}
	type args struct {
		ctx context.Context
		key string
	}

	ctx := context.Background()
	cli := testutil.OpenRedisForTest(t)
	key := "TestKVS_Load"
	uid := 1234
	cli.Set(ctx, key, uid, 30*time.Minute)

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		{
			"ok",
			fields{Cli: cli},
			args{
				ctx: ctx,
				key: key,
			},
			uid,
			false,
		},
		// 見つからない場合
		{
			"notFound",
			fields{Cli: cli},
			args{
				ctx: ctx,
				key: "TestKVS_Load_not_found", // keyが不正
			},
			0,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := &KVS{
				Cli: tt.fields.Cli,
			}
			got, err := k.Load(tt.args.ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("KVS.Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("KVS.Load() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKVS_Delete(t *testing.T) {
	type fields struct {
		Cli *redis.Client
	}
	type args struct {
		ctx context.Context
		key string
	}

	ctx := context.Background()
	cli := testutil.OpenRedisForTest(t)
	key := "TestKVS_Delete"
	uid := 1234

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"ok",
			fields{Cli: cli},
			args{
				ctx: ctx,
				key: key,
			},
			true,
		},
		// 見つからない場合
		{
			"notFound",
			fields{Cli: cli},
			args{
				ctx: ctx,
				key: "TestKVS_Delete_not_found", // keyが不正
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := &KVS{
				Cli: tt.fields.Cli,
			}
			// 正常値を登録する
			if err := k.Save(tt.args.ctx, key, uid); err != nil {
				t.Errorf("KVS.Save() error = %v", err)
			}
			t.Cleanup(func() {
				k.Cli.Del(tt.args.ctx, key)
			})

			err := k.Delete(tt.args.ctx, tt.args.key)
			if err != nil {
				t.Errorf("KVS.Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// 削除されているか確認する
			_, err = k.Load(tt.args.ctx, key)
			if (err != nil) != tt.wantErr {
				t.Errorf("KVS.Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
