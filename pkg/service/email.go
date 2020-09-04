package service

import (
	"fmt"
	jewerly "github.com/zhashkevych/jewelry-shop-backend"
	"github.com/zhashkevych/jewelry-shop-backend/pkg/email"
	"time"
)

type EmailDeps struct {
	SupportEmail string
	SenderEmail  string
	SupportName  string
	SenderName   string

	OrderInfoSupportTemplate string
	OrderInfoSupportSubject  string

	OrderInfoCustomerTemplate string
	OrderInfoCustomerSubject  string

	PaymentInfoSupportTemplate string
	PaymentInfoSupportSubject  string

	PaymentInfoCustomerTemplate string
	PaymentInfoCustomerSubject  string
}

type EmailService struct {
	client email.Sender
	EmailDeps
}

func NewEmailService(client email.Sender, deps EmailDeps) *EmailService {
	return &EmailService{client: client, EmailDeps: deps}
}

func (s *EmailService) SendOrderInfoSupport(inp jewerly.OrderInfoEmailInput) error {
	message := email.Email{
		ToName:    s.SupportName,
		ToEmail:   s.SupportEmail,
		FromEmail: s.SenderEmail,
		FromName:  s.SenderName,
		Subject:   fmt.Sprintf(s.OrderInfoSupportSubject, inp.OrderId, inp.TransactionStatus),
	}

	inp.OrderedAtFormated = inp.OrderedAt.Format(time.RFC822)

	if err := message.GenerateBodyFromHTML(s.OrderInfoSupportTemplate, inp); err != nil {
		return err
	}

	return s.client.Send(message)
}

func (s *EmailService) SendOrderInfoCustomer(inp jewerly.OrderInfoEmailInput) error {
	message := email.Email{
		ToName:    inp.FirstName,
		ToEmail:   inp.Email,
		FromEmail: s.SenderEmail,
		FromName:  s.SenderName,
		Subject:   fmt.Sprintf(s.OrderInfoCustomerSubject, inp.OrderId),
	}

	inp.OrderedAtFormated = inp.OrderedAt.Format(time.RFC822)

	if err := message.GenerateBodyFromHTML(s.OrderInfoCustomerTemplate, inp); err != nil {
		return err
	}

	return s.client.Send(message)
}

func (s *EmailService) SendPaymentInfoSupport(inp jewerly.PaymentInfoEmailInput) error {
	message := email.Email{
		ToName:    s.SupportName,
		ToEmail:   s.SupportEmail,
		FromEmail: s.SenderEmail,
		FromName:  s.SenderName,
		Subject:   fmt.Sprintf(s.PaymentInfoSupportSubject, inp.OrderId, inp.Status),
	}

	if err := message.GenerateBodyFromHTML(s.PaymentInfoSupportTemplate, inp); err != nil {
		return err
	}

	return s.client.Send(message)
}

func (s *EmailService) SendPaymentInfoCustomer(inp jewerly.PaymentInfoEmailInput) error {
	message := email.Email{
		ToName:    inp.BuyerName,
		ToEmail:   inp.BuyerEmail,
		FromEmail: s.SenderEmail,
		FromName:  s.SenderName,
		Subject:   fmt.Sprintf(s.PaymentInfoCustomerSubject, inp.OrderId, inp.Status),
	}

	if err := message.GenerateBodyFromHTML(s.PaymentInfoCustomerTemplate, inp); err != nil {
		return err
	}

	return s.client.Send(message)
}