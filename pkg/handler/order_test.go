package handler

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	jewerly "github.com/zhashkevych/jewelry-shop-backend"
	"github.com/zhashkevych/jewelry-shop-backend/pkg/service"
	mock_service "github.com/zhashkevych/jewelry-shop-backend/pkg/service/mocks"
	"io/ioutil"
	"net/http/httptest"
	"testing"
)

func TestHandler_placeOrder(t *testing.T) {
	type mockBehaviour func(s *mock_service.MockOrder, input jewerly.CreateOrderInput)

	testTable := []struct {
		name                 string
		fixturePath          string
		orderInput           jewerly.CreateOrderInput
		mockBehavior         mockBehaviour
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "Ok",
			fixturePath: "./fixtures/orders/ok.json",
			orderInput: jewerly.CreateOrderInput{
				Items: []jewerly.OrderItem{
					{1, 3},
					{2, 3},
				},
				FirstName:      "Vasya",
				LastName:       "Pupkin",
				AdditionalName: "Aleksandrovich",
				Email:          "vasya@pupkin.com",
				Phone:          "+380950515344",
				Country:        "UA",
				Address:        "st. Khreshatyk, Kiev",
				PostalCode:     "12303",
			},
			mockBehavior: func(s *mock_service.MockOrder, input jewerly.CreateOrderInput) {
				s.EXPECT().Create(input).Return("http://payment.link", nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"url":"http://payment.link"}`,
		},
		{
			name:                 "Address Empty",
			fixturePath:          "./fixtures/orders/address.empty.json",
			mockBehavior:         func(s *mock_service.MockOrder, input jewerly.CreateOrderInput) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"invalid input body"}`,
		},
		{
			name:                 "Country Empty",
			fixturePath:          "./fixtures/orders/country.empty.json",
			mockBehavior:         func(s *mock_service.MockOrder, input jewerly.CreateOrderInput) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"invalid input body"}`,
		},
		{
			name:                 "Email Empty",
			fixturePath:          "./fixtures/orders/email.empty.json",
			mockBehavior:         func(s *mock_service.MockOrder, input jewerly.CreateOrderInput) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"invalid input body"}`,
		},
		{
			name:                 "Email Invalid",
			fixturePath:          "./fixtures/orders/email.invalid.json",
			mockBehavior:         func(s *mock_service.MockOrder, input jewerly.CreateOrderInput) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"invalid input body"}`,
		},
		{
			name:                 "Items Empty",
			fixturePath:          "./fixtures/orders/items.empty.json",
			mockBehavior:         func(s *mock_service.MockOrder, input jewerly.CreateOrderInput) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"order should have at least 1 item"}`,
		},
		{
			name:                 "Items Invalid [product_id=0]",
			fixturePath:          "./fixtures/orders/items.invalid.product_id_0.json",
			mockBehavior:         func(s *mock_service.MockOrder, input jewerly.CreateOrderInput) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"order item is invalid"}`,
		},
		{
			name:                 "Items Invalid [product_id=negative]",
			fixturePath:          "./fixtures/orders/items.invalid.product_id_negative.json",
			mockBehavior:         func(s *mock_service.MockOrder, input jewerly.CreateOrderInput) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"order item is invalid"}`,
		},
		{
			name:                 "Items Invalid [quantity=0]",
			fixturePath:          "./fixtures/orders/items.invalid.product_id_negative.json",
			mockBehavior:         func(s *mock_service.MockOrder, input jewerly.CreateOrderInput) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"order item is invalid"}`,
		},
		{
			name:                 "Items Invalid [quantity=negative]",
			fixturePath:          "./fixtures/orders/items.invalid.product_id_negative.json",
			mockBehavior:         func(s *mock_service.MockOrder, input jewerly.CreateOrderInput) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"order item is invalid"}`,
		},
		{
			name:                 "Name Empty",
			fixturePath:          "./fixtures/orders/name.empty.json",
			mockBehavior:         func(s *mock_service.MockOrder, input jewerly.CreateOrderInput) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"invalid input body"}`,
		},
		{
			name:                 "Postal Code Empty",
			fixturePath:          "./fixtures/orders/postal_code.empty.json",
			mockBehavior:         func(s *mock_service.MockOrder, input jewerly.CreateOrderInput) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"invalid input body"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			// Init Deps
			c := gomock.NewController(t)
			defer c.Finish()

			order := mock_service.NewMockOrder(c)
			testCase.mockBehavior(order, testCase.orderInput)

			services := &service.Services{Order: order}
			handler := Handler{services}

			// Init Endpoint
			r := gin.New()
			r.POST("/order", handler.placeOrder)

			// Input Body
			fixtureData, err := ioutil.ReadFile(testCase.fixturePath)
			assert.NoError(t, err)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/order", bytes.NewBuffer(fixtureData))
			assert.NoError(t, err)

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, w.Body.String())
		})
	}
}
