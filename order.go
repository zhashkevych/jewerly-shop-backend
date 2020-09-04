package jewerly

import (
	"gopkg.in/guregu/null.v3"
	"time"
)

const (
	TransactionStatusCreated    = "Created"
	TransactionStatusPaid       = "Paid"
	TransactionStatusFailed     = "Payment Failed"
	TransactionStatusAuthorized = "Payment Authorized"
	TransactionStatusRefunded   = "Payment Refunded"
	TransactionStatusChargeback = "Payment Chargeback"
	TransactionStatusReverted   = "Payment Reverted"
)

type CreateOrderInput struct {
	Items          []OrderItem `json:"items" binding:"required"`
	FirstName      string      `json:"first_name"  binding:"required"`
	LastName       string      `json:"last_name"  binding:"required"`
	AdditionalName string      `json:"additional_name"`
	Email          string      `json:"email"  binding:"email,required"`
	Phone          string      `json:"phone"`
	Country        string      `json:"country"  binding:"required"`
	Address        string      `json:"address"  binding:"required"`
	PostalCode     string      `json:"postal_code"  binding:"required"`
	TransactionID  string
	TotalCost      float32
}

type OrderItem struct {
	ProductId int `json:"product_id" db:"product_id"  binding:"required"`
	Quantity  int `json:"quantity" db:"quantity" binding:"required"`
}

type TransactionCallbackInput struct {
	StatusCode         int    `form:"status_code"`
	StatusErrorCode    int    `form:"status_error_code"`
	StatusErrorDetails string `form:"status_error_details"`
	NotifyType         string `form:"notify_type"`
	SaleCreated        string `form:"sale_created"`
	TransactionID      string `form:"transaction_id"`
	SaleStatus         string `form:"sale_status"`
	CardBrand          string `form:"payme_transaction_card_brand"`
	Currency           string `form:"currency"`
	Price              int    `form:"price"`
	BuyerCardMask      string `form:"buyer_card_mask"`
	BuyerCardExp       string `form:"buyer_card_exp"`
	BuyerName          string `form:"buyer_name"`
	BuyerEmail         string `form:"buyer_email"`
	BuyerPhone         string `form:"buyer_phone"`
	SalePaidDate       string `form:"sale_paid_date"`
	SaleReleaseDate    string `form:"sale_release_date"`
	SaleInvoiceURL     string `form:"sale_invoice_url"`
}

type Order struct {
	Id             int           `json:"id" db:"id"`
	OrderedAt      time.Time     `json:"ordered_at" db:"ordered_at"`
	FirstName      string        `json:"first_name" db:"first_name"`
	LastName       string        `json:"last_name" db:"last_name"`
	AdditionalName string        `json:"additional_name" db:"additional_name"`
	Country        string        `json:"country" db:"country"`
	Address        string        `json:"address" db:"address"`
	Email          string        `json:"email" db:"email"`
	PostalCode     string        `json:"postal_code" db:"postal_code"`
	TotalCost      float32       `json:"total_cost" db:"total_cost"`
	Items          []OrderItem   `json:"items"`
	Transactions   []Transaction `json:"transactions"`
}

type Transaction struct {
	TransactionId string      `json:"transaction_id" db:"uuid"`
	CardMask      null.String `json:"card_mask" db:"card_mask"`
	Status        string      `json:"status" db:"status"`
	CreatedAt     time.Time   `json:"created_at" db:"created_at"`
}

type OrderList struct {
	Data  []Order `json:"data"`
	Total int     `json:"total"`
}

type GetAllOrdersFilters struct {
	Offset int
	Limit  int
}
