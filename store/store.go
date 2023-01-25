package store

import "github.com/kwtryo/go-sample/model"

type UserStore struct {
	LastID int
	Users  map[int]*model.User
}

func (us *UserStore) Add(u *model.User) (int, error) {
	us.LastID++
	u.Id = us.LastID
	us.Users[u.Id] = u
	return u.Id, nil
}

// ソート済のユーザー一覧を返す
func (us *UserStore) All() model.Users {
	users := make([]*model.User, len(us.Users))
	for i, u := range us.Users {
		users[i-1] = u
	}
	return users
}
