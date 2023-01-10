package entity

type User struct {
	Id       int
	Name     string
	UserName string
	Password string
	Role     string
	Email    string
	Address  string
	Phone    string
	Website  string
	Company  string
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
