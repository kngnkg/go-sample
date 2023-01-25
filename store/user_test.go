package store

import (
	"context"
	"fmt"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/kwtryo/go-sample/clock"
	"github.com/kwtryo/go-sample/model"
	"github.com/kwtryo/go-sample/testutil"
	"github.com/stretchr/testify/assert"
)

type userStoreTest struct {
	ctx  context.Context
	tx   *sqlx.Tx
	repo *Repository
}

func prepareTest(t *testing.T) *userStoreTest {
	t.Helper()

	ctx := context.Background()
	tx, err := testutil.OpenDbForTest(t).BeginTxx(ctx, nil)
	t.Cleanup(func() { _ = tx.Rollback() })
	if err != nil {
		t.Fatal(err)
	}
	repo := &Repository{
		Clocker: clock.FixedClocker{},
	}

	ust := &userStoreTest{
		ctx:  ctx,
		tx:   tx,
		repo: repo,
	}
	if err := ust.repo.DeleteUserAll(ust.ctx, ust.tx); err != nil {
		t.Logf("failed to initialize task: %v", err)
	}
	return ust
}

// TODO: 異常系もかく
func TestRegisterUser(t *testing.T) {
	type want struct {
		err error
	}
	type test struct {
		// 登録するユーザー
		user *model.User
		want want
	}

	tests := map[string]test{
		// 正常系
		"ok": {
			user: getTestUser(),
			want: want{
				err: nil,
			},
		},
		// 既に登録されていた場合
		"errAlreadyEntry": {
			user: getTestUser(),
			want: want{
				err: fmt.Errorf("cannot create same name user: %w", ErrAlreadyEntry),
			},
		},
	}

	for n, tst := range tests {
		t.Run(n, func(t *testing.T) {
			tstName := n
			tst := tst
			// CIワークフローでデッドロックが起こるので、暫定策としてコメントアウト
			// t.Parallel()

			ust := prepareTest(t)

			registeredUser, err := ust.repo.RegisterUser(ust.ctx, ust.tx, tst.user)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tstName == "errAlreadyEntry" {
				// 異常系のテストの場合のみ再度登録する
				_, err = ust.repo.RegisterUser(ust.ctx, ust.tx, tst.user)
				assert.Equal(t, tst.want.err, err)
			} else {
				// 正常系
				got, err := ust.repo.GetUser(ust.ctx, ust.tx, tst.user.UserName)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}

				t.Logf("The user ID obtained is: %d", got.Id)

				assert.Equal(t, registeredUser, got)
			}
		})
	}
}

func TestGetUser(t *testing.T) {
	ust := prepareTest(t)

	want := prepareUser(ust.ctx, t, ust.tx)
	got, err := ust.repo.GetUser(ust.ctx, ust.tx, want.UserName)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	t.Logf("The user ID obtained is: %d", got.Id)
	assert.Equal(t, want, got)
}

// テストに使用するユーザーを返す
func getTestUser() *model.User {
	return &model.User{
		Name:     "testUserFullName",
		UserName: "testUser",
		Password: "hashedTestPassword",
		Role:     "admin",
		Email:    "test@example.com",
		Address:  "testAddress",
		Phone:    "000-0000-0000",
		Website:  "ttp://test.com",
		Company:  "testCompany",
	}
}

func prepareUser(ctx context.Context, t *testing.T, con DBConnection) *model.User {
	t.Helper()

	c := clock.FixedClocker{}
	want := getTestUser()
	now := c.Now()
	want.Created = now
	want.Modified = now

	result, err := con.ExecContext(
		ctx,
		`INSERT INTO user (
			name, user_name, password,
			role, email, address,
			phone, website, company,
			created, modified
		)
		VALUES (
			?, ?, ?,
			?, ?, ?,
			?, ?, ?,
			?, ?
		);`,
		want.Name, want.UserName, want.Password,
		want.Role, want.Email, want.Address,
		want.Phone, want.Website, want.Company,
		want.Created, want.Modified,
	)
	if err != nil {
		t.Fatal(err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}
	want.Id = int(id)
	return want
}
