package testutil

import (
	"testing"

	"github.com/kwtryo/go-sample/model"
)

// テスト用ユーザーを返す
func GetTestUser(t *testing.T) *model.User {
	t.Helper()

	return &model.User{
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
}
