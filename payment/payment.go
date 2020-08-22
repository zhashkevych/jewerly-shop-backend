package payment

type GenerateSaleInput struct {
	Price int
	Currency string
	ProductName string
	TransactionID string
}

type Provider interface {
	GenerateSale(inp GenerateSaleInput) (string, error)
}
