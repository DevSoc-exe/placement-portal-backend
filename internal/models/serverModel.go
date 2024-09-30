package models

type Store interface {
	CreateUser(user *User) error
	VerifyUser(userId string, token string) error
	UpdateUserRefreshToken(refresh_token, userId string) error
	RevokeUserRefreshToken(refresh_token string) error

	GetUserByEmail(email string) (*User, error)
	GetUserByID(id string) (*User, error)
}
