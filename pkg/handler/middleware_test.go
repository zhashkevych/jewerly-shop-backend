package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/zhashkevych/jewelry-shop-backend/pkg/service"
	mock_service "github.com/zhashkevych/jewelry-shop-backend/pkg/service/mocks"
	"net/http/httptest"
	"testing"
)

func TestHandler_adminIdentity(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mock_service.MockAdmin, token string)

	testTable := []struct {
		name                 string
		headerName           string
		headerValue          string
		token                string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "Ok",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehavior: func(r *mock_service.MockAdmin, token string) {
				r.EXPECT().ParseToken(token).Return( nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: "ok",
		},
		{
			name:                 "Invalid Header Name",
			headerName:           "",
			headerValue:          "Bearer token",
			token:                "token",
			mockBehavior:         func(r *mock_service.MockAdmin, token string) {},
			expectedStatusCode:   401,
			expectedResponseBody: `{"error":"empty auth header"}`,
		},
		{
			name:                 "Invalid Header Value",
			headerName:           "Authorization",
			headerValue:          "Bearr token",
			token:                "token",
			mockBehavior:         func(r *mock_service.MockAdmin, token string) {},
			expectedStatusCode:   401,
			expectedResponseBody: `{"error":"invalid auth header"}`,
		},
		{
			name:                 "Empty Token",
			headerName:           "Authorization",
			headerValue:          "Bearer ",
			token:                "token",
			mockBehavior:         func(r *mock_service.MockAdmin, token string) {},
			expectedStatusCode:   401,
			expectedResponseBody: `{"error":"invalid token"}`,
		},
		{
			name:        "Parse Error",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehavior: func(r *mock_service.MockAdmin, token string) {
				r.EXPECT().ParseToken(token).Return(errors.New("invalid token"))
			},
			expectedStatusCode:   401,
			expectedResponseBody: `{"error":"invalid token"}`,
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			admin := mock_service.NewMockAdmin(c)
			test.mockBehavior(admin, test.token)

			services := &service.Services{Admin: admin}
			handler := Handler{services}

			// Init Endpoint
			r := gin.New()
			r.GET("/identity", handler.adminIdentity, func(c *gin.Context) {
				c.String(200, "ok")
			})

			// Init Test Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/identity", nil)
			req.Header.Set(test.headerName, test.headerValue)

			r.ServeHTTP(w, req)

			// Asserts
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}
