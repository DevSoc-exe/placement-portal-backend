package models

type Store interface {
	CreateUser(user *User)	error
	GetUserByEmail(email string) (*User, error)
	VerifyUser(userId string, token string) error
}
