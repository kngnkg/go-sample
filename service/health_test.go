package service

import (
	"context"
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
			return nil
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hs := &HealthService{
				DB:   tt.fields.DB,
				Repo: tt.fields.Repo,
			}
			if err := hs.HealthCheck(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("HealthService.HealthCheck() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
