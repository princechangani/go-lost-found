package services

import (
	"gopkg.in/gomail.v2"
)

type EmailService struct {
	Host     string
	Port     int
	Username string
	Password string
}

func NewEmailService(host string, port int, username, password string) *EmailService {
	return &EmailService{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
	}
}

func (e *EmailService) SendEmail(subject, toEmail, body string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", e.Username)
	msg.SetHeader("To", toEmail)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", body)

	dialer := gomail.NewDialer(e.Host, e.Port, e.Username, e.Password)
	return dialer.DialAndSend(msg)
}
