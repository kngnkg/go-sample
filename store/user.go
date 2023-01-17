package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/kwtryo/go-sample/model"
)

// ユーザーをDBに登録し、登録したユーザーを返す。
func (r *Repository) RegisterUser(ctx context.Context, db Execer, u *model.User) (*model.User, error) {
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
		var mysqlError *mysql.MySQLError
		if errors.As(err, &mysqlError) && mysqlError.Number == ErrCodeMySQLDuplicateEntry {
			err = fmt.Errorf("cannot create same name user: %w", ErrAlreadyEntry)
			return nil, err
		}
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	u.Id = int(id)
	return u, nil
}

// ユーザーネームからユーザーを取得する。
func (r *Repository) GetUser(ctx context.Context, db Queryer, userName string) (*model.User, error) {
	u := &model.User{}
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

// DBのユーザーを全て削除する
func (r *Repository) DeleteUserAll(ctx context.Context, db Execer) error {
	sql := `DELETE FROM user;`
	if _, err := db.ExecContext(ctx, sql); err != nil {
		return err
	}
	return nil
}
