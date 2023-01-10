package store

import "github.com/kwtryo/go-sample/entity"

// DBから全てのユーザー一覧を取得する
func (r *Repository) ListUsers(db Queryer) (entity.Users, error) {
	users := entity.Users{}
	sql := `SELECT
				id, name, user_name,
				role, email, address,
				phone, website, company,
				created, modified
			FROM user;`
	if err := db.Select(&users, sql); err != nil {
		return nil, err
	}
	return users, nil
}

// ユーザーをDBに登録する
func (r *Repository) AddUser(db Execer, u *entity.User) error {
	u.Created = r.Clocker.Now()
	u.Modified = r.Clocker.Now()
	sql := `INSERT INTO user (
				name, user_name, password,
				role, email, address,
				phone, website, company,
				created, modified
			)
			VALUES ( ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`
	result, err := db.Exec(
		sql,
		u.Name, u.UserName, u.Password,
		u.Role, u.Email, u.Address,
		u.Phone, u.Website, u.Company,
		u.Created, u.Modified,
	)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	u.Id = int(id)
	return nil
}
