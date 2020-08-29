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

	OrderInfoTemplate string
	OrderInfoSubject  string
}

type EmailService struct {
	client email.Sender
	EmailDeps
}

func NewEmailService(client email.Sender, deps EmailDeps) *EmailService {
	return &EmailService{client: client, EmailDeps: deps}
}

func (s *EmailService) SendOrderInfo(inp jewerly.OrderInfoEmailInput) error {
	message := email.Email{
		ToName:    s.SupportName,
		ToEmail:   s.SupportEmail,
		FromEmail: s.SenderEmail,
		FromName:  s.SenderName,
		Subject:   fmt.Sprintf(s.OrderInfoSubject, inp.OrderId, inp.TransactionStatus),
	}

	inp.OrderedAtFormated = inp.OrderedAt.Format(time.RFC3339)

	if err := message.GenerateBodyFromHTML(s.OrderInfoTemplate, inp); err != nil {
		return err
	}

	return s.client.Send(message)
}
