package email

import (
	"bytes"
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

func (m *Email) EmailBytes() []byte {
	to := mail.Address{m.ToName, m.ToEmail}
	from := mail.Address{m.FromName, m.FromEmail}

	header := make(map[string]string)
	header["To"] = to.String()
	header["From"] = from.String()
	header["Subject"] = m.Subject
	header["Content-Type"] = `text/html; charset="UTF-8"`
	msg := ""

	for k, v := range header {
		msg += fmt.Sprintf("%s: %s\r\n", k, v)
	}

	msg += "\r\n" + m.Body

	return []byte(msg)
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
