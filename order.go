package jewerly

// todo: input validation

type CreateOrderInput struct {
	Items          []OrderItem `json:"items" binding:"required"`
	FirstName      string      `json:"first_name"  binding:"required"`
	LastName       string      `json:"last_name"  binding:"required"`
	AdditionalName string      `json:"additional_name"`
	Email          string      `json:"email"  binding:"required"`
	Phone          string      `json:"phone"`
	Country        string      `json:"country"  binding:"required"`
	Address        string      `json:"address"  binding:"required"`
	PostalCode     string      `json:"postal_code"  binding:"required"`
	TransactionID  string
	TotalCost      float32
}

type OrderItem struct {
	ProductId int `json:"product_id"  binding:"required"`
	Quantity  int `json:"quantity"  binding:"required"`
}

type TransactionCallbackInput struct {
	StatusCode         int    `form:"status_code"`
	StatusErrorCode    int    `form:"status_error_code"`
	StatusErrorDetails string `form:"status_error_details"`
	NotifyType         string `form:"notify_type"`
	SaleCreated        string `form:"sale_created"`
	TransactionID      string `form:"transaction_id"`
	SaleStatus         string `form:"sale_status"`
	BuyerCardMask      string `form:"buyer_card_mask"`
	BuyerCardExp       string `form:"buyer_card_exp"`
	BuyerName          string `form:"buyer_name"`
	BuyerEmail         string `form:"buyer_email"`
	BuyerPhone         string `form:"buyer_phone"`
	SalePaidDate       string `form:"sale_paid_date"`
	SaleReleaseDate    string `form:"sale_release_date"`
	SaleInvoiceURL     string `form:"sale_invoice_url"`
}
