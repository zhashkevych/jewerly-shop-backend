package jewerly

import (
	"time"
)

type OrderInfoEmailInput struct {
	OrderId           int
	FirstName         string
	LastName          string
	Country           string
	Address           string
	PostalCode        string
	Email             string
	CardMask          string
	TotalCost         string
	TransactionId     string
	TransactionStatus string
	OrderedAt         time.Time
	OrderedAtFormated string
	Products          []ProductInfo
}

type ProductInfo struct {
	Id       int
	Title    string
	Quantity int
	Price    float32
	ImageURL string
}

type PaymentInfoEmailInput struct {
	TransactionId string
	OrderId       int
	CardMask      string
	CardBrand     string
	Price         float32
	Currency      string
	BuyerName     string
	BuyerEmail    string
	Status        string
}
