package service

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/kwtryo/go-sample/model"
	"github.com/kwtryo/go-sample/store"
)

type UserRepository interface {
	RegisterUser(ctx context.Context, db store.Execer, u *model.User) (*model.User, error)
	GetUser(ctx context.Context, db store.Queryer, userName string) (*model.User, error)
}

type UserService struct {
	DB   *sqlx.DB
	Repo UserRepository
}

// ユーザーを登録し、登録したユーザーを返す
func (us *UserService) RegisterUser(ctx context.Context, form *model.FormRequest) (*model.User, error) {
	u := &model.User{
		Name:     form.Name,
		UserName: form.UserName,
		Password: form.Password,
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
