package service

import (
	"context"

	"github.com/kwtryo/go-sample/model"
	"github.com/kwtryo/go-sample/store"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	RegisterUser(ctx context.Context, db store.DBConnection, u *model.User) (*model.User, error)
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
		return nil, err
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
		return nil, err
	}
	return result, nil
}

// ユーザーネームからユーザーを取得し、該当ユーザーを返す
func (us *UserService) GetUser(ctx context.Context, userName string) (*model.User, error) {
	user, err := us.Repo.GetUserByUserName(ctx, us.DB, userName)
	if err != nil {
		return nil, err
	}
	return user, nil
}
