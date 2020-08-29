package service

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	jewerly "github.com/zhashkevych/jewelry-shop-backend"
	"github.com/zhashkevych/jewelry-shop-backend/pkg/payment"
	"github.com/zhashkevych/jewelry-shop-backend/pkg/repository"
)

const (
	defaultCurrency = "USD"
	timeFomrat      = "2006-01-02 15:04:05"
)

type OrderService struct {
	repo            repository.Order
	paymentProvider payment.Provider
}

func NewOrderService(repo repository.Order, paymentProvider payment.Provider) *OrderService {
	return &OrderService{repo: repo, paymentProvider: paymentProvider}
}

func (s *OrderService) Create(input jewerly.CreateOrderInput) (string, error) {
	totalCost, err := s.getOrderTotalCost(input.Items)
	if err != nil {
		logrus.Errorf("failed to calculate total order cost: %s", err.Error())
		return "", err
	}
	input.TotalCost = totalCost

	transactionId, err := s.generateTransactionId()
	if err != nil {
		logrus.Errorf("failed to generate transactionID: %s", err.Error())
		return "", err
	}
	input.TransactionID = transactionId

	// create order & transaction
	orderId, err := s.repo.Create(input)
	if err != nil {
		logrus.Errorf("failed to create order & transaction: %s", err.Error())
		return "", err
	}

	// generate form with transaction id
	url, err := s.paymentProvider.GenerateSale(payment.GenerateSaleInput{
		Price:         int(input.TotalCost * 100),
		ProductName:   fmt.Sprintf("Order #%d", orderId),
		TransactionID: input.TransactionID,
		Currency:      defaultCurrency, // todo: implement currency input
	})
	if err != nil {
		logrus.Errorf("failed to generate sale form: %s", err.Error())
		return "", err
	}

	url = fmt.Sprintf("%s?first_name=%s&last_name=%s&phone=%s&email=%s&zip_code=%s",
		url, input.FirstName, input.LastName, input.Phone, input.Email, input.PostalCode)

	return url, nil
}

func (s *OrderService) ProcessCallback(inp jewerly.TransactionCallbackInput) error {
	return s.repo.CreateTransaction(inp.TransactionID, inp.BuyerCardMask, inp.NotifyType)
}

func (s *OrderService) GetAll(input jewerly.GetAllOrdersFilters) (jewerly.OrderList, error) {
	return s.repo.GetAll(input)
}

func (s *OrderService) GetById(id int) (jewerly.Order, error) {
	return s.repo.GetById(id)
}

func (s *OrderService) getOrderTotalCost(orderItems []jewerly.OrderItem) (float32, error) {
	products, err := s.repo.GetOrderProducts(orderItems)
	if err != nil {
		return 0, err
	}

	productsPriceList := make(map[int]float32)
	for _, product := range products {
		productsPriceList[product.Id] = product.CurrentPrice
	}

	var totalCost float32
	for _, item := range orderItems {
		totalCost += productsPriceList[item.ProductId] * float32(item.Quantity)
	}

	return totalCost, nil
}

func (s *OrderService) generateTransactionId() (string, error) {
	transactionId, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}

	return transactionId.String(), nil
}
