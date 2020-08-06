package service

import (
	"context"
	jewerly "github.com/zhashkevych/jewelry-shop-backend"
	"github.com/zhashkevych/jewelry-shop-backend/pkg/repository"
	"github.com/zhashkevych/jewelry-shop-backend/storage"
	"io"
)

// Authorization
type SignUpInput struct {
	FirstName string
	LastName  string
	Email     string
	Password  string
}

type Auth interface {
	SignUp(inp SignUpInput) error
	SignIn(email, password string) (string, error)
	ParseToken(token string) (jewerly.User, error)
}

// Users
type User interface {
	GetById(id int64) (jewerly.User, error)
}

type Product interface {
	Create(jewerly.CreateProductInput) error
	GetAll(filters jewerly.GetAllProductsFilters) (jewerly.ProductsList, error)
	GetById(id int, language string) (jewerly.ProductResponse, error)
	Delete(id int) error
	UploadImage(ctx context.Context, file io.Reader, size int64, contentType string) (int, error)
}

// Services Interface, Constructor & Dependencies
type Dependencies struct {
	Repos       *repository.Repository
	FileStorage storage.Storage
	HashSalt    string
	SigningKey  []byte
}

type Services struct {
	Auth
	User
	Product
}

func NewServices(deps Dependencies) *Services {
	return &Services{
		Auth:    NewAuthorization(deps.Repos.User, deps.HashSalt, deps.SigningKey),
		Product: NewProductService(deps.Repos.Product, deps.FileStorage),
	}
}
