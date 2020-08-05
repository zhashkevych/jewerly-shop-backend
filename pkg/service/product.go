package service

import (
	jewerly "github.com/zhashkevych/jewelry-shop-backend"
	"github.com/zhashkevych/jewelry-shop-backend/pkg/repository"
)

type ProductService struct {
	repo repository.Product
}

func NewProductService(repo repository.Product) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) Create(product jewerly.CreateProductInput) error {
	return s.repo.Create(product)
}

func (s *ProductService) GetAll(filters jewerly.GetAllProductsFilters) (jewerly.ProductsList, error) {
	return s.repo.GetAll(filters)
}

func (s *ProductService) Delete(id int) error {
	return s.repo.Delete(id)
}

func (s *ProductService) GetById(id int, language string) (jewerly.ProductResponse, error) {
	return s.repo.GetById(id, language)
}
