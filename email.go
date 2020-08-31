package jewerly

import (
	"time"
)

// todo: bind transaction statuses with custom statuses

type OrderInfoEmailInput struct {
	OrderId           int
	FirstName         string
	LastName          string
	Country           string
	Address           string
	PostalCode        string
	Email             string
	CardMask          string
	TotalCost         float32
	TransactionId     string
	TransactionStatus string
	OrderedAt         time.Time
	OrderedAtFormated string
	Products          []ProductInfo
}

type ProductInfo struct {
	Title    string
	Quantity int
	Price    float32
	ImageURL string
}

type PaymentInfoEmailInput struct {
	TransactionId string
	OrderId       int
	CardMask      string
	Status        string
}
