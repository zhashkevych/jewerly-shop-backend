package repository

import (
	"github.com/jmoiron/sqlx"
	jewerly "github.com/zhashkevych/jewelry-shop-backend"
	"github.com/zhashkevych/jewelry-shop-backend/pkg/repository/postgres"
)

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

type Settings interface {
	GetImages() ([]jewerly.HomepageImage, error)
	CreateImage(imageID int) error
	UpdateImage(id, imageID int) error

	GetTextBlocks() ([]jewerly.TextBlock, error)
	GetTextBlockById(id int) (jewerly.TextBlock, error)
	CreateTextBlock(block jewerly.TextBlock) error
	UpdateTextBlock(id int, block jewerly.UpdateTextBlockInput) error
}

type PageText interface {
	Create(page string, input jewerly.MultiLanguageInput) error
	Update()
	Get()
}

type Repository struct {
	Admin
	Product
	Order
	Settings
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Admin:    postgres.NewAdminRepository(db),
		Product:  postgres.NewProductRepository(db),
		Order:    postgres.NewOrderRepository(db),
		Settings: postgres.NewSettingsRepository(db),
	}
}
