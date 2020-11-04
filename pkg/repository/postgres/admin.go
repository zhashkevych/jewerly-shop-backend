package postgres

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	jewerly "github.com/zhashkevych/jewelry-shop-backend"
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
	var admin jewerly.AdminUser
	query := fmt.Sprintf("SELECT * FROM %s WHERE login=$1 AND password_hash=$2", adminUsersTable)
	err := r.db.Get(&admin, query, email, passwordHash)

	return err
}
