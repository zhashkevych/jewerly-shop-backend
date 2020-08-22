package repository

import (
	"github.com/jmoiron/sqlx"
	jewerly "github.com/zhashkevych/jewelry-shop-backend"
)

type User interface {
	Create(user jewerly.User) error
	GetByCredentials(email, passwordHash string) (jewerly.User, error)
	GetById(id int64) (jewerly.User, error)
}

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
}

type Repository struct {
	User
	Admin
	Product
	Order
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		User:    NewUserRepository(db),
		Admin:   NewAdminRepository(db),
		Product: NewProductRepository(db),
		Order:   NewOrderRepository(db),
	}
}
