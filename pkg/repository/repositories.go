package repository

import (
	"github.com/jmoiron/sqlx"
	jewerly "github.com/zhashkevych/jewelry-shop-backend"
)

type User interface {
	Create(user jewerly.User) error
	GetByCredentials(email, passwordHash string) (jewerly.User, error)
	GetById(id int64) (jewerly.User, error)
	//GetAll(id int64) (jewerly.User, error)
	//Update(id int64, newUser jewerly.User) error
}

type Product interface {
	Create(product jewerly.CreateProductInput) error
	GetAll(filters jewerly.GetAllProductsFilters) (jewerly.ProductsList, error)
	GetById(id int, language string) (jewerly.ProductResponse, error)
	Delete(id int) error
	CreateImage(url, altText string) (int, error)
	GetProductImages(productId int) ([]jewerly.Image, error)
}

type Repository struct {
	User
	Product
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		User:    NewUserRepository(db),
		Product: NewProductRepository(db),
	}
}
