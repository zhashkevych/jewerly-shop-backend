package handler

import (
	"bytes"
	"errors"
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

func TestHandler_adminSignIn(t *testing.T) {
	// Init Test Data
	type mockBehavior func(r *mock_service.MockAdmin, login, password string)

	testCases := []struct {
		name                 string
		login                string
		password             string
		inputBody            string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			login:     "login",
			password:  "qwerty",
			inputBody: `{"login": "login", "password": "qwerty"}`,
			mockBehavior: func(r *mock_service.MockAdmin, login, password string) {
				r.EXPECT().SignIn(login, password).Return("token", nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"token":"token"}`,
		},
		{
			name:                 "Empty Password",
			login:                "login",
			inputBody:            `{"login": "login"}`,
			mockBehavior:         func(r *mock_service.MockAdmin, login, password string) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"invalid input body"}`,
		},
		{
			name:                 "No Input",
			inputBody:            `{}`,
			mockBehavior:         func(r *mock_service.MockAdmin, login, password string) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"invalid input body"}`,
		},
		{
			name:      "Service Failure",
			login:     "login",
			password:  "qwerty",
			inputBody: `{"login": "login", "password": "qwerty"}`,
			mockBehavior: func(r *mock_service.MockAdmin, login, password string) {
				r.EXPECT().SignIn(login, password).Return("", errors.New("invalid credentials"))
			},
			expectedStatusCode:   401,
			expectedResponseBody: `{"error":"invalid credentials"}`,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// Init Deps
			c := gomock.NewController(t)
			defer c.Finish()

			admin := mock_service.NewMockAdmin(c)
			test.mockBehavior(admin, test.login, test.password)

			services := &service.Services{Admin: admin}
			handler := Handler{services}

			// Init Endpoint
			r := gin.New()
			r.POST("/sign-in", handler.adminSignIn)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/sign-in",
				bytes.NewBufferString(test.inputBody))

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_createProduct(t *testing.T) {
	// Init Test Data
	type mockBehavior func(r *mock_service.MockProduct, product jewerly.CreateProductInput)

	testCases := []struct {
		name                 string
		fixturePath          string
		inputProduct         jewerly.CreateProductInput
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			fixturePath: "./fixtures/products/create.ok.json", // /pkg/handler/
			inputProduct: jewerly.CreateProductInput{
				Titles: jewerly.MultiLanguageInput{
					English: "Product 1",
					Ukrainian: "Продукт 1",
					Russian: "Продукт 1",
				},
				Descriptions: jewerly.MultiLanguageInput{
					English: "Description",
					Ukrainian: "Опис",
					Russian: "Описание",
				},
				Material: jewerly.MultiLanguageInput{
					English: "Material",
					Ukrainian: "Матеріал",
					Russian: "Материал",
				},
				Price: 199.99,
				Code: "ABC123",
				ImageIds: []int{1},
				CategoryId: jewerly.CategoryBracelets,
			},
			mockBehavior: func(r *mock_service.MockProduct, product jewerly.CreateProductInput) {
				r.EXPECT().Create(product).Return(nil)
			},
			expectedStatusCode:   200,
		},
		{
			name:      "Missing Titles",
			fixturePath: "./fixtures/products/create.no_titles.json", // /pkg/handler/
			mockBehavior: func(r *mock_service.MockProduct, product jewerly.CreateProductInput) {},
			expectedStatusCode:   400,
			expectedResponseBody:   `{"error":"invalid input body"}`,
		},
		{
			name:      "Missing Descriptions",
			fixturePath: "./fixtures/products/create.no_descriptions.json", // /pkg/handler/
			mockBehavior: func(r *mock_service.MockProduct, product jewerly.CreateProductInput) {},
			expectedStatusCode:   400,
			expectedResponseBody:   `{"error":"invalid input body"}`,
		},
		{
			name:      "Missing Materials",
			fixturePath: "./fixtures/products/create.no_materials.json", // /pkg/handler/
			mockBehavior: func(r *mock_service.MockProduct, product jewerly.CreateProductInput) {},
			expectedStatusCode:   400,
			expectedResponseBody:   `{"error":"invalid input body"}`,
		},
		{
			name:      "Missing Price",
			fixturePath: "./fixtures/products/create.no_price.json", // /pkg/handler/
			mockBehavior: func(r *mock_service.MockProduct, product jewerly.CreateProductInput) {},
			expectedStatusCode:   400,
			expectedResponseBody:   `{"error":"invalid input body"}`,
		},
		{
			name:      "Invalid Price",
			fixturePath: "./fixtures/products/create.no_price.json", // /pkg/handler/
			mockBehavior: func(r *mock_service.MockProduct, product jewerly.CreateProductInput) {},
			expectedStatusCode:   400,
			expectedResponseBody:   `{"error":"invalid input body"}`,
		},
		{
			name:      "Missing Code",
			fixturePath: "./fixtures/products/create.no_code.json", // /pkg/handler/
			mockBehavior: func(r *mock_service.MockProduct, product jewerly.CreateProductInput) {},
			expectedStatusCode:   400,
			expectedResponseBody:   `{"error":"invalid input body"}`,
		},
		{
			name:      "Empty Ids",
			fixturePath: "./fixtures/products/create.empty_ids.json", // /pkg/handler/
			mockBehavior: func(r *mock_service.MockProduct, product jewerly.CreateProductInput) {},
			expectedStatusCode:   400,
			expectedResponseBody:   `{"error":"invalid input body"}`,
		},
		{
			name:      "Missing Ids",
			fixturePath: "./fixtures/products/create.no_ids.json", // /pkg/handler/
			mockBehavior: func(r *mock_service.MockProduct, product jewerly.CreateProductInput) {},
			expectedStatusCode:   400,
			expectedResponseBody:   `{"error":"invalid input body"}`,
		},
		{
			name:      "Missing Category",
			fixturePath: "./fixtures/products/create.no_category.json", // /pkg/handler/
			mockBehavior: func(r *mock_service.MockProduct, product jewerly.CreateProductInput) {},
			expectedStatusCode:   400,
			expectedResponseBody:   `{"error":"invalid input body"}`,
		},
		{
			name:      "Invalid Category",
			fixturePath: "./fixtures/products/create.invalid_category.json", // /pkg/handler/
			mockBehavior: func(r *mock_service.MockProduct, product jewerly.CreateProductInput) {},
			expectedStatusCode:   400,
			expectedResponseBody:   `{"error":"invalid category"}`,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// Init Deps
			c := gomock.NewController(t)
			defer c.Finish()

			product := mock_service.NewMockProduct(c)
			test.mockBehavior(product, test.inputProduct)

			services := &service.Services{Product: product}
			handler := Handler{services}

			// Init Endpoint
			r := gin.New()
			r.POST("/product", handler.createProduct)

			// Input Body
			fixtureData, err := ioutil.ReadFile(test.fixturePath)
			assert.NoError(t, err)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/product",
				bytes.NewBuffer(fixtureData))

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedResponseBody, w.Body.String())
		})
	}
}
