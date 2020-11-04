package service

import (
	"context"
	"github.com/hashicorp/go-uuid"
	"github.com/sirupsen/logrus"
	jewerly "github.com/zhashkevych/jewelry-shop-backend"
	"github.com/zhashkevych/jewelry-shop-backend/pkg/repository"
	"github.com/zhashkevych/jewelry-shop-backend/pkg/storage"
	"io"
)

type ProductService struct {
	repo        repository.Product
	fileStorage storage.Storage
}

func NewProductService(repo repository.Product, fileStorage storage.Storage) *ProductService {
	return &ProductService{repo: repo, fileStorage: fileStorage}
}

func (s *ProductService) Create(product jewerly.CreateProductInput) error {
	return s.repo.Create(product)
}

func (s *ProductService) GetAll(filters jewerly.GetAllProductsFilters) (jewerly.ProductsList, error) {
	productList, err := s.repo.GetAll(filters)
	if err != nil {
		return productList, err
	}

	for i, product := range productList.Products {
		images, err := s.repo.GetProductImages(product.Id)
		if err != nil {
			logrus.Errorf("failed to get images for product id %d: %s", product.Id, err.Error())
			continue
		}

		productList.Products[i].Images = images
	}

	return productList, nil
}

func (s *ProductService) Update(id int, inp jewerly.UpdateProductInput) error {
	return s.repo.Update(id, inp)
}

func (s *ProductService) Delete(id int) error {
	return s.repo.Delete(id)
}

func (s *ProductService) GetById(id int, language string) (jewerly.ProductResponse, error) {
	product, err := s.repo.GetById(id, language)
	if err != nil {
		return product, err
	}

	images, err := s.repo.GetProductImages(product.Id)
	if err != nil {
		logrus.Errorf("failed to get images for product id %d: %s", product.Id, err.Error())
		return product, err
	}

	product.Images = images

	return product, nil
}

func (s *ProductService) UploadImage(ctx context.Context, file io.Reader, size int64, contentType string) (int, error) {
	filename, err := generateFileName()
	if err != nil {
		return 0, err
	}

	url, err := s.fileStorage.Upload(ctx, storage.UploadInput{
		File:        file,
		Name:        filename,
		Size:        size,
		ContentType: contentType,
	})
	if err != nil {
		return 0, err
	}

	return s.repo.CreateImage(url, "")
}

func generateFileName() (string, error) {
	return uuid.GenerateUUID()
}
