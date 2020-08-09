package service

import (
	"crypto/sha1"
	"errors"
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

func (a *Authorization) ParseToken(token string) (jewerly.User, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return a.signingKey, nil
	})

	if err != nil {
		return jewerly.User{}, err
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return jewerly.User{}, fmt.Errorf("error get user claims from token")
	}

	sub, ex := claims["sub"].(string)
	if !ex {
		return jewerly.User{}, errors.New("token is invalid")
	}

	id, err := strconv.ParseInt(sub, 10, 64)
	if err != nil {
		return jewerly.User{}, fmt.Errorf("error convert user id from string to int: err `%s`", err)
	}

	return jewerly.User{Id: id}, nil
}

func (a *Authorization) getPasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(a.salt))

	return fmt.Sprintf("%x", sha1.Sum([]byte(password)))
}
