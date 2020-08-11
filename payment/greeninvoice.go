package payment

import (
	"bytes"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

const (
	authorizationEndpoint = "account/token"
)

type GreenInvoiceProvider struct {
	endpoint  string
	apiKey    string
	secretKey string

	client      http.Client
	accessToken string
}

func NewGreenInvoiceProvider(endpoint string, apiKey string, secretKey string) (*GreenInvoiceProvider, error) {
	provider := &GreenInvoiceProvider{
		endpoint:  endpoint,
		apiKey:    apiKey,
		secretKey: secretKey,
		client: http.Client{
			Timeout: time.Second * 5,
		},
	}

	err := provider.authorization()

	return provider, err
}

func (p *GreenInvoiceProvider) GetPaymentForm() {

}

func (p *GreenInvoiceProvider) GetAvailablePaymentTypes() {

}

type authorizationInput struct {
	Id     string `json:"id"`
	Secret string `json:"secret"`
}

type authorizationResponse struct {
	Token     string `json:"token"`
	ExpiresIn int64  `json:"expires"`
}

func (p *GreenInvoiceProvider) authorization() error {
	body, err := json.Marshal(authorizationInput{
		Id:     p.apiKey,
		Secret: p.secretKey,
	})
	if err != nil {
		logrus.Errorf("Error occurred while marshaling body: %s\n", err.Error())
		return err
	}

	req, err := http.NewRequest(http.MethodPost, p.endpoint+authorizationEndpoint, bytes.NewBuffer(body))
	if err != nil {
		logrus.Errorf("Error occurred while forming request: %s\n", err.Error())
		return err
	}

	req.Header.Set("Content-type", "application/json")

	logrus.Debugf("req: %+v\n", req)

	resp, err := p.client.Do(req)
	if err != nil {
		logrus.Errorf("Error occurred while sending request: %s\n", err.Error())
		return err
	}

	logrus.Debugf("auth status: %d\n", resp.StatusCode)
	logrus.Debugf("auth resp: %s\n", resp.Body)
	logrus.Debugf("auth header: %s\n", resp.Header.Get("X-Authorization-Bearer"))


	var out authorizationResponse
	err = json.NewDecoder(resp.Body).Decode(&out)
	if err != nil {
		logrus.Errorf("Error occurred while decoding response: %s\n", err.Error())
		return err
	}

	logrus.Debugf("respo: %+v\n", out)

	p.accessToken = out.Token

	logrus.Debugf("AccessToken: %s\n", p.accessToken)

	return nil
}
