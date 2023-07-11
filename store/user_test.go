package store

import (
	"context"
	"errors"
	"testing"

	"github.com/kngnkg/go-sample/clock"
	"github.com/kngnkg/go-sample/model"
	"github.com/kngnkg/go-sample/testutil"
	"github.com/kngnkg/go-sample/testutil/fixture"
	"github.com/stretchr/testify/assert"
)

func TestRepository_RegisterUser(t *testing.T) {
	type fields struct {
		Clocker clock.Clocker
	}
	user := fixture.User(&model.User{
		Id: 0, // 未設定
	})
	tests := []struct {
		name    string
		fields  fields
		want    *model.User
		wantErr error
	}{
		// 正常系
		{
			"ok",
			fields{&clock.FixedClocker{}},
			user,
			nil,
		},
		// 既に登録されていた場合
		{
			"errAlreadyEntry",
			fields{&clock.FixedClocker{}},
			nil,
			ErrAlreadyEntry,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Repository{
				Clocker: tt.fields.Clocker,
			}

			ctx := context.TODO()
			tx, err := testutil.OpenDbForTest(t).BeginTxx(ctx, nil)
			t.Cleanup(func() { _ = tx.Rollback() })
			if err != nil {
				t.Fatal(err)
			}

			testutil.DeleteUserAll(ctx, t, tx)
			if tt.name == "errAlreadyEntry" {
				// 先にユーザーを登録する
				_ = testutil.PrepareUser(ctx, t, r.Clocker, tx, user)
			}
			got, err := r.RegisterUser(ctx, tx, user)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("want err: %v but got: %v", tt.wantErr, err)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestRepository_GetByUserName(t *testing.T) {
	const userName = "testUserName"

	type fields struct {
		Clocker clock.Clocker
	}
	user := fixture.User(&model.User{
		Id:       0, // 未設定
		UserName: userName,
	})
	tests := []struct {
		name     string
		fields   fields
		userName string
		want     *model.User
		wantErr  error
	}{
		// 正常系
		{
			"ok",
			fields{&clock.FixedClocker{}},
			userName,
			user,
			nil,
		},
		// ユーザーが存在しない場合
		{
			"notFound",
			fields{&clock.FixedClocker{}},
			"invalidUserName",
			nil,
			ErrNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Repository{
				Clocker: tt.fields.Clocker,
			}

			ctx := context.TODO()
			tx, err := testutil.OpenDbForTest(t).BeginTxx(ctx, nil)
			t.Cleanup(func() { _ = tx.Rollback() })
			if err != nil {
				t.Fatal(err)
			}

			testutil.DeleteUserAll(ctx, t, tx)
			// ユーザーを登録する
			_ = testutil.PrepareUser(ctx, t, r.Clocker, tx, user)

			got, err := r.GetUserByUserName(ctx, tx, tt.userName)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("want err: %v but got: %v", tt.wantErr, err)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestRepository_GetAllUsers(t *testing.T) {
	type fields struct {
		Clocker clock.Clocker
	}
	// ランダムなユーザーを5人生成する
	users := model.Users{}
	for i := 0; i < 5; i++ {
		user := fixture.User(&model.User{})
		users = append(users, user)
	}
	tests := []struct {
		name    string
		fields  fields
		users   model.Users // 登録するユーザー
		want    []*model.User
		wantErr error
	}{
		// 正常系
		{
			"ok",
			fields{&clock.FixedClocker{}},
			users,
			users,
			nil,
		},
		// DBにユーザーが1人もいない場合
		{
			"notExists",
			fields{&clock.FixedClocker{}},
			model.Users{},
			nil,
			ErrNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Repository{
				Clocker: tt.fields.Clocker,
			}

			ctx := context.TODO()
			tx, err := testutil.OpenDbForTest(t).BeginTxx(ctx, nil)
			t.Cleanup(func() { _ = tx.Rollback() })
			if err != nil {
				t.Fatal(err)
			}

			testutil.DeleteUserAll(ctx, t, tx)
			// ユーザーを登録する
			for _, v := range tt.users {
				_ = testutil.PrepareUser(ctx, t, r.Clocker, tx, v)
			}

			got, err := r.GetAllUsers(ctx, tx)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("want err: %v but got: %v", tt.wantErr, err)
			}

			if tt.name != "ok" {
				// 終了
				return
			}
			for i, v := range got {
				assert.Equal(t, tt.want[i].Name, v.Name)
				assert.Equal(t, tt.want[i].UserName, v.UserName)
				assert.Empty(t, v.Password) // Passwordは取得しない
				assert.Equal(t, tt.want[i].Role, v.Role)
				assert.Equal(t, tt.want[i].Email, v.Email)
				assert.Equal(t, tt.want[i].Address, v.Address)
				assert.Equal(t, tt.want[i].Phone, v.Phone)
				assert.Equal(t, tt.want[i].Website, v.Website)
				assert.Equal(t, tt.want[i].Company, v.Company)
				assert.Equal(t, tt.want[i].Created, v.Created)
				assert.Equal(t, tt.want[i].Modified, v.Modified)
			}
		})
	}
}
