package handler

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	jewerly "github.com/zhashkevych/jewelry-shop-backend"
	"github.com/zhashkevych/jewelry-shop-backend/pkg/service"
	mock_service "github.com/zhashkevych/jewelry-shop-backend/pkg/service/mocks"
	"gopkg.in/guregu/null.v3"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
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
			name:        "Ok",
			fixturePath: "./fixtures/products/create.ok.json",
			inputProduct: jewerly.CreateProductInput{
				Titles: jewerly.MultiLanguageInput{
					English:   "Product 1",
					Ukrainian: "Продукт 1",
					Russian:   "Продукт 1",
				},
				Descriptions: jewerly.MultiLanguageInput{
					English:   "Description",
					Ukrainian: "Опис",
					Russian:   "Описание",
				},
				Material: jewerly.MultiLanguageInput{
					English:   "Material",
					Ukrainian: "Матеріал",
					Russian:   "Материал",
				},
				Price:      199.99,
				Code:       "ABC123",
				ImageIds:   []int{1},
				CategoryId: jewerly.CategoryBracelets,
			},
			mockBehavior: func(r *mock_service.MockProduct, product jewerly.CreateProductInput) {
				r.EXPECT().Create(product).Return(nil)
			},
			expectedStatusCode: 200,
		},
		{
			name:                 "Missing Titles",
			fixturePath:          "./fixtures/products/create.no_titles.json", // /pkg/handler/
			mockBehavior:         func(r *mock_service.MockProduct, product jewerly.CreateProductInput) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"invalid input body"}`,
		},
		{
			name:                 "Missing Descriptions",
			fixturePath:          "./fixtures/products/create.no_descriptions.json", // /pkg/handler/
			mockBehavior:         func(r *mock_service.MockProduct, product jewerly.CreateProductInput) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"invalid input body"}`,
		},
		{
			name:                 "Missing Materials",
			fixturePath:          "./fixtures/products/create.no_materials.json", // /pkg/handler/
			mockBehavior:         func(r *mock_service.MockProduct, product jewerly.CreateProductInput) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"invalid input body"}`,
		},
		{
			name:                 "Missing Price",
			fixturePath:          "./fixtures/products/create.no_price.json", // /pkg/handler/
			mockBehavior:         func(r *mock_service.MockProduct, product jewerly.CreateProductInput) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"invalid input body"}`,
		},
		{
			name:                 "Invalid Price",
			fixturePath:          "./fixtures/products/create.no_price.json", // /pkg/handler/
			mockBehavior:         func(r *mock_service.MockProduct, product jewerly.CreateProductInput) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"invalid input body"}`,
		},
		{
			name:                 "Missing Code",
			fixturePath:          "./fixtures/products/create.no_code.json", // /pkg/handler/
			mockBehavior:         func(r *mock_service.MockProduct, product jewerly.CreateProductInput) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"invalid input body"}`,
		},
		{
			name:                 "Empty Ids",
			fixturePath:          "./fixtures/products/create.empty_ids.json", // /pkg/handler/
			mockBehavior:         func(r *mock_service.MockProduct, product jewerly.CreateProductInput) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"invalid input body"}`,
		},
		{
			name:                 "Missing Ids",
			fixturePath:          "./fixtures/products/create.no_ids.json", // /pkg/handler/
			mockBehavior:         func(r *mock_service.MockProduct, product jewerly.CreateProductInput) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"invalid input body"}`,
		},
		{
			name:                 "Missing Category",
			fixturePath:          "./fixtures/products/create.no_category.json", // /pkg/handler/
			mockBehavior:         func(r *mock_service.MockProduct, product jewerly.CreateProductInput) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"invalid input body"}`,
		},
		{
			name:                 "Invalid Category",
			fixturePath:          "./fixtures/products/create.invalid_category.json", // /pkg/handler/
			mockBehavior:         func(r *mock_service.MockProduct, product jewerly.CreateProductInput) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"invalid category"}`,
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

func TestHandler_updateProduct(t *testing.T) {
	// Init Test Data
	type mockBehavior func(r *mock_service.MockProduct, id int, product jewerly.UpdateProductInput)

	testCases := []struct {
		name                 string
		fixturePath          string
		inputProduct         jewerly.UpdateProductInput
		id                   int
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "Ok",
			fixturePath: "./fixtures/products/update.ok.json",
			id:          1,
			inputProduct: jewerly.UpdateProductInput{
				Titles: &jewerly.MultiLanguageInput{
					English:   "Product 1",
					Ukrainian: "Продукт 1",
					Russian:   "Продукт 1",
				},
				Descriptions: &jewerly.MultiLanguageInput{
					English:   "Description",
					Ukrainian: "Опис",
					Russian:   "Описание",
				},
				Material: &jewerly.MultiLanguageInput{
					English:   "Material",
					Ukrainian: "Матеріал",
					Russian:   "Материал",
				},
				Price:      null.NewFloat(199.99, true),
				Code:       null.NewString("ABC123", true),
				CategoryId: newCategory(jewerly.CategoryBracelets),
			},
			mockBehavior: func(r *mock_service.MockProduct, id int, product jewerly.UpdateProductInput) {
				r.EXPECT().Update(id, product).Return(nil)
			},
			expectedStatusCode: 200,
		},
		{
			name:        "Missing Titles",
			fixturePath: "./fixtures/products/update.no_titles.json",
			id:          1,
			inputProduct: jewerly.UpdateProductInput{
				Descriptions: &jewerly.MultiLanguageInput{
					English:   "Description",
					Ukrainian: "Опис",
					Russian:   "Описание",
				},
				Material: &jewerly.MultiLanguageInput{
					English:   "Material",
					Ukrainian: "Матеріал",
					Russian:   "Материал",
				},
				Price:      null.NewFloat(199.99, true),
				Code:       null.NewString("ABC123", true),
				CategoryId: newCategory(jewerly.CategoryBracelets),
			},
			mockBehavior: func(r *mock_service.MockProduct, id int, product jewerly.UpdateProductInput) {
				r.EXPECT().Update(id, product).Return(nil)
			},
			expectedStatusCode: 200,
		},
		{
			name:        "Missing Description",
			fixturePath: "./fixtures/products/update.no_description.json",
			id:          1,
			inputProduct: jewerly.UpdateProductInput{
				Material: &jewerly.MultiLanguageInput{
					English:   "Material",
					Ukrainian: "Матеріал",
					Russian:   "Материал",
				},
				Price:      null.NewFloat(199.99, true),
				Code:       null.NewString("ABC123", true),
				CategoryId: newCategory(jewerly.CategoryBracelets),
			},
			mockBehavior: func(r *mock_service.MockProduct, id int, product jewerly.UpdateProductInput) {
				r.EXPECT().Update(id, product).Return(nil)
			},
			expectedStatusCode: 200,
		},
		{
			name:        "Missing Materials",
			fixturePath: "./fixtures/products/update.no_materials.json",
			id:          1,
			inputProduct: jewerly.UpdateProductInput{
				Price:      null.NewFloat(199.99, true),
				Code:       null.NewString("ABC123", true),
				CategoryId: newCategory(jewerly.CategoryBracelets),
			},
			mockBehavior: func(r *mock_service.MockProduct, id int, product jewerly.UpdateProductInput) {
				r.EXPECT().Update(id, product).Return(nil)
			},
			expectedStatusCode: 200,
		},
		{
			name:        "Missing Price",
			fixturePath: "./fixtures/products/update.no_price.json",
			id:          1,
			inputProduct: jewerly.UpdateProductInput{
				Code:       null.NewString("ABC123", true),
				CategoryId: newCategory(jewerly.CategoryBracelets),
			},
			mockBehavior: func(r *mock_service.MockProduct, id int, product jewerly.UpdateProductInput) {
				r.EXPECT().Update(id, product).Return(nil)
			},
			expectedStatusCode: 200,
		},
		{
			name:                 "Invalid Price",
			fixturePath:          "./fixtures/products/update.invalid_price.json",
			mockBehavior:         func(r *mock_service.MockProduct, id int, product jewerly.UpdateProductInput) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"price can't be negative or zero"}`,
		},
		{
			name:        "Missing Code",
			fixturePath: "./fixtures/products/update.no_code.json",
			inputProduct: jewerly.UpdateProductInput{
				CategoryId: newCategory(jewerly.CategoryBracelets),
			},
			mockBehavior: func(r *mock_service.MockProduct, id int, product jewerly.UpdateProductInput) {
				r.EXPECT().Update(id, product).Return(nil)
			},
			expectedStatusCode: 200,
		},
		{
			name:        "Empty",
			fixturePath: "./fixtures/products/update.empty.json",
			inputProduct: jewerly.UpdateProductInput{
				CategoryId: newCategory(jewerly.CategoryBracelets),
			},
			mockBehavior:         func(r *mock_service.MockProduct, id int, product jewerly.UpdateProductInput) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"empty update product input"}`,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// Init Deps
			c := gomock.NewController(t)
			defer c.Finish()

			product := mock_service.NewMockProduct(c)
			test.mockBehavior(product, test.id, test.inputProduct)

			services := &service.Services{Product: product}
			handler := Handler{services}

			// Init Endpoint
			r := gin.New()
			r.PUT("/product/:id", handler.updateProduct)

			// Input Body
			fixtureData, err := ioutil.ReadFile(test.fixturePath)
			assert.NoError(t, err)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("PUT", fmt.Sprintf("/product/%d", test.id),
				bytes.NewBuffer(fixtureData))

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_deleteProduct(t *testing.T) {
	// Init Test Data
	type mockBehavior func(r *mock_service.MockProduct, id int)

	testCases := []struct {
		name                 string
		id                   int
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "Ok",
			id:   1,
			mockBehavior: func(r *mock_service.MockProduct, id int) {
				r.EXPECT().Delete(id).Return(nil)
			},
			expectedStatusCode: 200,
		},
		{
			name:                 "Id Zero",
			mockBehavior:         func(r *mock_service.MockProduct, id int) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"id can't be zero"}`,
		},
		{
			name: "Service Error",
			id:   1,
			mockBehavior: func(r *mock_service.MockProduct, id int) {
				r.EXPECT().Delete(id).Return(errors.New("no rows"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"error":"no rows"}`,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// Init Deps
			c := gomock.NewController(t)
			defer c.Finish()

			product := mock_service.NewMockProduct(c)
			test.mockBehavior(product, test.id)

			services := &service.Services{Product: product}
			handler := Handler{services}

			// Init Endpoint
			r := gin.New()
			r.DELETE("/product/:id", handler.deleteProduct)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", fmt.Sprintf("/product/%d", test.id), nil)

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_getProduct(t *testing.T) {
	// Init Test Data
	type mockBehavior func(r *mock_service.MockProduct, product jewerly.ProductResponse, id int, language string)

	testCases := []struct {
		name                 string
		id                   int
		language             string
		languageQuery        string
		product              jewerly.ProductResponse
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:          "Ok",
			id:            1,
			language:      jewerly.English,
			languageQuery: "en",
			mockBehavior: func(r *mock_service.MockProduct, product jewerly.ProductResponse, id int, language string) {
				r.EXPECT().GetById(id, language).Return(product, nil)
			},
			product: jewerly.ProductResponse{
				Id:          1,
				Title:       "product",
				Description: "description",
				Material:    "material",
				Price:       199.99,
				Code:        null.NewString("ABC123", true),
				Images: []jewerly.Image{
					{
						Id:  1,
						URL: "http://image",
					},
				},
				CategoryId: jewerly.CategoryRings,
				InStock:    true,
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"id":1,"title":"product","description":"description","material":"material","price":199.99,"code":"ABC123","images":[{"id":1,"url":"http://image","alt_text":null}],"category_id":1,"in_stock":true}`,
		},
		{
			name:     "No Language Query",
			id:       1,
			language: jewerly.English,
			mockBehavior: func(r *mock_service.MockProduct, product jewerly.ProductResponse, id int, language string) {
				r.EXPECT().GetById(id, language).Return(product, nil)
			},
			product: jewerly.ProductResponse{
				Id:          1,
				Title:       "product",
				Description: "description",
				Material:    "material",
				Price:       199.99,
				Code:        null.NewString("ABC123", true),
				Images: []jewerly.Image{
					{
						Id:  1,
						URL: "http://image",
					},
				},
				CategoryId: jewerly.CategoryRings,
				InStock:    true,
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"id":1,"title":"product","description":"description","material":"material","price":199.99,"code":"ABC123","images":[{"id":1,"url":"http://image","alt_text":null}],"category_id":1,"in_stock":true}`,
		},
		{
			name:                 "Id is 0",
			mockBehavior:         func(r *mock_service.MockProduct, product jewerly.ProductResponse, id int, language string) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"id can't be zero"}`,
		},
		{
			name:     "Service Error",
			id:       1,
			language: jewerly.English,
			mockBehavior: func(r *mock_service.MockProduct, product jewerly.ProductResponse, id int, language string) {
				r.EXPECT().GetById(id, language).Return(product, errors.New("failed to get product"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"error":"failed to get product"}`,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// Init Deps
			c := gomock.NewController(t)
			defer c.Finish()

			product := mock_service.NewMockProduct(c)
			test.mockBehavior(product, test.product, test.id, test.language)

			services := &service.Services{Product: product}
			handler := Handler{services}

			// Init Endpoint
			r := gin.New()
			r.GET("/product/:id", handler.getProduct)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", fmt.Sprintf("/product/%d?language=%s", test.id, test.languageQuery), nil)

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_uploadImage(t *testing.T) {
	// Init Test Data
	type uploadInput struct {
		size        int64
		contentType string
	}

	type mockBehavior func(r *mock_service.MockProduct, input uploadInput, id int)

	testCases := []struct {
		name                 string
		filePath             string
		uploadInput          uploadInput
		imageId              int
		mockBehavior         mockBehavior
		expectedContentType  string
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:                "Ok PNG",
			filePath:            "./fixtures/images/ok.png",
			expectedContentType: "image/png",
			imageId:             1,
			mockBehavior: func(r *mock_service.MockProduct, input uploadInput, id int) {
				r.EXPECT().UploadImage(gomock.Any(), gomock.Any(), input.size, input.contentType).Return(id, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"id":1}`,
		},
		{
			name:                "Ok JPEG",
			filePath:            "./fixtures/images/ok2.jpg",
			expectedContentType: "image/jpeg",
			imageId:             1,
			mockBehavior: func(r *mock_service.MockProduct, input uploadInput, id int) {
				r.EXPECT().UploadImage(gomock.Any(), gomock.Any(), input.size, input.contentType).Return(id, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"id":1}`,
		},
		{
			name:                "File Too Large (> 5MB)",
			filePath:            "./fixtures/images/large.png",
			mockBehavior: func(r *mock_service.MockProduct, input uploadInput, id int) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"http: request body too large"}`,
		},
		{
			name:                "Wrong File Format",
			filePath:            "./fixtures/images/wrong.gif",
			mockBehavior: func(r *mock_service.MockProduct, input uploadInput, id int) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"invalid file type"}`,
		},
		{
			name:                "Service Error",
			filePath:            "./fixtures/images/ok2.jpg",
			expectedContentType: "image/jpeg",
			imageId:             1,
			mockBehavior: func(r *mock_service.MockProduct, input uploadInput, id int) {
				r.EXPECT().UploadImage(gomock.Any(), gomock.Any(), input.size, input.contentType).Return(id, errors.New("error processing image"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"error":"error processing image"}`,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// Init Deps
			c := gomock.NewController(t)
			defer c.Finish()

			// Init Upload Image
			file, err := os.Open(test.filePath)
			assert.NoError(t, err)

			defer file.Close()

			stat, err := file.Stat()
			assert.NoError(t, err)

			test.uploadInput.size = stat.Size()
			test.uploadInput.contentType = test.expectedContentType

			product := mock_service.NewMockProduct(c)
			test.mockBehavior(product, test.uploadInput, test.imageId)

			services := &service.Services{Product: product}
			handler := Handler{services}

			// Init Endpoint
			r := gin.New()
			r.POST("/image", handler.uploadImage)

			// Create Request
			w := httptest.NewRecorder()
			req, err := uploadRequest("/image", file, stat)
			assert.NoError(t, err)

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedResponseBody, w.Body.String())
		})
	}
}

func newCategory(category jewerly.Category) *jewerly.Category {
	return &category
}

func uploadRequest(url string, file *os.File, stat os.FileInfo) (*http.Request, error) {
	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("image", stat.Name())
	if err != nil {
		return nil, err
	}

	part.Write(fileContents)

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return &http.Request{}, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	return req, nil
}
