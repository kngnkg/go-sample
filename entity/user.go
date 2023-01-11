package entity

import "time"

type User struct {
	Id       int       `json:"id" db:"id"`
	Name     string    `json:"name" db:"name"`
	UserName string    `json:"username" db:"user_name"`
	Password string    `json:"password" db:"password"`
	Role     string    `json:"role" db:"role"`
	Email    string    `json:"email" db:"email"`
	Address  string    `json:"address" db:"address"`
	Phone    string    `json:"phone" db:"phone"`
	Website  string    `json:"website" db:"website"`
	Company  string    `json:"company" db:"company"`
	Created  time.Time `json:"created" db:"created"`
	Modified time.Time `json:"modified" db:"modified"`
}

type Users []*User

// Returns all users
func AllUser() (users []User, err error) {
	user1 := User{
		Id:       1,
		Name:     "Leanne Graham",
		UserName: "Bret",
		Email:    "Sincere@april.biz",
		Address:  "test address",
		Phone:    "1-770-736-8031 x56442",
		Website:  "hildegard.org",
		Company:  "test company",
	}
	users = append(users, user1)
	user2 := User{
		Id:       2,
		Name:     "Leanne Graham2",
		UserName: "Bret2",
		Email:    "Sincere@april.biz",
		Address:  "test address",
		Phone:    "1-770-736-8031 x56442",
		Website:  "hildegard.org",
		Company:  "test company",
	}
	users = append(users, user2)
	return
}
