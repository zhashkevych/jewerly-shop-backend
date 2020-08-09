package service

import "github.com/zhashkevych/jewelry-shop-backend/pkg/repository"

type OrderService struct {
	repo repository.Order
}

func NewOrderService(repo repository.Order) *OrderService {
	return &OrderService{repo: repo}
}

func (s *OrderService) Create(userId int64, productIds []int) error {
	return s.repo.Create(userId, productIds)
}

