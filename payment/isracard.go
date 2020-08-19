package payment

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

const (
	statusFail           = 1
	generateSaleEndpoint = "generate-sale"
	defaultLanguage      = "en"
)

type IsracardProvider struct {
	endpoint string
	apiKey   string

	callbackURL string
	returnURL   string

	client http.Client
}

func NewIsracardProvider(endpoint, apiKey, returnURL, callbackURL string) *IsracardProvider {
	return &IsracardProvider{
		endpoint:    endpoint,
		apiKey:      apiKey,
		returnURL:   returnURL,
		callbackURL: callbackURL,
		client: http.Client{
			Timeout: time.Second * 5,
		}}
}

type generateSaleInput struct {
	SellerPaymeID string `json:"seller_payme_id"`
	Price         int    `json:"sale_price"`
	Currency      string `json:"currency"`
	ProductName   string `json:"product_name"`
	TransactionID string `json:"transaction_id"`
	CallbackURL   string `json:"sale_callback_url"`
	ReturnURL     string `json:"sale_return_url"`
	Language      string `json:"language"`
}

type generateSaleResponse struct {
	StatusCode int    `json:"status_code"`
	SaleURL    string `json:"sale_url"`
	SaleID     string `json:"payme_sale_id"`
	SaleCode   int    `json:"payme_sale_code"`
}

func (p *IsracardProvider) GenerateSale(inp GenerateSaleInput) (string, error) {
	input := &generateSaleInput{
		SellerPaymeID: p.apiKey,
		Price:         inp.Price,
		ProductName:   inp.ProductName,
		Currency:      inp.Currency,
		TransactionID: inp.TransactionID,
		Language:      defaultLanguage,
		CallbackURL:   p.callbackURL,
		ReturnURL:     p.returnURL,
	}
	out := new(generateSaleResponse)

	err := p.do(http.MethodPost, generateSaleEndpoint, input, out)
	if err != nil {
		return "", err
	}

	logrus.Debugf("resp: %+v\n", out)

	if out.StatusCode == statusFail {
		return "", errors.New("generate sail fail")
	}

	return out.SaleURL, nil
}

// http client
func (p *IsracardProvider) do(method, endpoint string, input, out interface{}) error {
	body, err := json.Marshal(input)
	if err != nil {
		logrus.Errorf("Error occurred while marshaling body: %s\n", err.Error())
		return err
	}

	req, err := http.NewRequest(method, p.endpoint+endpoint, bytes.NewBuffer(body))
	if err != nil {
		logrus.Errorf("Error occurred while forming request: %s\n", err.Error())
		return err
	}

	req.Header.Set("Content-type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		logrus.Errorf("Error occurred while sending request: %s\n", err.Error())
		return err
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return json.NewDecoder(resp.Body).Decode(&out)
	default:
		logrus.Errorf("Error occurred while sending request, status code: %d\n", resp.StatusCode)
		return errors.New("request unsuccessful")
	}
}
