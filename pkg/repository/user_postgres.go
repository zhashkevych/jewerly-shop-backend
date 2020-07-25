package repository

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	jewerly "github.com/zhashkevych/jewelry-shop-backend"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Create(user jewerly.User) error {
	_, err := r.db.Exec(fmt.Sprintf("INSERT INTO %s (first_name, last_name, email, password_hash) VALUES ($1, $2, $3, $4)",
		userTableName), user.FirstName, user.LastName, user.Email, user.PasswordHash)
	return err
}

func (r *UserRepository) GetByCredentials(email, passwordHash string) (jewerly.User, error) {
	var user jewerly.User
	query := fmt.Sprintf("SELECT * FROM %s WHERE email=$1 AND password_hash=$2", userTableName)
	err := r.db.Get(&user, query, email, passwordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return jewerly.User{}, jewerly.ErrUserNotFound
		}

		return jewerly.User{}, err
	}

	return user, nil
}

func (r *UserRepository) GetById(id int64) (jewerly.User, error) {
	var user jewerly.User
	query := fmt.Sprintf("SELECT * FROM %s WHERE id=$1", userTableName)
	err := r.db.Get(&user, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return jewerly.User{}, jewerly.ErrUserNotFound
		}

		return jewerly.User{}, err
	}

	return user, nil
}
