package store

import (
	"context"
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

func TestRegisterUser(t *testing.T) {
	ust := prepareTest(t)

	want := &model.User{
		Name:     "testUserFullName",
		UserName: "testUser",
		Password: "testPassword",
		Role:     "admin",
		Email:    "test@example.com",
		Address:  "testAddress",
		Phone:    "000-0000-0000",
		Website:  "ttp://test.com",
		Company:  "testCompany",
	}
	if err := ust.repo.RegisterUser(ust.ctx, ust.tx, want); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := ust.repo.GetUser(ust.ctx, ust.tx, want.UserName)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	t.Logf("The user ID obtained is: %d", got.Id)

	assert.Equal(t, want, got)
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

func prepareUser(ctx context.Context, t *testing.T, con Execer) *model.User {
	t.Helper()

	c := clock.FixedClocker{}
	want := &model.User{
		Name:     "testUserFullName",
		UserName: "testUser",
		Password: "testPassword",
		Role:     "admin",
		Email:    "test@example.com",
		Address:  "testAddress",
		Phone:    "000-0000-0000",
		Website:  "ttp://test.com",
		Company:  "testCompany",
		Created:  c.Now(),
		Modified: c.Now(),
	}
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
