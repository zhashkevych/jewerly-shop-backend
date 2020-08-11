package payment

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type httpClient struct {
	endpoint  string
	apiKey    string
	secretKey string

	client      http.Client
	accessToken string
}

func newHttpClient(endpoint, apiKey, secretKey string) (*httpClient, error) {
	client := &httpClient{
		endpoint:  endpoint,
		apiKey:    apiKey,
		secretKey: secretKey,
		client: http.Client{
			Timeout: time.Second * 5,
		}}

	err := client.authorize()
	return client, err
}

func (c *httpClient) request(method, endpoint string, input, out interface{}) error {
	err := c.do(method, endpoint, input, out)

	switch err {
	case errUnauthorized:
		if err := c.authorize(); err != nil {
			logrus.Errorf("authorization failure: %s", err.Error())
			return err
		}

		return c.do(method, endpoint, input, out)
	default:
		return nil
	}
}

func (c *httpClient) do(method, endpoint string, input, out interface{}) error {
	body, err := json.Marshal(input)
	if err != nil {
		logrus.Errorf("Error occurred while marshaling body: %s\n", err.Error())
		return err
	}

	req, err := http.NewRequest(method, c.endpoint+endpoint, bytes.NewBuffer(body))
	if err != nil {
		logrus.Errorf("Error occurred while forming request: %s\n", err.Error())
		return err
	}

	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.accessToken))

	resp, err := c.client.Do(req)
	if err != nil {
		logrus.Errorf("Error occurred while sending request: %s\n", err.Error())
		return err
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return json.NewDecoder(resp.Body).Decode(&out)
	case http.StatusUnauthorized:
		return errUnauthorized
	default:
		logrus.Errorf("Error occurred while sending request, status code: %d\n", resp.StatusCode)
		return errors.New("request unsuccessful")
	}
}

type authorizationInput struct {
	Id     string `json:"id"`
	Secret string `json:"secret"`
}

type authorizationResponse struct {
	Token     string `json:"token"`
	ExpiresIn int64  `json:"expires"`
}

func (c *httpClient) authorize() error {
	body, err := json.Marshal(authorizationInput{
		Id:     c.apiKey,
		Secret: c.secretKey,
	})
	if err != nil {
		logrus.Errorf("Error occurred while marshaling body: %s\n", err.Error())
		return err
	}

	req, err := http.NewRequest(http.MethodPost, c.endpoint+authorizationEndpoint, bytes.NewBuffer(body))
	if err != nil {
		logrus.Errorf("Error occurred while forming request: %s\n", err.Error())
		return err
	}

	req.Header.Set("Content-type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		logrus.Errorf("Error occurred while sending request: %s\n", err.Error())
		return err
	}

	if resp.StatusCode != http.StatusOK {
		logrus.Errorf("Error occurred while sending request, status code: %d\n", resp.StatusCode)
		return errors.New("request unsuccessful")
	}

	var out authorizationResponse
	err = json.NewDecoder(resp.Body).Decode(&out)
	if err != nil {
		logrus.Errorf("Error occurred while decoding response: %s\n", err.Error())
		return err
	}

	c.accessToken = out.Token

	return nil
}
