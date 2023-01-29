package service

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/kwtryo/go-sample/model"
	"github.com/kwtryo/go-sample/store"
	"github.com/kwtryo/go-sample/testutil"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

type key int

const (
	// コンテキストに入れるテストの名前のKey
	TEST_NAME_KEY key = iota
)

func TestUserService_RegisterUser(t *testing.T) {
	type fields struct {
		DB   store.DBConnection
		Repo UserRepository
	}
	type args struct {
		ctx  context.Context
		form *model.FormRequest
	}

	moqDb := &DBConnectionMock{}
	moqRepo := &UserRepositoryMock{}
	moqRepo.RegisterUserFunc =
		func(ctx context.Context, db store.DBConnection, u *model.User) (*model.User, error) {
			// コンテキストからテストの名前を取得する
			testName, ok := ctx.Value(TEST_NAME_KEY).(string)
			if !ok {
				t.Fatal("unexpected error")
			}
			if testName == "ok" {
				u.Id = 1
				return u, nil
			}
			return nil, errors.New("error")
		}

	ctx := context.Background()

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.User
		wantErr bool
	}{
		{
			"ok",
			fields{DB: moqDb, Repo: moqRepo},
			args{
				ctx:  ctx,
				form: getValidTestFormRequest(t),
			},
			getValidTestUser(t),
			false,
		},
		// フォームリクエストが不正な場合
		{
			"invalidFormRequest",
			fields{DB: moqDb, Repo: moqRepo},
			args{
				ctx:  ctx,
				form: getInvalidTestFormRequest(t),
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			us := &UserService{
				DB:   tt.fields.DB,
				Repo: tt.fields.Repo,
			}

			// コンテキストに現在のテストの名前を入れる
			ctx := context.WithValue(tt.args.ctx, TEST_NAME_KEY, tt.name)
			got, err := us.RegisterUser(ctx, tt.args.form)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserService.RegisterUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// 以下、正常系のみチェック
			if tt.name == "ok" {
				assert.Equal(t, got.Name, tt.want.Name)
				assert.Equal(t, got.UserName, tt.want.UserName)
				// パスワードの確認
				if err := bcrypt.CompareHashAndPassword([]byte(got.Password), []byte(tt.want.Password)); err != nil {
					t.Fatalf("password is wrong: %v", err)
				}
				assert.Equal(t, got.Role, tt.want.Role)
				assert.Equal(t, got.Email, tt.want.Email)
				assert.Equal(t, got.Address, tt.want.Address)
				assert.Equal(t, got.Phone, tt.want.Phone)
				assert.Equal(t, got.Website, tt.want.Website)
				assert.Equal(t, got.Company, tt.want.Company)
			}
		})
	}
}

func TestUserService_GetUser(t *testing.T) {
	type fields struct {
		DB   store.DBConnection
		Repo UserRepository
	}
	type args struct {
		ctx      context.Context
		userName string
	}

	moqDb := &DBConnectionMock{}
	moqRepo := &UserRepositoryMock{}
	moqRepo.GetUserByUserNameFunc =
		func(ctx context.Context, db store.DBConnection, userName string) (*model.User, error) {
			if userName == testutil.VALID_USER_NAME {
				return getValidTestUser(t), nil
			}
			return nil, store.ErrNotFound
		}

	ctx := context.Background()

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.User
		wantErr bool
	}{
		{
			"ok",
			fields{DB: moqDb, Repo: moqRepo},
			args{
				ctx:      ctx,
				userName: testutil.VALID_USER_NAME,
			},
			getValidTestUser(t),
			false,
		},
		// 見つからない場合
		{
			"notFound",
			fields{DB: moqDb, Repo: moqRepo},
			args{
				ctx:      ctx,
				userName: testutil.INVALID_USER_NAME, // 存在しないユーザー名
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			us := &UserService{
				DB:   tt.fields.DB,
				Repo: tt.fields.Repo,
			}
			got, err := us.GetUser(tt.args.ctx, tt.args.userName)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserService.GetUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserService.GetUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func getValidTestUser(t *testing.T) *model.User {
	t.Helper()
	u := testutil.GetTestUser(t)
	u.Id = 1
	return u
}

// テスト用FormRequest構造体を返す
func getValidTestFormRequest(t *testing.T) *model.FormRequest {
	t.Helper()

	return &model.FormRequest{
		Name:     "testUserFullName",
		UserName: "testUser",
		Password: "testPassword",
		Role:     "admin",
		Email:    "test@example.com",
		Address:  "testAddress",
		Phone:    "000-0000-0000",
		Website:  "ttp://test.com",
		Company:  "testCompany",
	}
}

func getInvalidTestFormRequest(t *testing.T) *model.FormRequest {
	t.Helper()
	fr := getValidTestFormRequest(t)
	fr.UserName = "invalidUser"
	return fr
}
