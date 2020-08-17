package payment

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

// todo: use default consts for language
// todo: implement callback endpoints and set them in request

const (
	generateSaleEndpoint = "generate-sale"

	StatusSuccess = 0
	StatusFail    = 1
)

type IsracardProvider struct {
	endpoint string
	apiKey   string

	client http.Client
}

func NewIsracardProvider(endpoint string, apiKey string) *IsracardProvider {
	return &IsracardProvider{
		endpoint: endpoint,
		apiKey:   apiKey,
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
		Language:      "en",
		CallbackURL:   "https://www.example.com/payment/callback",
		ReturnURL:     "https://www.example.com/payment/success",
	}

	out := new(generateSaleResponse)

	logrus.Debugf("Sending: %+v\n", input)

	err := p.do(http.MethodPost, generateSaleEndpoint, input, out)
	if err != nil {
		return "", err
	}

	if out.StatusCode == StatusFail {
		return "", errors.New("generate sail fail")
	}

	logrus.Debugf("resp: %+v\n", out)

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
