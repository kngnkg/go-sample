package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/kwtryo/go-sample/model"
)

// ユーザーをDBに登録し、登録したユーザーを返す。
func (r *Repository) RegisterUser(ctx context.Context, db DBConnection, u *model.User) (*model.User, error) {
	u.Created = r.Clocker.Now()
	u.Modified = r.Clocker.Now()
	query := `INSERT INTO user (
				name, user_name, password, role, email, address,
				phone, website, company, created, modified
			)
			VALUES (
				?, ?, ?, ?, ?, ?,
				?, ?, ?, ?, ?);`
	result, err := db.ExecContext(
		ctx, query,
		u.Name, u.UserName, u.Password,
		u.Role, u.Email, u.Address,
		u.Phone, u.Website, u.Company,
		u.Created, u.Modified,
	)
	if err != nil {
		var mysqlError *mysql.MySQLError
		if errors.As(err, &mysqlError) && mysqlError.Number == ErrCodeMySQLDuplicateEntry {
			err = fmt.Errorf("cannot create same name user: %w", ErrAlreadyEntry)
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
func (r *Repository) GetUser(ctx context.Context, db DBConnection, userName string) (*model.User, error) {
	u := &model.User{}
	query := `SELECT
				id, name, user_name,
				password, role, email,
				address, phone, website,
				company, created, modified
			FROM user
			WHERE user_name = ?;`
	if err := db.GetContext(ctx, u, query, userName); err != nil {
		if err == sql.ErrNoRows {
			err = fmt.Errorf("user not found: %w", ErrNotFound)
		}
		return nil, err
	}
	return u, nil
}

// DBのユーザーを全て削除する
func (r *Repository) DeleteUserAll(ctx context.Context, db DBConnection) error {
	query := `DELETE FROM user;`
	if _, err := db.ExecContext(ctx, query); err != nil {
		return err
	}
	return nil
}
