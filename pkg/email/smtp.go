package email

import "net/smtp"

type SMTPClient struct {
	hostname  string
	port      string
	emailFrom string
	auth      smtp.Auth
}

func NewSMTPClient(hostname, port, sender, pass string) *SMTPClient {
	return &SMTPClient{
		hostname:  hostname,
		port:      port,
		emailFrom: sender,
		auth:      smtp.PlainAuth("", sender, pass, hostname),
	}
}

func (c *SMTPClient) Send(m Email) error {
	message, err := m.EmailBytes()
	if err != nil {
		return err
	}

	return smtp.SendMail(c.hostname+c.port, c.auth, c.emailFrom, []string{m.ToEmail}, message)
}
