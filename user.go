package jewerly

import (
	"errors"
	"time"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type User struct {
	Id           int64     `db:"id"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
	FirstName    string    `db:"first_name"`
	LastName     string    `db:"last_name"`
	RegisteredAt time.Time `db:"registered_at"`
}

type AdminUser struct {
	Id           int64  `db:"id"`
	Login        string `db:"login"`
	PasswordHash string `db:"password_hash"`
}
