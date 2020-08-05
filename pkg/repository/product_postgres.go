package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	jewerly "github.com/zhashkevych/jewelry-shop-backend"
)

type ProductRepository struct {
	db *sqlx.DB
}

func NewProductRepository(db *sqlx.DB) *ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

func (r *ProductRepository) Create(product jewerly.CreateProductInput) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	// insert titles
	var titleId int
	query, args := multiLanguageInsertQuery(titlesTable, product.Titles)
	row := tx.QueryRow(query, args...)

	err = row.Scan(&titleId)
	if err != nil {
		logrus.Errorf("[Create Product] create title error: %s", err.Error())
		tx.Rollback()
		return err
	}

	// insert descriptions
	var descriptionId int
	query, args = multiLanguageInsertQuery(descriptionsTable, product.Titles)
	row = tx.QueryRow(query, args...)

	err = row.Scan(&descriptionId)
	if err != nil {
		logrus.Errorf("[Create Product] create description error: %s", err.Error())
		tx.Rollback()
		return err
	}

	// insert product
	var productId int
	row = tx.QueryRow(fmt.Sprintf(`INSERT INTO %s (title_id, description_id, current_price, previous_price, code, category_id)
								VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`, productsTable),
								titleId, descriptionId, product.CurrentPrice, product.PreviousPrice, product.Code, product.CategoryId)
	err = row.Scan(&productId)
	if err != nil {
		logrus.Errorf("[Create Product] create product error: %s", err.Error())
		tx.Rollback()
		return err
	}

	// insert product images
	var imageValues string
	for i, id := range product.ImageIds {
		if i == len(product.ImageIds) - 1 {
			imageValues += fmt.Sprintf("($1, %d)", id)
		} else {
			imageValues += fmt.Sprintf("($1, %d), ", id)
		}
	}

	_, err = tx.Exec(fmt.Sprintf("INSERT INTO %s (product_id, image_id) VALUES %s", productImagesTable, imageValues), productId)
	if err != nil {
		logrus.Errorf("[Create Product] create product images error: %s", err.Error())
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func multiLanguageInsertQuery(table string, input jewerly.MultiLanguageInput) (string, []interface{}) {
	query := fmt.Sprintf("INSERT INTO %s (english, russian, ukrainian) values ($1, $2, $3) RETURNING id", table)
	args := []interface{}{input.English, input.Russian, input.Ukrainian}

	return query, args
}
