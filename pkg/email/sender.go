package email

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"html/template"
	"net/mail"
)

type Email struct {
	ToName    string
	ToEmail   string
	FromEmail string
	FromName  string

	Subject string
	Body    string
}

type Sender interface {
	Send(m Email) error
}

func (m *Email) EmailBytes() ([]byte, error) {
	if err := m.validate(); err != nil {
		return nil, err
	}

	var msg string

	headers := m.generateHeaders()
	for k, v := range headers {
		msg += fmt.Sprintf("%s: %s\r\n", k, v)
	}

	msg += "\r\n" + m.Body

	return []byte(msg), nil
}

func (m *Email) GenerateBodyFromHTML(templateFileName string, data interface{}) error {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		logrus.Errorf("failed to parse file %s:%s\n", templateFileName, err.Error())
		return err
	}

	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return err
	}

	m.Body = buf.String()

	return nil
}

func (m *Email) validate() error {
	if m.ToName == "" || m.ToEmail == "" || m.FromEmail == "" || m.FromName == "" {
		return errors.New("empty from/to")
	}

	if m.Subject == "" || m.Body == "" {
		return errors.New("empty subject/body")
	}

	if !isEmailValid(m.ToEmail) {
		return errors.New("invalid to email")
	}

	if !isEmailValid(m.FromEmail) {
		return errors.New("invalid from email")
	}

	return nil
}

func (m *Email) generateHeaders() map[string]string {
	to := mail.Address{m.ToName, m.ToEmail}
	from := mail.Address{m.FromName, m.FromEmail}

	header := make(map[string]string)

	header["To"] = to.String()
	header["From"] = from.String()
	header["Subject"] = m.Subject
	header["Content-Type"] = `text/html; charset="UTF-8"`

	return header
}
