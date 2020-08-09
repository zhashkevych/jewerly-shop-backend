package repository

import (
	//"fmt"
	"github.com/jmoiron/sqlx"
)

type OrderRepository struct {
	db *sqlx.DB
}

func NewOrderRepository(db *sqlx.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(userId int64, productIds []int) error {
	//insertValuesQuery := ""
	//insertValues := make([]interface{}, len(productIds))
	//
	//for i, id := range productIds {
	//	if i == len(productIds) - 1 {
	//		insertValues += fmt.Sprintf("($1, $%d)", i+1)
	//	} else {
	//		insertValues += fmt.Sprintf("($1, $%d), ", i+1)
	//	}
	//
	//	insertValues[i] = id
	//}
	//
	//r.db.Exec(fmt.Sprintf("INSERT INTO %s"))
	return nil
}

