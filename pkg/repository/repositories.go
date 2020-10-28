package repository

import (
	"github.com/jmoiron/sqlx"
	jewerly "github.com/zhashkevych/jewelry-shop-backend"
)

// todo refactor package structure
// todo remove user from DB schema

//go:generate mockgen -source=repositories.go -destination=mocks/mock.go

type Admin interface {
	Authorize(login, passwordHash string) error
}

type Product interface {
	Create(product jewerly.CreateProductInput) error
	GetAll(filters jewerly.GetAllProductsFilters) (jewerly.ProductsList, error)
	GetById(id int, language string) (jewerly.ProductResponse, error)
	Update(id int, inp jewerly.UpdateProductInput) error
	Delete(id int) error
	CreateImage(url, altText string) (int, error)
	GetProductImages(productId int) ([]jewerly.Image, error)
}

type Order interface {
	Create(input jewerly.CreateOrderInput) (int, error)
	GetOrderProducts(items []jewerly.OrderItem) ([]jewerly.ProductResponse, error)
	CreateTransaction(transactionId, cardMask, status string) error
	GetOrderId(transactionId string) (int, error)
	GetAll(jewerly.GetAllOrdersFilters) (jewerly.OrderList, error)
	GetById(id int) (jewerly.Order, error)
}

type Repository struct {
	Admin
	Product
	Order
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Admin:   NewAdminRepository(db),
		Product: NewProductRepository(db),
		Order:   NewOrderRepository(db),
	}
}
