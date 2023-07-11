package testutil

import (
	"context"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/kngnkg/go-sample/clock"
	"github.com/kngnkg/go-sample/model"
)

// テスト用のユーザを登録する。
func PrepareUser(ctx context.Context, t *testing.T, c clock.Clocker, tx *sqlx.Tx, user *model.User) *model.User {
	t.Helper()

	now := c.Now()
	user.Created = now
	user.Modified = now

	result, err := tx.ExecContext(
		ctx,
		`INSERT INTO user (
			name, user_name, password, role, email, address,
			phone, website, company, created, modified
		)
		VALUES (
			?, ?, ?, ?, ?, ?,
			?, ?, ?, ?, ?);`,
		user.Name, user.UserName, user.Password,
		user.Role, user.Email, user.Address,
		user.Phone, user.Website, user.Company,
		user.Created, user.Modified,
	)
	if err != nil {
		t.Fatal(err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}
	user.Id = int(id)
	return user
}

// DBのユーザーを全て削除する
func DeleteUserAll(ctx context.Context, t *testing.T, tx *sqlx.Tx) {
	query := `DELETE FROM user;`
	if _, err := tx.ExecContext(ctx, query); err != nil {
		t.Errorf("cannot delete users: %v", err)
	}
}
