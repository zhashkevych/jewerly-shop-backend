package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	jewerly "github.com/zhashkevych/jewelry-shop-backend"
	"strings"
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
	row.Scan(&titleId)
	if err != nil {
		logrus.Errorf("[Create Product] create title error: %s", err.Error())
		tx.Rollback()
		return err
	}

	// insert descriptions
	var descriptionId int
	query, args = multiLanguageInsertQuery(descriptionsTable, product.Descriptions)
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
	row.Scan(&materialId)
	if err != nil {
		logrus.Errorf("[Create Product] create materials error: %s", err.Error())
		tx.Rollback()
		return err
	}

	// insert prices
	var currentPriceId int
	query, args = multiCurrencyInsertQuery(product.Price)
	row = tx.QueryRow(query, args...)
	row.Scan(&currentPriceId)
	if err != nil {
		logrus.Errorf("[Create Product] create materials error: %s", err.Error())
		tx.Rollback()
		return err
	}

	//insert product
	var productId int
	row = tx.QueryRow(fmt.Sprintf(`INSERT INTO %s
								(code, category_id, title_id, description_id, material_id, price_id)
								VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`, productsTable),
		product.Code, product.CategoryId, titleId, descriptionId, materialId, currentPriceId)
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

func multiCurrencyInsertQuery(input jewerly.MultiCurrencyInput) (string, []interface{}) {
	query := fmt.Sprintf("INSERT INTO %s (usd, eur, uah) values ($1, $2, $3) RETURNING id", pricesTable)
	args := []interface{}{input.USD, input.EUR, input.UAH}

	return query, args
}

func (r *ProductRepository) GetAll(filters jewerly.GetAllProductsFilters) (jewerly.ProductsList, error) {
	var products jewerly.ProductsList

	selectQuery := fmt.Sprintf(`SELECT p.id, t.%[1]s as title, d.%[1]s as description, m.%[1]s as material, pr.%[2]s as price,
							p.code, p.category_id, p.in_stock`, filters.Language, filters.Currency)
	fromQuery := fmt.Sprintf(` FROM %[1]s p
							JOIN %[2]s t on t.id = p.title_id
							JOIN %[3]s d on d.id = p.description_id
							JOIN %[4]s m on m.id = p.material_id
							JOIN %[5]s pr on pr.id = p.price_id`,
		productsTable, titlesTable, descriptionsTable, materialsTable, pricesTable)

	// build where query
	var whereQuery string

	argId := 1
	args := make([]interface{}, 0)
	if filters.CategoryId.Valid {
		whereQuery = fmt.Sprintf("WHERE p.category_id=$%d", argId)
		args = append(args, filters.CategoryId.Int64)
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
	var countQuery string
	args = make([]interface{}, 0)
	if whereQuery == "" {
		countQuery = fmt.Sprintf("SELECT count(*) %s", fromQuery)
	} else {
		countQuery = fmt.Sprintf("SELECT count(*) %s %s", fromQuery, whereQuery)
		args = append(args, filters.CategoryId.Int64)
	}

	err = r.db.Get(&products.Total, countQuery, args...)

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

func (r *ProductRepository) GetById(id int, language, currency string) (jewerly.ProductResponse, error) {
	var product jewerly.ProductResponse

	query := fmt.Sprintf(`SELECT p.id, t.%[1]s as title, d.%[1]s as description, m.%[1]s as material, 
							pr.%[2]s as price, p.code, p.category_id, p.in_stock FROM %[3]s p
							JOIN %[4]s t on t.id = p.title_id
							JOIN %[5]s d on d.id = p.description_id
							JOIN %[6]s m on m.id = p.material_id
							JOIN %[7]s pr on pr.id = p.price_id WHERE p.id = $1`,
		language, currency, productsTable, titlesTable, descriptionsTable, materialsTable, pricesTable)
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
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	if inp.Material != nil {
		query, args := multiLanguageUpdateQuery(materialsTable, *inp.Material, id)

		if _, err := r.db.Exec(query, args...); err != nil {
			logrus.Errorf("[Update Product] update materials error: %s", err.Error())
			tx.Rollback()
			return err
		}
	}

	if inp.Titles != nil {
		query, args := multiLanguageUpdateQuery(titlesTable, *inp.Titles, id)

		if _, err := r.db.Exec(query, args...); err != nil {
			logrus.Errorf("[Update Product] update titles error: %s", err.Error())
			tx.Rollback()
			return err
		}
	}

	if inp.Descriptions != nil {
		query, args := multiLanguageUpdateQuery(descriptionsTable, *inp.Descriptions, id)

		if _, err := r.db.Exec(query, args...); err != nil {
			logrus.Errorf("[Update Product] update titles error: %s", err.Error())
			tx.Rollback()
			return err
		}
	}

	if inp.Price != nil {
		query, args := multiCurrencyUpdateQuery(*inp.Price, id)

		if _, err := r.db.Exec(query, args...); err != nil {
			logrus.Errorf("[Update Product] update price error: %s", err.Error())
			tx.Rollback()
			return err
		}
	}

	// update product query
	argId := 1
	args := make([]interface{}, 0)
	updateValues := make([]string, 0)

	if inp.Code.Valid {
		updateValues = append(updateValues, fmt.Sprintf("code=$%d", argId))
		args = append(args, inp.Code.String)
		argId++
	}

	if inp.InStock.Valid {
		updateValues = append(updateValues, fmt.Sprintf("in_stock=$%d", argId))
		args = append(args, inp.InStock.Bool)
		argId++
	}

	if inp.CategoryId != nil {
		updateValues = append(updateValues, fmt.Sprintf("category_id=$%d", argId))
		args = append(args, *inp.CategoryId)
		argId++
	}

	updateProductQuery := fmt.Sprintf("UPDATE %s SET %s WHERE id = $%d", productsTable, strings.Join(updateValues, ", "), argId)
	args = append(args, id)
	argId++

	_, err = r.db.Exec(updateProductQuery, args...)
	if err != nil {
		logrus.Errorf("[Update Product] update product error: %s", err.Error())
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func multiLanguageUpdateQuery(table string, input jewerly.MultiLanguageInput, productId int) (string, []interface{}) {
	var fieldName string
	switch table {
	case titlesTable:
		fieldName = "title_id"
	case descriptionsTable:
		fieldName = "description_id"
	case materialsTable:
		fieldName = "material_id"
	}

	query := fmt.Sprintf("UPDATE %s l SET l.english=$1, l.russian=$2, l.ukrainian=$3 JOIN %s p ON p.%s = l.id WHERE p.id = $4",
		table, productsTable, fieldName)
	args := []interface{}{input.English, input.Russian, input.Ukrainian, productId}

	return query, args
}

func multiCurrencyUpdateQuery(input jewerly.MultiCurrencyInput, productId int) (string, []interface{}) {
	query := fmt.Sprintf("UPDATE %s pr SET pr.usd=$1, pr.eur=$2, pr.uah=$3 JOIN %s p ON p.price_id = pr.id WHERE p.id = $4",
		pricesTable, productsTable)
	args := []interface{}{input.USD, input.EUR, input.UAH, productId}

	return query, args
}
