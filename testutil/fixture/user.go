package fixture

import (
	"log"
	"math/rand"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/kwtryo/go-sample/model"
	"golang.org/x/crypto/bcrypt"
)

// Passwordを設定した場合、そのパスワードを元としてハッシュ化された値が設定される。
func User(u *model.User) *model.User {
	pw := "testPassword"
	if u.Password != "" {
		pw = u.Password
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("unexpected error: %v", err)
	}
	random := strconv.Itoa(rand.Int())[:3]

	result := &model.User{
		Id:       rand.Int(),
		Name:     "testUserFullName" + random,
		UserName: "testUserName" + random,
		Password: string(hashed),
		Role:     "admin",
		Email:    "test@example.com",
		Address:  "testAddress",
		Phone:    "000-0000-0000",
		Website:  "ttp://test.com",
		Company:  "testCompany",
		Created:  time.Now(),
		Modified: time.Now(),
	}

	if u == nil {
		return result
	}
	if u.Id != 0 {
		result.Id = u.Id
	}
	if u.Name != "" {
		result.Name = u.Name
	}
	if u.UserName != "" {
		result.UserName = u.UserName
	}
	// Passwordは設定済み
	if u.Role != "" {
		result.Role = u.Role
	}
	if u.Email != "" {
		result.Email = u.Email
	}
	if u.Address != "" {
		result.Address = u.Address
	}
	if u.Phone != "" {
		result.Phone = u.Phone
	}
	if u.Website != "" {
		result.Website = u.Website
	}
	if u.Company != "" {
		result.Company = u.Company
	}
	if !u.Created.IsZero() {
		result.Created = u.Created
	}
	if !u.Modified.IsZero() {
		result.Modified = u.Modified
	}
	return result
}

// ログイン時に送信されるボディ
func LoginFormBody(l *model.Login) *strings.Reader {
	random := strconv.Itoa(rand.Int())[:5]
	login := model.Login{
		Username: "testUserName" + random,
		Password: "testPassword",
	}
	if l.Username != "" {
		login.Username = l.Username
	}
	if l.Password != "" {
		login.Password = l.Password
	}
	form := url.Values{}
	form.Add("username", l.Username)
	form.Add("password", l.Password)
	body := strings.NewReader(form.Encode())
	return body
}

func UserFormRequest(u *model.User) *model.FormRequest {
	user := User(u)
	return &model.FormRequest{
		Name:     user.Name,
		UserName: user.UserName,
		Password: u.Password, // ハッシュ化しない
		Role:     user.Role,
		Email:    user.Email,
		Address:  user.Address,
		Phone:    user.Phone,
		Website:  user.Website,
		Company:  user.Company,
	}
}

// ユーザー登録時に送信されるボディ
func RegisterUserBody(u *model.User) url.Values {
	user := User(u)

	form := url.Values{}
	form.Add("name", user.Name)
	form.Add("username", user.UserName)
	form.Add("password", user.Password)
	form.Add("role", user.Role)
	form.Add("email", user.Email)
	form.Add("address", user.Address)
	form.Add("phone", user.Phone)
	form.Add("website", user.Website)
	form.Add("company", user.Company)
	return form
}
