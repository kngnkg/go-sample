package service

import (
	"context"
	"errors"
	"testing"

	"github.com/kwtryo/go-sample/store"
)

func TestHealthService_HealthCheck(t *testing.T) {
	type fields struct {
		DB   store.DBConnection
		Repo HealthRepository
	}
	type args struct {
		ctx context.Context
	}

	ctx := context.Background()
	moqDb := &DBConnectionMock{}
	moqRepo := &HealthRepositoryMock{
		PingFunc: func(ctx context.Context, db store.DBConnection) error {
			// コンテキストからテストの名前を取得する
			testName, ok := ctx.Value(TEST_NAME_KEY).(string)
			if !ok {
				t.Fatal("unexpected error")
			}
			if testName == "ok" {
				return nil
			}
			return errors.New("error")
		},
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"ok",
			fields{DB: moqDb, Repo: moqRepo},
			args{ctx: ctx},
			false,
		},
		{
			"error",
			fields{DB: moqDb, Repo: moqRepo},
			args{ctx: ctx},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hs := &HealthService{
				DB:   tt.fields.DB,
				Repo: tt.fields.Repo,
			}

			// コンテキストに現在のテストの名前を入れる
			ctx := context.WithValue(tt.args.ctx, TEST_NAME_KEY, tt.name)
			if err := hs.HealthCheck(ctx); (err != nil) != tt.wantErr {
				t.Errorf("HealthService.HealthCheck() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
