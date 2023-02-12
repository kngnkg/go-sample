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
			err = fmt.Errorf("cannot create same name user, username=%s: %w", u.UserName, ErrAlreadyEntry)
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
func (r *Repository) GetUserByUserName(ctx context.Context, db DBConnection, userName string) (*model.User, error) {
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
			err = fmt.Errorf("user not found, username=%s: %w", userName, ErrNotFound)
		}
		return nil, err
	}
	return u, nil
}

// 全ユーザーを取得する。
func (r *Repository) GetAllUsers(ctx context.Context, db DBConnection) ([]*model.User, error) {
	users := model.Users{}
	query := `SELECT
				id, name, user_name,
				role, email, address,
				phone, website, company,
				created, modified
			FROM user;`
	if err := db.SelectContext(ctx, &users, query); err != nil {
		return nil, err
	}
	// ユーザーが一人も存在しない場合
	if len(users) == 0 {
		return nil, fmt.Errorf("none of the users exist: %w", ErrNotFound)
	}
	return users, nil
}
