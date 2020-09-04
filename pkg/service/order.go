package service

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	jewerly "github.com/zhashkevych/jewelry-shop-backend"
	"github.com/zhashkevych/jewelry-shop-backend/pkg/payment"
	"github.com/zhashkevych/jewelry-shop-backend/pkg/repository"
	"time"
)

const (
	defaultCurrency = "USD"
	timeFomrat      = "2006-01-02 15:04:05"
)

var paymentStatuses = map[string]string{
	"sale-complete":          jewerly.TransactionStatusPaid,
	"sale-authorized":        jewerly.TransactionStatusAuthorized,
	"refund":                 jewerly.TransactionStatusRefunded,
	"sale-failure":           jewerly.TransactionStatusFailed,
	"sale-chargeback":        jewerly.TransactionStatusChargeback,
	"sale-chargeback-refund": jewerly.TransactionStatusReverted,
}

type OrderService struct {
	repo            repository.Order
	paymentProvider payment.Provider
	emailService    Email
}

func NewOrderService(repo repository.Order, paymentProvider payment.Provider, emailService Email) *OrderService {
	return &OrderService{repo: repo, paymentProvider: paymentProvider, emailService: emailService}
}

func (s *OrderService) Create(input jewerly.CreateOrderInput) (string, error) {
	totalCost, products, err := s.getOrderTotalCost(input.Items)
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

	// send order email to support
	go s.sendOrderEmails(jewerly.OrderInfoEmailInput{
		OrderId:           orderId,
		FirstName:         input.FirstName,
		LastName:          input.LastName,
		Country:           input.Country,
		Address:           input.Address,
		PostalCode:        input.PostalCode,
		Email:             input.Email,
		TotalCost:         input.TotalCost,
		TransactionId:     transactionId,
		OrderedAt:         time.Now(),
		TransactionStatus: jewerly.TransactionStatusCreated,
		Products:          createOrderProductsList(input.Items, products),
	})

	url = urlWithParameters(url, input)

	return url, nil
}

func (s *OrderService) ProcessCallback(inp jewerly.TransactionCallbackInput) {
	err := s.repo.CreateTransaction(inp.TransactionID, inp.BuyerCardMask, inp.NotifyType)
	if err != nil {
		logrus.Errorf("failed to create transaction on callback: %s", err.Error())
		return
	}

	go s.sendPaymentEmail(inp)
}

func (s *OrderService) GetAll(input jewerly.GetAllOrdersFilters) (jewerly.OrderList, error) {
	return s.repo.GetAll(input)
}

func (s *OrderService) GetById(id int) (jewerly.Order, error) {
	return s.repo.GetById(id)
}

func (s *OrderService) getOrderTotalCost(orderItems []jewerly.OrderItem) (float32, []jewerly.ProductResponse, error) {
	products, err := s.repo.GetOrderProducts(orderItems)
	if err != nil {
		return 0, products, err
	}

	productsPriceList := make(map[int]float32)
	for _, product := range products {
		productsPriceList[product.Id] = product.CurrentPrice
	}

	var totalCost float32
	for _, item := range orderItems {
		totalCost += productsPriceList[item.ProductId] * float32(item.Quantity)
	}

	return totalCost, products, nil
}

func (s *OrderService) generateTransactionId() (string, error) {
	transactionId, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}

	return transactionId.String(), nil
}

func (s *OrderService) sendOrderEmails(inp jewerly.OrderInfoEmailInput) {
	if err := s.emailService.SendOrderInfoSupport(inp); err != nil {
		logrus.Errorf("failed to send order info email: %s", err.Error())
	}

	if err := s.emailService.SendOrderInfoCustomer(inp); err != nil {
		logrus.Errorf("failed to send order info email: %s", err.Error())
	}
}

func (s *OrderService) sendPaymentEmail(inp jewerly.TransactionCallbackInput) {
	status, err := getPaymentStatus(inp.NotifyType)
	if err != nil {
		logrus.Errorf("transactionId: %s, err: %s", inp.TransactionID, err.Error())
		return
	}

	if status != jewerly.TransactionStatusPaid {
		logrus.Errorf("transactionId: %s, status: %s", inp.TransactionID, status)
		return
	}

	orderId, err := s.repo.GetOrderId(inp.TransactionID)
	if err != nil {
		logrus.Errorf("failed to get order by transaction id: %s", err.Error())
		return
	}

	emailInput := jewerly.PaymentInfoEmailInput{
		TransactionId: inp.TransactionID,
		OrderId:       orderId,
		CardMask:      inp.BuyerCardMask,
		CardBrand:     inp.CardBrand,
		BuyerName:     inp.BuyerName,
		BuyerEmail:    inp.BuyerEmail,
		Price:         float32(inp.Price) / 100,
		Status:        status,
	}

	if err := s.emailService.SendPaymentInfoSupport(emailInput); err != nil {
		logrus.Errorf("failed to send payment info support email: %s", err.Error())
	}

	if err := s.emailService.SendPaymentInfoCustomer(emailInput); err != nil {
		logrus.Errorf("failed to send payment info customer email: %s", err.Error())
	}
}

func createOrderProductsList(orderItems []jewerly.OrderItem, products []jewerly.ProductResponse) []jewerly.ProductInfo {
	quantityList := make(map[int]int)
	for i := range orderItems {
		quantityList[orderItems[i].ProductId] = orderItems[i].Quantity
	}

	items := make([]jewerly.ProductInfo, len(products))
	for i := range products {
		items[i].Id = products[i].Id
		items[i].Title = products[i].Title
		items[i].Price = products[i].CurrentPrice
		items[i].Quantity = quantityList[products[i].Id]

		if len(products[i].Images) > 0 {
			items[i].ImageURL = products[i].Images[0].URL
		}
	}

	return items
}

func urlWithParameters(url string, input jewerly.CreateOrderInput) string {
	return fmt.Sprintf("%s?first_name=%s&last_name=%s&phone=%s&email=%s&zip_code=%s",
		url, input.FirstName, input.LastName, input.Phone, input.Email, input.PostalCode)
}

func getPaymentStatus(notifyType string) (string, error) {
	status, ok := paymentStatuses[notifyType]
	if !ok {
		return jewerly.TransactionStatusFailed, errors.New(fmt.Sprintf("notify type %s missing", notifyType))
	}

	return status, nil
}
