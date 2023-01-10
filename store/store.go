package store

import "github.com/kwtryo/go-sample/entity"

type UserStore struct {
	LastID int
	Users  map[int]*entity.User
}

func (us *UserStore) Add(u *entity.User) (int, error) {
	us.LastID++
	u.Id = us.LastID
	us.Users[u.Id] = u
	return u.Id, nil
}

// ソート済のユーザー一覧を返す
func (us *UserStore) All() entity.Users {
	users := make([]*entity.User, len(us.Users))
	for i, u := range us.Users {
		users[i-1] = u
	}
	return users
}
