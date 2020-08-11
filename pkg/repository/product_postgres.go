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

	// insert product
	var productId int
	row := tx.QueryRow(fmt.Sprintf(`INSERT INTO %s 
								(current_price, previous_price, code, category_id)
								VALUES ($1, $2, $3, $4) RETURNING id`, productsTable),
		product.CurrentPrice, product.PreviousPrice, product.Code, product.CategoryId)
	err = row.Scan(&productId)
	if err != nil {
		logrus.Errorf("[Create Product] create product error: %s", err.Error())
		tx.Rollback()
		return err
	}

	// insert titles
	query, args := multiLanguageInsertQuery(titlesTable, product.Titles, productId)
	_, err = tx.Exec(query, args...)
	if err != nil {
		logrus.Errorf("[Create Product] create title error: %s", err.Error())
		tx.Rollback()
		return err
	}

	// insert descriptions
	query, args = multiLanguageInsertQuery(descriptionsTable, product.Descriptions, productId)
	_, err = tx.Exec(query, args...)
	if err != nil {
		logrus.Errorf("[Create Product] create description error: %s", err.Error())
		tx.Rollback()
		return err
	}

	// insert meterial
	query, args = multiLanguageInsertQuery(materialsTable, product.Material, productId)
	_, err = tx.Exec(query, args...)
	if err != nil {
		logrus.Errorf("[Create Product] create materials error: %s", err.Error())
		tx.Rollback()
		return err
	}

	// insert product images
	var imageValues string
	for i, id := range product.ImageIds {
		if i == len(product.ImageIds)-1 {
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

func multiLanguageInsertQuery(table string, input jewerly.MultiLanguageInput, productId int) (string, []interface{}) {
	query := fmt.Sprintf("INSERT INTO %s (english, russian, ukrainian, product_id) values ($1, $2, $3, $4)", table)
	args := []interface{}{input.English, input.Russian, input.Ukrainian, productId}

	return query, args
}

func (r *ProductRepository) GetAll(filters jewerly.GetAllProductsFilters) (jewerly.ProductsList, error) {
	var products jewerly.ProductsList

	selectQuery := fmt.Sprintf(`SELECT p.id, t.%[1]s as title, d.%[1]s as description, m.%[1]s as material, p.current_price,
							p.previous_price, p.code, p.category_id, p.in_stock`, filters.Language)
	fromQuery := fmt.Sprintf(` FROM %[1]s p
							JOIN %[2]s t on t.product_id = p.id
							JOIN %[3]s d on d.product_id = p.id
							JOIN %[4]s m on m.product_id = p.id`, productsTable, titlesTable, descriptionsTable, materialsTable)

	// build where query
	var whereQuery string

	argId := 1
	args := make([]interface{}, 0)
	if filters.CategoryId.Valid {
		whereQuery = fmt.Sprintf("WHERE p.category_id=$%d", argId)
		args = append(args, filters.CategoryId)
		argId++
	}

	args = append(args, filters.Offset, filters.Limit)
	limitQuery := fmt.Sprintf(" OFFSET $%d LIMIT $%d", argId, argId+1)

	// BUILD FINAL QUERY
	var query string
	if whereQuery == "" {
		query = fmt.Sprintf("%s %s %s", selectQuery, fromQuery, limitQuery)
	} else {
		query = fmt.Sprintf("%s %s %s %s", selectQuery, fromQuery, whereQuery, limitQuery)
	}

	// select products
	err := r.db.Select(&products.Products, query, args...)

	// total count
	err = r.db.Get(&products.Total, fmt.Sprintf("SELECT count(*) %s", fromQuery))

	return products, err
}

func (r *ProductRepository) Delete(id int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(fmt.Sprintf("DELETE FROM %s WHERE product_id=$1", productImagesTable), id)
	if err != nil {
		logrus.Errorf("[Delete Product] delete images error: %s", err.Error())
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(fmt.Sprintf("DELETE FROM %s WHERE id=$1", productsTable), id)
	if err != nil {
		logrus.Errorf("[Delete Product] delete product error: %s", err.Error())
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *ProductRepository) GetById(id int, language string) (jewerly.ProductResponse, error) {
	var product jewerly.ProductResponse

	query := fmt.Sprintf(`SELECT p.id, t.%[1]s as title, d.%[1]s as description, m.%[1]s as material, p.current_price,
							p.previous_price, p.code, p.category_id, p.in_stock FROM %[2]s p
							JOIN %[3]s t on t.product_id = p.id
							JOIN %[4]s d on d.product_id = p.id
							JOIN %[5]s m on m.product_id = p.id WHERE p.id = $1`, language, productsTable, titlesTable, descriptionsTable, materialsTable)
	err := r.db.Get(&product, query, id)

	return product, err
}

func (r *ProductRepository) CreateImage(url, altText string) (int, error) {
	var id int

	row := r.db.QueryRow(fmt.Sprintf("INSERT INTO %s (url, alt_text) values ($1, $2) RETURNING id", imagesTable), url, altText)
	err := row.Scan(&id)
	return id, err
}

func (r *ProductRepository) GetProductImages(productId int) ([]jewerly.Image, error) {
	var images []jewerly.Image

	err := r.db.Select(&images, fmt.Sprintf("SELECT i.id, i.url, i.alt_text FROM %s i JOIN %s pi ON pi.image_id = i.id WHERE pi.product_id = $1",
		imagesTable, productImagesTable), productId)

	return images, err
}

func (r *ProductRepository) Update(id int, inp jewerly.UpdateProductInput) error {
	//tx, err := r.db.Begin()
	//if err != nil {
	//	return err
	//}
	//
	//if inp.Material != nil {
	//	tx.Exec(fmt.Sprintf("UPDATE %s SET "))
	//}
	return nil
}

func multiLanguageUpdateQuery(table string, input jewerly.MultiLanguageInput) (string, []interface{}) {
	query := fmt.Sprintf("UPDATE %s SET english=$1, russian=$2, ukrainian=$3 JOIN products WHERE products.mate ", table)
	args := []interface{}{input.English, input.Russian, input.Ukrainian}

	return query, args
}
