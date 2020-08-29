package service

import (
	"context"
	jewerly "github.com/zhashkevych/jewelry-shop-backend"
	"github.com/zhashkevych/jewelry-shop-backend/pkg/payment"
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

type Admin interface {
	SignIn(login, password string) (string, error)
	ParseToken(token string) error
}

// Users
type User interface {
	GetById(id int64) (jewerly.User, error)
}

type Product interface {
	Create(jewerly.CreateProductInput) error
	GetAll(jewerly.GetAllProductsFilters) (jewerly.ProductsList, error)
	GetById(id int, language string) (jewerly.ProductResponse, error)
	Update(id int, inp jewerly.UpdateProductInput) error
	Delete(id int) error
	UploadImage(ctx context.Context, file io.Reader, size int64, contentType string) (int, error)
}

type Order interface {
	Create(jewerly.CreateOrderInput) (string, error)
	ProcessCallback(jewerly.TransactionCallbackInput) error
	GetAll(jewerly.GetAllOrdersFilters) (jewerly.OrderList, error)
	GetById(id int) (jewerly.Order, error)
}

// Services Interface, Constructor & Dependencies
type Dependencies struct {
	Repos           *repository.Repository
	FileStorage     storage.Storage
	HashSalt        string
	SigningKey      []byte
	PaymentProvider payment.Provider
}

type Services struct {
	Auth
	Admin
	User
	Product
	Order
}

func NewServices(deps Dependencies) *Services {
	return &Services{
		Auth:    NewAuthorization(deps.Repos.User, deps.HashSalt, deps.SigningKey),
		Admin:   NewAdminService(deps.Repos.Admin, deps.HashSalt, deps.SigningKey),
		Product: NewProductService(deps.Repos.Product, deps.FileStorage),
		Order:   NewOrderService(deps.Repos.Order, deps.PaymentProvider),
	}
}
