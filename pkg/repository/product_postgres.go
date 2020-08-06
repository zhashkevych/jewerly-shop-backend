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

	// insert meterial
	var materialId int
	query, args = multiLanguageInsertQuery(materialsTable, product.Material)
	row = tx.QueryRow(query, args...)

	err = row.Scan(&materialId)
	if err != nil {
		logrus.Errorf("[Create Product] create materials error: %s", err.Error())
		tx.Rollback()
		return err
	}

	// insert product
	var productId int
	row = tx.QueryRow(fmt.Sprintf(`INSERT INTO %s 
								(title_id, description_id, material_id, current_price, previous_price, code, category_id)
								VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`, productsTable),
		titleId, descriptionId, materialId, product.CurrentPrice, product.PreviousPrice, product.Code, product.CategoryId)
	err = row.Scan(&productId)
	if err != nil {
		logrus.Errorf("[Create Product] create product error: %s", err.Error())
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

func multiLanguageInsertQuery(table string, input jewerly.MultiLanguageInput) (string, []interface{}) {
	query := fmt.Sprintf("INSERT INTO %s (english, russian, ukrainian) values ($1, $2, $3) RETURNING id", table)
	args := []interface{}{input.English, input.Russian, input.Ukrainian}

	return query, args
}

func (r *ProductRepository) GetAll(filters jewerly.GetAllProductsFilters) (jewerly.ProductsList, error) {
	var products jewerly.ProductsList

	query := fmt.Sprintf(`SELECT p.id, t.%[1]s as title, d.%[1]s as description, m.%[1]s as material, p.current_price,
							p.previous_price, p.code, p.category_id FROM %[2]s p
							JOIN %[3]s t on t.id = p.title_id
							JOIN %[4]s d on d.id = p.description_id
							JOIN %[5]s m on m.id = p.material_id`, filters.Language, productsTable, titlesTable, descriptionsTable, materialsTable)
	err := r.db.Select(&products.Products, query)

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
							p.previous_price, p.code, p.category_id FROM %[2]s p
							JOIN %[3]s t on t.id = p.title_id
							JOIN %[4]s d on d.id = p.description_id
							JOIN %[5]s m on m.id = p.material_id WHERE p.id = $1`, language, productsTable, titlesTable, descriptionsTable, materialsTable)
	err := r.db.Get(&product, query, id)

	return product, err
}

func (r *ProductRepository) CreateImage(url, altText string) (int, error) {
	var id int

	row := r.db.QueryRow(fmt.Sprintf("INSERT INTO %s (url, alt_text) values ($1, $2) RETURNING id", imagesTable), url, altText)
	err := row.Scan(&id)
	return id, err
}
