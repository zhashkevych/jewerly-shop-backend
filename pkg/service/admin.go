package service

import (
	"crypto/sha1"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/zhashkevych/jewelry-shop-backend/pkg/repository"
	"time"
)

const (
	tokenTTL = 12 * time.Hour
)

type AdminService struct {
	repo       repository.Admin
	salt       string
	signingKey []byte
}

func NewAdminService(repo repository.Admin, salt string, signingKey []byte) *AdminService {
	return &AdminService{repo: repo, salt: salt, signingKey: signingKey}
}

func (s *AdminService) SignIn(login, password string) (string, error) {
	err := s.repo.Authorize(login, s.getPasswordHash(password))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(tokenTTL).Unix(),
		IssuedAt:  time.Now().Unix(),
	})

	return token.SignedString(s.signingKey)
}

func (s *AdminService) ParseToken(token string) error {
	t, err := jwt.Parse(token, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.signingKey, nil
	})

	if err != nil {
		return err
	}

	_, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return fmt.Errorf("error get user claims from token")
	}

	return nil
}

func (s *AdminService) getPasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(s.salt))

	return fmt.Sprintf("%x", sha1.Sum([]byte(password)))
}
