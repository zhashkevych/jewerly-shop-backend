package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

type AdminRepository struct {
	db *sqlx.DB
}

func NewAdminRepository(db *sqlx.DB) *AdminRepository {
	return &AdminRepository{
		db: db,
	}
}

func (r *AdminRepository) Authorize(email, passwordHash string) error {
	query := fmt.Sprintf("SELECT * FROM %s WHERE login=$1 AND password_hash=$2", adminUsersTable)
	_, err := r.db.Exec(query, email, passwordHash)

	return err
}
