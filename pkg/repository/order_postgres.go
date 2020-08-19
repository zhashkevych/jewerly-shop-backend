package repository

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/jmoiron/sqlx"
	jewerly "github.com/zhashkevych/jewelry-shop-backend"
)

type OrderRepository struct {
	db *sqlx.DB
}

func NewOrderRepository(db *sqlx.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(input jewerly.CreateOrderInput) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	// create order
	var orderId int
	createOrderQuery := fmt.Sprintf(`INSERT INTO %s (first_name, last_name, additional_name, country, address, postal_code, email, total_cost)
									VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`, ordersTable)
	row := tx.QueryRow(createOrderQuery, input.FirstName, input.LastName, input.AdditionalName, input.Country, input.Address,
		input.PostalCode, input.Email, input.TotalCost)

	err = row.Scan(&orderId)
	if err != nil {
		logrus.Errorf("failed to create new order: %s", err.Error())
		tx.Rollback()
		return 0, err
	}

	// create order items
	items := ""
	values := []interface{}{}
	values = append(values, orderId)
	argId := 2

	for i, item := range input.Items {
		values = append(values, item.ProductId, item.Quantity)

		if i == len(input.Items)-1 {
			items += fmt.Sprintf("($1, $%d, $%d)", argId, argId+1)
			break
		}

		items += fmt.Sprintf("($1, $%d, $%d), ", argId, argId+1)

		argId += 2
	}

	createOrderItemsQuery := fmt.Sprintf("INSERT INTO %s (order_id, product_id, quantity) VALUES %s", orderItemsTable, items)
	logrus.Debug(createOrderItemsQuery)

	_, err = tx.Exec(createOrderItemsQuery, values...)
	if err != nil {
		logrus.Errorf("failed to create order items: %s", err.Error())
		tx.Rollback()
		return 0, err
	}

	//create transaction
	_, err = tx.Exec(fmt.Sprintf("INSERT INTO %s (order_id, uuid) VALUES ($1, $2)", transactionsTable),
		orderId, input.TransactionID)
	if err != nil {
		logrus.Errorf("failed to create transaction: %s", err.Error())
		tx.Rollback()
		return 0, err
	}

	return orderId, tx.Commit()
}

func (r *OrderRepository) GetOrderProducts(items []jewerly.OrderItem) ([]jewerly.ProductResponse, error) {
	var products []jewerly.ProductResponse

	ids := ""
	values := make([]interface{}, len(items))

	for i := range items {
		values[i] = items[i].ProductId

		if i == len(items)-1 {
			ids += fmt.Sprintf("$%d", i+1)
			break
		}

		ids += fmt.Sprintf("$%d, ", i+1)
	}

	err := r.db.Select(&products, fmt.Sprintf("SELECT * FROM %s WHERE id IN (%s)", productsTable, ids), values...)
	return products, err
}

func (r *OrderRepository) UpdateTransaction(transactionId, cardMask, status string) error {
	_, err := r.db.Exec(fmt.Sprintf("UPDATE %s SET card_mask=$1, status=$2 where uuid=$3", transactionsTable), cardMask, status, transactionId)
	return err
}
