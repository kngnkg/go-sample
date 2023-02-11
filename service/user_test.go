package service

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/kwtryo/go-sample/model"
	"github.com/kwtryo/go-sample/store"
	"github.com/kwtryo/go-sample/testutil/fixture"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
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

	userName := "testUserName"
	password := "password"
	user := fixture.User(&model.User{
		UserName: userName,
		Password: password,
	})
	user.Password = password
	userForm := fixture.UserFormRequest(user)

	moqErr := errors.New("err from mock")
	moqDb := &DBConnectionMock{}
	moqRepo := &UserRepositoryMock{}
	moqRepo.RegisterUserFunc =
		func(ctx context.Context, db store.DBConnection, u *model.User) (*model.User, error) {
			if u.UserName != userName {
				return nil, moqErr
			}
			u.Id = user.Id
			return u, nil
		}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.User
		wantErr error
	}{
		{
			"ok",
			fields{DB: moqDb, Repo: moqRepo},
			args{
				context.TODO(),
				userForm,
			},
			user,
			nil,
		},
		// ユーザー登録に失敗
		{
			"failedToRegister",
			fields{DB: moqDb, Repo: moqRepo},
			args{
				context.TODO(),
				fixture.UserFormRequest(&model.User{
					UserName: "failedToRegister",
				}),
			},
			nil,
			moqErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			us := &UserService{
				DB:   tt.fields.DB,
				Repo: tt.fields.Repo,
			}
			got, err := us.RegisterUser(tt.args.ctx, tt.args.form)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("want err: %v but got: %v", tt.wantErr, err)
			}
			if tt.name != "ok" {
				return
			}

			assert.Equal(t, got.Name, tt.want.Name)
			assert.Equal(t, got.UserName, tt.want.UserName)
			// パスワードの確認
			if err := bcrypt.CompareHashAndPassword([]byte(got.Password), []byte(password)); err != nil {
				t.Fatalf("password is wrong: %v", err)
			}
			assert.Equal(t, got.Role, tt.want.Role)
			assert.Equal(t, got.Email, tt.want.Email)
			assert.Equal(t, got.Address, tt.want.Address)
			assert.Equal(t, got.Phone, tt.want.Phone)
			assert.Equal(t, got.Website, tt.want.Website)
			assert.Equal(t, got.Company, tt.want.Company)
		})
	}
}

func TestUserService_GetAllUsers(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	// ランダムなユーザーを5人生成する
	users := model.Users{}
	for i := 0; i < 5; i++ {
		user := fixture.User(&model.User{})
		users = append(users, user)
	}

	moqErr := errors.New("err from mock")
	tests := []struct {
		name    string
		args    args
		want    model.Users
		wantErr error
	}{
		// 正常系
		{
			"ok",
			args{ctx: context.TODO()},
			users,
			nil,
		},
		// 取得に失敗
		{
			"failedToGetUsers",
			args{ctx: context.TODO()},
			nil,
			moqErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			moqDb := &DBConnectionMock{}
			moqRepo := &UserRepositoryMock{}
			moqRepo.GetAllUsersFunc =
				func(ctx context.Context, db store.DBConnection) ([]*model.User, error) {
					if tt.name == "failedToGetUsers" {
						return nil, moqErr
					}
					return users, nil
				}
			us := &UserService{
				DB:   moqDb,
				Repo: moqRepo,
			}

			got, err := us.GetAllUsers(tt.args.ctx)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("want err: %v but got: %v", tt.wantErr, err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserService.GetAllUsers() = %v, want %v", got, tt.want)
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
	search := "testUserName"
	user := fixture.User(&model.User{
		UserName: search,
	})

	moqErr := errors.New("err from mock")
	moqDb := &DBConnectionMock{}
	moqRepo := &UserRepositoryMock{}
	moqRepo.GetUserByUserNameFunc =
		func(ctx context.Context, db store.DBConnection, userName string) (*model.User, error) {
			if userName != search {
				return nil, moqErr
			}
			return user, nil
		}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.User
		wantErr error
	}{
		{
			"ok",
			fields{DB: moqDb, Repo: moqRepo},
			args{
				context.TODO(),
				search,
			},
			user,
			nil,
		},
		// 取得に失敗
		{
			"failedToGetUser",
			fields{DB: moqDb, Repo: moqRepo},
			args{
				context.TODO(),
				"invalid",
			},
			nil,
			moqErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			us := &UserService{
				DB:   tt.fields.DB,
				Repo: tt.fields.Repo,
			}
			got, err := us.GetUser(tt.args.ctx, tt.args.userName)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("want err: %v but got: %v", tt.wantErr, err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserService.GetUser() = %v, want %v", got, tt.want)
			}
		})
	}
}
