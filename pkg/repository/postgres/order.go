package postgres

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	jewerly "github.com/zhashkevych/jewelry-shop-backend"
	"strings"
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

	orderId, err := r.createOrder(tx, input)
	if err != nil {
		return 0, err
	}

	err = r.createOrderItems(tx, orderId, input.Items)
	if err != nil {
		return 0, err
	}

	err = r.createOrderTransactionRecords(tx, orderId, input.TransactionID)
	if err != nil {
		return 0, err
	}

	return orderId, tx.Commit()
}

func (r *OrderRepository) GetOrderProducts(items []jewerly.OrderItem) ([]jewerly.ProductResponse, error) {
	products, err := r.getOrderProducts(items)
	if err != nil {
		return nil, err
	}

	err = r.getProductsImages(products)

	return products, err
}

func (r *OrderRepository) getOrderProducts(items []jewerly.OrderItem) ([]jewerly.ProductResponse, error) {
	var products []jewerly.ProductResponse

	ids := make([]string, len(items))
	values := make([]interface{}, len(items))

	for i := range items {
		values[i] = items[i].ProductId
		ids[i] = fmt.Sprintf("$%d", i+1)
	}

	err := r.db.Select(&products, fmt.Sprintf("SELECT p.id, p.price FROM %s p INNER JOIN %s t ON t.id = p.title_id WHERE p.id IN (%s) and p.in_stock=true",
		productsTable, titlesTable, strings.Join(ids, ",")), values...)

	return products, err
}

func (r *OrderRepository) getProductsImages(products []jewerly.ProductResponse) error {
	for i := range products {
		err := r.db.Select(&products[i].Images,
			fmt.Sprintf("SELECT i.id, i.url, i.alt_text FROM %s i INNER JOIN %s pi on pi.image_id = i.id WHERE pi.product_id=$1", imagesTable, productImagesTable),
			products[i].Id)
		if err != nil {
			return err
		}
	}

	return nil
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

func (r *OrderRepository) GetById(id int) (jewerly.Order, error) {
	var order jewerly.Order

	selectOrdersQuery := fmt.Sprintf(`SELECT id, ordered_at, first_name, last_name, additional_name, country,
										address, email, postal_code, total_cost FROM %s WHERE id=$1`, ordersTable)
	err := r.db.Get(&order, selectOrdersQuery, id)
	if err != nil {
		logrus.Errorf("failed to get orders: %s", err.Error())
		return order, err
	}

	selectOrderItemsQuery := fmt.Sprintf("SELECT product_id, quantity FROM %s WHERE order_id = $1", orderItemsTable)
	err = r.db.Select(&order.Items, selectOrderItemsQuery, id)
	if err != nil {
		logrus.Errorf("failed to get order items for order id %d, error: %s", id, err.Error())
		return order, err
	}

	selectTransactionsQuery := fmt.Sprintf(`SELECT th.uuid, th.created_at, th.status, th.card_mask FROM %s th 
											INNER JOIN %s t on t.uuid = th.uuid WHERE t.order_id = $1`, transactionsHistoryTable, transactionsTable)
	err = r.db.Select(&order.Transactions, selectTransactionsQuery, id)
	if err != nil {
		logrus.Errorf("failed to get transactions for order id %d, error: %s", id, err.Error())
		return order, err
	}

	return order, nil
}

func (r *OrderRepository) createOrder(tx *sql.Tx, input jewerly.CreateOrderInput) (int, error) {
	var orderId int
	createOrderQuery := fmt.Sprintf(`INSERT INTO %s (first_name, last_name, additional_name, country, address, postal_code, email, total_cost)
									VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`, ordersTable)
	row := tx.QueryRow(createOrderQuery, input.FirstName, input.LastName, input.AdditionalName, input.Country, input.Address,
		input.PostalCode, input.Email, input.TotalCost)
	err := row.Scan(&orderId)
	if err != nil {
		logrus.Errorf("failed to create new order: %s", err.Error())
		tx.Rollback()
		return 0, err
	}

	return orderId, nil
}

func (r *OrderRepository) createOrderItems(tx *sql.Tx, orderId int, orderItems []jewerly.OrderItem) error {
	items := []string{}
	values := []interface{}{}
	values = append(values, orderId)
	argId := 2

	for _, item := range orderItems {
		values = append(values, item.ProductId, item.Quantity)
		items = append(items, fmt.Sprintf("($1, $%d, $%d)", argId, argId+1))

		argId += 2
	}

	createOrderItemsQuery := fmt.Sprintf("INSERT INTO %s (order_id, product_id, quantity) VALUES %s", orderItemsTable, strings.Join(items, ","))

	_, err := tx.Exec(createOrderItemsQuery, values...)
	if err != nil {
		logrus.Errorf("failed to create order items: %s", err.Error())
		tx.Rollback()
		return err
	}

	return nil
}

func (r *OrderRepository) createOrderTransactionRecords(tx *sql.Tx, orderId int, transactionId string) error {
	_, err := tx.Exec(fmt.Sprintf("INSERT INTO %s (order_id, uuid) VALUES ($1, $2)", transactionsTable),
		orderId, transactionId)
	if err != nil {
		logrus.Errorf("failed to create transaction: %s", err.Error())
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(fmt.Sprintf("INSERT INTO %s (uuid) VALUES ($1)", transactionsHistoryTable), transactionId)
	if err != nil {
		logrus.Errorf("failed to insert transaction history record: %s", err.Error())
		tx.Rollback()
		return err
	}

	return nil
}
