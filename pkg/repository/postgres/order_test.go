package postgres

import (
	"database/sql/driver"
	"errors"
	"github.com/stretchr/testify/assert"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
	jewerly "github.com/zhashkevych/jewelry-shop-backend"
	"testing"
)

func TestOrderRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	type mockBehavior func(input jewerly.CreateOrderInput, orderId int)

	testTable := []struct {
		name         string
		input        jewerly.CreateOrderInput
		orderId      int
		mockBehavior mockBehavior
		shouldFail   bool
	}{
		{
			name: "OK",
			input: jewerly.CreateOrderInput{
				Items: []jewerly.OrderItem{
					{1, 3},
					{18, 5},
				},
				FirstName:      "Test",
				LastName:       "Test",
				AdditionalName: "",
				Email:          "test@test.com",
				Phone:          "+380506663212",
				Country:        "UA",
				Address:        "Kreshatyk st.",
				PostalCode:     "32012",
				TransactionID:  "1111-2222-3333-4444-asdas",
			},
			orderId: 42,
			mockBehavior: func(input jewerly.CreateOrderInput, orderId int) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id"}).AddRow(orderId)
				mock.ExpectQuery("INSERT INTO orders").WithArgs(input.FirstName, input.LastName, input.AdditionalName, input.Country, input.Address,
					input.PostalCode, input.Email, input.TotalCost).WillReturnRows(rows)

				args := []driver.Value{orderId}
				for _, item := range input.Items {
					args = append(args, item.ProductId, item.Quantity)
				}
				mock.ExpectExec("INSERT INTO order_items").WithArgs(args...).WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectExec("INSERT INTO transactions").WithArgs(orderId, input.TransactionID).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectExec("INSERT INTO transactions_history").WithArgs(input.TransactionID).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
		},
		{
			name: "Insert Order Error",
			input: jewerly.CreateOrderInput{
				Items: []jewerly.OrderItem{
					{1, 3},
					{18, 5},
				},
				FirstName:      "Test",
				LastName:       "Test",
				AdditionalName: "",
				Email:          "test@test.com",
				Phone:          "+380506663212",
				Country:        "UA",
				Address:        "Kreshatyk st.",
				PostalCode:     "32012",
				TransactionID:  "1111-2222-3333-4444-asdas",
			},
			orderId: 42,
			mockBehavior: func(input jewerly.CreateOrderInput, orderId int) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id"}).AddRow(orderId).CloseError(errors.New("fail"))
				mock.ExpectQuery("INSERT INTO orders").WithArgs(input.FirstName, input.LastName, input.AdditionalName, input.Country, input.Address,
					input.PostalCode, input.Email, input.TotalCost).WillReturnRows(rows)

				mock.ExpectRollback()
			},
			shouldFail: true,
		},
		{
			name: "Insert Order Items Fail",
			input: jewerly.CreateOrderInput{
				Items: []jewerly.OrderItem{
					{1, 3},
					{18, 5},
				},
				FirstName:      "Test",
				LastName:       "Test",
				AdditionalName: "",
				Email:          "test@test.com",
				Phone:          "+380506663212",
				Country:        "UA",
				Address:        "Kreshatyk st.",
				PostalCode:     "32012",
				TransactionID:  "1111-2222-3333-4444-asdas",
			},
			orderId: 42,
			mockBehavior: func(input jewerly.CreateOrderInput, orderId int) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id"}).AddRow(orderId)
				mock.ExpectQuery("INSERT INTO orders").WithArgs(input.FirstName, input.LastName, input.AdditionalName, input.Country, input.Address,
					input.PostalCode, input.Email, input.TotalCost).WillReturnRows(rows)

				args := []driver.Value{orderId}
				for _, item := range input.Items {
					args = append(args, item.ProductId, item.Quantity)
				}
				mock.ExpectExec("INSERT INTO order_items").WithArgs(args...).WillReturnError(errors.New("fail"))

				mock.ExpectRollback()
			},
			shouldFail: true,
		},
		{
			name: "Insert Tranasction Fail",
			input: jewerly.CreateOrderInput{
				Items: []jewerly.OrderItem{
					{1, 3},
					{18, 5},
				},
				FirstName:      "Test",
				LastName:       "Test",
				AdditionalName: "",
				Email:          "test@test.com",
				Phone:          "+380506663212",
				Country:        "UA",
				Address:        "Kreshatyk st.",
				PostalCode:     "32012",
				TransactionID:  "1111-2222-3333-4444-asdas",
			},
			orderId: 42,
			mockBehavior: func(input jewerly.CreateOrderInput, orderId int) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id"}).AddRow(orderId)
				mock.ExpectQuery("INSERT INTO orders").WithArgs(input.FirstName, input.LastName, input.AdditionalName, input.Country, input.Address,
					input.PostalCode, input.Email, input.TotalCost).WillReturnRows(rows)

				args := []driver.Value{orderId}
				for _, item := range input.Items {
					args = append(args, item.ProductId, item.Quantity)
				}
				mock.ExpectExec("INSERT INTO order_items").WithArgs(args...).WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectExec("INSERT INTO transactions").WithArgs(orderId, input.TransactionID).
					WillReturnError(errors.New("fail"))

				mock.ExpectRollback()
			},
			shouldFail: true,
		},
		{
			name: "Insert Transaction History Fail",
			input: jewerly.CreateOrderInput{
				Items: []jewerly.OrderItem{
					{1, 3},
					{18, 5},
				},
				FirstName:      "Test",
				LastName:       "Test",
				AdditionalName: "",
				Email:          "test@test.com",
				Phone:          "+380506663212",
				Country:        "UA",
				Address:        "Kreshatyk st.",
				PostalCode:     "32012",
				TransactionID:  "1111-2222-3333-4444-asdas",
			},
			orderId: 42,
			mockBehavior: func(input jewerly.CreateOrderInput, orderId int) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id"}).AddRow(orderId)
				mock.ExpectQuery("INSERT INTO orders").WithArgs(input.FirstName, input.LastName, input.AdditionalName, input.Country, input.Address,
					input.PostalCode, input.Email, input.TotalCost).WillReturnRows(rows)

				args := []driver.Value{orderId}
				for _, item := range input.Items {
					args = append(args, item.ProductId, item.Quantity)
				}
				mock.ExpectExec("INSERT INTO order_items").WithArgs(args...).WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectExec("INSERT INTO transactions").WithArgs(orderId, input.TransactionID).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectExec("INSERT INTO transactions_history").WithArgs(input.TransactionID).
					WillReturnError(errors.New("fail"))

				mock.ExpectRollback()
			},
			shouldFail: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.input, testCase.orderId)

			r := NewOrderRepository(db)

			got, err := r.Create(testCase.input)
			if testCase.shouldFail {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.orderId, got)
			}
		})
	}
}