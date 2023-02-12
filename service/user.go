package service

import (
	"context"
	"fmt"

	"github.com/kwtryo/go-sample/model"
	"github.com/kwtryo/go-sample/store"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	RegisterUser(ctx context.Context, db store.DBConnection, u *model.User) (*model.User, error)
	GetAllUsers(ctx context.Context, db store.DBConnection) ([]*model.User, error)
	GetUserByUserName(ctx context.Context, db store.DBConnection, userName string) (*model.User, error)
}

type UserService struct {
	DB   store.DBConnection
	Repo UserRepository
}

// ユーザーを登録し、登録したユーザーを返す
func (us *UserService) RegisterUser(ctx context.Context, form *model.FormRequest) (*model.User, error) {
	// パスワードをbcryptでハッシュ化
	pw, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to generate hashedpassword: %w", err)
	}

	u := &model.User{
		Name:     form.Name,
		UserName: form.UserName,
		Password: string(pw),
		Role:     form.Role,
		Email:    form.Email,
		Address:  form.Address,
		Phone:    form.Phone,
		Website:  form.Website,
		Company:  form.Company,
	}

	result, err := us.Repo.RegisterUser(ctx, us.DB, u)
	if err != nil {
		return nil, fmt.Errorf("failed to register user, username=%s: %w", u.UserName, err)
	}
	return result, nil
}

// 全てのユーザーを取得する
func (us *UserService) GetAllUsers(ctx context.Context) (model.Users, error) {
	users, err := us.Repo.GetAllUsers(ctx, us.DB)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	return users, nil
}

// ユーザーネームからユーザーを取得し、該当ユーザーを返す
func (us *UserService) GetUser(ctx context.Context, userName string) (*model.User, error) {
	user, err := us.Repo.GetUserByUserName(ctx, us.DB, userName)
	if err != nil {
		return nil, fmt.Errorf("failed to get user, username=%s: %w", userName, err)
	}
	return user, nil
}
