package model

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

// リクエストをバインドする構造体
type FormRequest struct {
	Name     string `form:"name" binding:"required"`
	UserName string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
	Role     string `form:"role" binding:"required"`
	Email    string `form:"email" binding:"required"`
	Address  string `form:"address" binding:"required"`
	Phone    string `form:"phone" binding:"required"`
	Website  string `form:"website" binding:"required"`
	Company  string `form:"company" binding:"required"`
}

// ログイン時のリクエストをバインドする
type Login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}
