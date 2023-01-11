package store

import (
	"context"

	"github.com/kwtryo/go-sample/entity"
)

// ユーザーをDBに登録する
func (r *Repository) RegisterUser(ctx context.Context, db Execer, u *entity.User) error {
	u.Created = r.Clocker.Now()
	u.Modified = r.Clocker.Now()
	sql := `INSERT INTO user (
				name, user_name, password,
				role, email, address,
				phone, website, company,
				created, modified
			)
			VALUES (
				?, ?, ?,
				?, ?, ?,
				?, ?, ?,
				?, ?);`
	result, err := db.ExecContext(
		ctx, sql,
		u.Name, u.UserName, u.Password,
		u.Role, u.Email, u.Address,
		u.Phone, u.Website, u.Company,
		u.Created, u.Modified,
	)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	u.Id = int(id)
	return nil
}

// DBからユーザーを取得する
func (r *Repository) GetUser(ctx context.Context, db Queryer, userName string) (*entity.User, error) {
	u := &entity.User{}
	sql := `SELECT
				id, name, user_name,
				password, role, email,
				address, phone, website,
				company, created, modified
			FROM user
			WHERE user_name = ?;`
	if err := db.GetContext(ctx, u, sql, userName); err != nil {
		return nil, err
	}
	return u, nil
}