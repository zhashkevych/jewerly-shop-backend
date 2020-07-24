package service

import (
	"crypto/sha1"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	jewerly "github.com/zhashkevych/jewelry-shop-backend"
	"github.com/zhashkevych/jewelry-shop-backend/pkg/repository"
	"strconv"
	"time"
)

const (
	tokenTTL = time.Hour * 12
)

type Authorization struct {
	repo       repository.User
	salt       string
	signingKey []byte
}

func NewAuthorization(repo repository.User, salt string, signingKey []byte) *Authorization {
	return &Authorization{
		repo:       repo,
		salt:       salt,
		signingKey: signingKey,
	}
}

func (a *Authorization) SignUp(inp SignUpInput) error {
	user := jewerly.User{
		FirstName:    inp.FirstName,
		LastName:     inp.LastName,
		Email:        inp.Email,
		PasswordHash: a.getPasswordHash(inp.Password),
	}

	return a.repo.Create(user)
}

func (a *Authorization) SignIn(email, password string) (string, error) {
	user, err := a.repo.GetByCredentials(email, a.getPasswordHash(password))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(tokenTTL).Unix(),
		IssuedAt:  time.Now().Unix(),
		Subject:   strconv.FormatInt(user.Id, 10),
	})

	return token.SignedString(a.signingKey)
}

func (a *Authorization) getPasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(a.salt))

	return fmt.Sprintf("%x", sha1.Sum([]byte(password)))
}
