package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/magiconair/properties/assert"
	jewerly "github.com/zhashkevych/jewelry-shop-backend"
	"gopkg.in/guregu/null.v3"
	"net/http"
	"net/url"
	"testing"
)

func Test_getProductFilters(t *testing.T) {
	testTable := []struct {
		name          string
		language      string
		limit, offset string
		category      string
		expected      jewerly.GetAllProductsFilters
	}{
		{
			name:     "Ok",
			language: "en",
			limit:    "10",
			offset:   "10",
			category: "1",
			expected: jewerly.GetAllProductsFilters{
				Language:   jewerly.English,
				Offset:     10,
				Limit:      10,
				CategoryId: null.NewInt(1, true),
			},
		},
		{
			name:     "Ok - Empty Language",
			language: "",
			limit:    "10",
			offset:   "10",
			category: "1",
			expected: jewerly.GetAllProductsFilters{
				Language:   jewerly.English,
				Offset:     10,
				Limit:      10,
				CategoryId: null.NewInt(1, true),
			},
		},
		{
			name:     "Ok - Empty Limit",
			language: "",
			limit:    "",
			offset:   "10",
			category: "1",
			expected: jewerly.GetAllProductsFilters{
				Language:   jewerly.English,
				Offset:     10,
				Limit:      20,
				CategoryId: null.NewInt(1, true),
			},
		},
		{
			name:     "Ok - Empty Offset",
			language: "",
			limit:    "",
			offset:   "",
			category: "1",
			expected: jewerly.GetAllProductsFilters{
				Language:   jewerly.English,
				Offset:     0,
				Limit:      20,
				CategoryId: null.NewInt(1, true),
			},
		},
		{
			name:     "Ok - Empty Category",
			language: "",
			limit:    "",
			offset:   "",
			category: "",
			expected: jewerly.GetAllProductsFilters{
				Language:   jewerly.English,
				Offset:     0,
				Limit:      20,
			},
		},
		{
			name:     "Ok - Language ru",
			language: "ru",
			limit:    "",
			offset:   "",
			category: "",
			expected: jewerly.GetAllProductsFilters{
				Language:   jewerly.Russian,
				Offset:     0,
				Limit:      20,
			},
		},
		{
			name:     "Ok - Language ua",
			language: "ua",
			limit:    "",
			offset:   "",
			category: "",
			expected: jewerly.GetAllProductsFilters{
				Language:   jewerly.Ukraininan,
				Offset:     0,
				Limit:      20,
			},
		},
		{
			name:     "Ok - Zero Limit",
			language: "ua",
			limit:    "0",
			offset:   "",
			category: "",
			expected: jewerly.GetAllProductsFilters{
				Language:   jewerly.Ukraininan,
				Offset:     0,
				Limit:      20,
			},
		},
		{
			name:     "Ok - Negative Limit",
			language: "ua",
			limit:    "-10",
			offset:   "",
			category: "",
			expected: jewerly.GetAllProductsFilters{
				Language:   jewerly.Ukraininan,
				Offset:     0,
				Limit:      20,
			},
		},
		{
			name:     "Ok - Negative Offset",
			language: "ua",
			limit:    "",
			offset:   "-10",
			category: "",
			expected: jewerly.GetAllProductsFilters{
				Language:   jewerly.Ukraininan,
				Offset:     0,
				Limit:      20,
			},
		},
		{
			name:     "Ok - Incorrect Category",
			language: "ua",
			limit:    "0",
			offset:   "",
			category: "10",
			expected: jewerly.GetAllProductsFilters{
				Language:   jewerly.Ukraininan,
				Offset:     0,
				Limit:      20,
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctx := &gin.Context{
				Request: &http.Request{
					URL: &url.URL{
						RawQuery: fmt.Sprintf("language=%s&limit=%s&offset=%s&category=%s",
							testCase.language, testCase.limit, testCase.offset, testCase.category),
					},
				},
			}

			result := getProductFilters(ctx)

			assert.Equal(t, result, testCase.expected)
		})
	}
}

func Test_getOrderFilters(t *testing.T) {
	testTable := []struct {
		name          string
		limit, offset string
		expected      jewerly.GetAllOrdersFilters
	}{
		{
			name:     "Ok",
			limit:    "10",
			offset:   "10",
			expected: jewerly.GetAllOrdersFilters{
				Offset:     10,
				Limit:      10,
			},
		},
		{
			name:     "Ok - Empty Limit",
			limit:    "",
			offset:   "10",
			expected: jewerly.GetAllOrdersFilters{
				Offset:     10,
				Limit:      20,
			},
		},
		{
			name:     "Ok - Empty Offset",
			limit:    "",
			offset:   "",
			expected: jewerly.GetAllOrdersFilters{
				Offset:     0,
				Limit:      20,
			},
		},
		{
			name:     "Ok - Zero Limit",
			limit:    "0",
			offset:   "",
			expected: jewerly.GetAllOrdersFilters{
				Offset:     0,
				Limit:      20,
			},
		},
		{
			name:     "Ok - Negative Limit",
			limit:    "-10",
			offset:   "",
			expected: jewerly.GetAllOrdersFilters{
				Offset:     0,
				Limit:      20,
			},
		},
		{
			name:     "Ok - Negative Offset",
			limit:    "10",
			offset:   "-10",
			expected: jewerly.GetAllOrdersFilters{
				Offset:     0,
				Limit:      10,
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctx := &gin.Context{
				Request: &http.Request{
					URL: &url.URL{
						RawQuery: fmt.Sprintf("limit=%s&offset=%s", testCase.limit, testCase.offset),
					},
				},
			}

			result := getOrderFilters(ctx)

			assert.Equal(t, result, testCase.expected)
		})
	}
}