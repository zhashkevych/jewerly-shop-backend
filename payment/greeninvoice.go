package payment

import (
	"errors"
)

const (
	authorizationEndpoint  = "account/token"
	getPaymentFormEndpoint = "payments/form"
)

var (
	errUnauthorized = errors.New("invalid access token")
)

type GreenInvoiceProvider struct {
	client *httpClient
}

func NewGreenInvoiceProvider(endpoint string, apiKey string, secretKey string) (*GreenInvoiceProvider, error) {
	client, err := newHttpClient(endpoint, apiKey, secretKey)
	if err != nil {
		return nil, err
	}

	return &GreenInvoiceProvider{
		client: client,
	}, nil
}

func (p *GreenInvoiceProvider) GetPaymentForm() {

}

func (p *GreenInvoiceProvider) GetAvailablePaymentTypes() {

}
