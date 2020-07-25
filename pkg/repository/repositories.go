package repository

import (
	"github.com/jmoiron/sqlx"
	jewerly "github.com/zhashkevych/jewelry-shop-backend"
)

type User interface {
	Create(user jewerly.User) error
	GetByCredentials(email, passwordHash string) (jewerly.User, error)
	GetById(id int64) (jewerly.User, error)
	//GetAll(id int64) (jewerly.User, error)
	//Update(id int64, newUser jewerly.User) error
}

type Repository struct {
	User
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		User: NewUserRepository(db),
	}
}