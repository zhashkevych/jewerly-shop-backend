package payment

type GetPaymentFormInput struct {
	Type       int
	Lang       string
	Currency   string
	Amount     float32
	Client     Client
	SuccessURL string
	FailureURL string
	NotifyURL  string
}

type Client struct {
	Name    string
	Emails  []string
	Address string
	City    string
	Zip     string
	Country string
	Mobile  string
}

type Provider interface {
	GetPaymentForm()
	GetAvailablePaymentTypes()
}
