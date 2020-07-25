package service

import (
	jewerly "github.com/zhashkevych/jewelry-shop-backend"
	"github.com/zhashkevych/jewelry-shop-backend/pkg/repository"
)

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

type User interface {
	GetById(id int64) (jewerly.User, error)
}

type Dependencies struct {
	Repos      *repository.Repository
	HashSalt   string
	SigningKey []byte
}

type Service struct {
	Auth
	User
}

func NewService(deps Dependencies) *Service {
	return &Service{
		Auth: NewAuthorization(deps.Repos.User, deps.HashSalt, deps.SigningKey),
	}
}
