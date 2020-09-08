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

	_, err = tx.Exec(fmt.Sprintf("INSERT INTO %s (uuid) VALUES ($1)", transactionsHistoryTable), input.TransactionID)
	if err != nil {
		logrus.Errorf("failed to insert transaction history record: %s", err.Error())
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

	err := r.db.Select(&products, fmt.Sprintf(`SELECT p.id, pr.usd as price, t.english as title 
							FROM %s p INNER JOIN %s t ON t.id = p.title_id INNER JOIN %s pr on pr.id = p.price_id
							WHERE p.id IN (%s) and p.in_stock=true`, productsTable, titlesTable, pricesTable, ids), values...)
	if err != nil {
		return products, err
	}

	for i := range products {
		err := r.db.Select(&products[i].Images,
			fmt.Sprintf("SELECT i.id, i.url, i.alt_text FROM %s i INNER JOIN %s pi on pi.image_id = i.id WHERE pi.product_id=$1", imagesTable, productImagesTable),
			products[i].Id)
		if err != nil {
			return products, err
		}
	}

	return products, nil
}

func (r *OrderRepository) CreateTransaction(transactionId, cardMask, status string) error {
	_, err := r.db.Exec(fmt.Sprintf("INSERT INTO %s (uuid, card_mask, status) VALUES ($1, $2, $3)", transactionsHistoryTable),
		transactionId, cardMask, status)
	return err
}

func (r *OrderRepository) GetOrderId(transactionId string) (int, error) {
	var id int
	err := r.db.Get(&id, fmt.Sprintf("SELECT order_id FROM %s WHERE uuid=$1", transactionsTable), transactionId)
	return id, err
}

func (r *OrderRepository) GetAll(input jewerly.GetAllOrdersFilters) (jewerly.OrderList, error) {
	var orders jewerly.OrderList

	selectOrdersQuery := fmt.Sprintf(`SELECT id, ordered_at, first_name, last_name, additional_name, country,
										address, email, postal_code, total_cost FROM %s OFFSET $1 LIMIT $2`, ordersTable)
	err := r.db.Select(&orders.Data, selectOrdersQuery, input.Offset, input.Limit)
	if err != nil {
		logrus.Errorf("failed to get orders: %s", err.Error())
		return orders, err
	}

	selectCountQuery := fmt.Sprintf(`SELECT count(*) FROM %s`, ordersTable)
	err = r.db.Get(&orders.Total, selectCountQuery)
	if err != nil {
		logrus.Errorf("failed to get orders count: %s", err.Error())
		return orders, err
	}

	selectOrderItemsQuery := fmt.Sprintf("SELECT product_id, quantity FROM %s WHERE order_id = $1", orderItemsTable)

	for i := range orders.Data {
		err = r.db.Select(&orders.Data[i].Items, selectOrderItemsQuery, orders.Data[i].Id)
		if err != nil {
			logrus.Errorf("failed to get order items for order id %d, error: %s", orders.Data[i].Id, err.Error())
			return orders, err
		}
	}

	selectTransactionsQuery := fmt.Sprintf(`SELECT th.uuid, th.created_at, th.status, th.card_mask FROM %s th 
											INNER JOIN %s t on t.uuid = th.uuid WHERE t.order_id = $1`, transactionsHistoryTable, transactionsTable)
	for i := range orders.Data {
		err = r.db.Select(&orders.Data[i].Transactions, selectTransactionsQuery, orders.Data[i].Id)
		if err != nil {
			logrus.Errorf("failed to get transactions for order id %d, error: %s", orders.Data[i].Id, err.Error())
			return orders, err
		}
	}

	return orders, nil
}

// todo implement
func (r *OrderRepository) GetById(id int) (jewerly.Order, error) {
	return jewerly.Order{}, nil
}
