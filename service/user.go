package service

import (
	"context"

	"github.com/kwtryo/go-sample/model"
	"github.com/kwtryo/go-sample/store"
	"golang.org/x/crypto/bcrypt"
)

//go:generate go run github.com/matryer/moq -out moq_test.go . UserRepository
type UserRepository interface {
	RegisterUser(ctx context.Context, db store.Execer, u *model.User) (*model.User, error)
	GetUser(ctx context.Context, db store.Queryer, userName string) (*model.User, error)
}

type UserService struct {
	DB   store.Execer
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

func (us *UserService) GetUser(ctx context.Context, db store.Queryer, userName string) (*model.User, error) {
	return nil, nil
}
