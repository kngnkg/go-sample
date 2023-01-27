package store

import (
	"context"
	"testing"

	"github.com/kwtryo/go-sample/clock"
	"github.com/kwtryo/go-sample/testutil"
)

func TestRepository_Ping(t *testing.T) {
	type fields struct {
		Clocker clock.Clocker
	}
	type args struct {
		ctx context.Context
		db  DBConnection
	}

	ctx := context.Background()
	db := testutil.OpenDbForTest(t)
	tx, _ := db.BeginTxx(ctx, nil)
	defer tx.Rollback()
	cl := clock.FixedClocker{}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"ok",
			fields{Clocker: cl},
			args{ctx: ctx, db: db},
			false,
		},
		// 引数が*sqlx.dbでない場合
		{
			"invalidArgs",
			fields{Clocker: cl},
			args{ctx: ctx, db: tx},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Repository{
				Clocker: tt.fields.Clocker,
			}
			if err := r.Ping(tt.args.ctx, tt.args.db); (err != nil) != tt.wantErr {
				t.Errorf("Repository.Ping() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
