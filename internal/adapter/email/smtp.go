package email

import (
	"context"
	"fmt"
	"log/slog"
	"net/smtp"
	"strings"
)

type SMTPMailer struct {
	host, port, username, password, from string
	logger                               *slog.Logger
}

func NewSMTPMailer(host, port, username, password, from string, logger *slog.Logger) *SMTPMailer {
	if host == "" || port == "" || from == "" {
		return &SMTPMailer{logger: logger}
	}
	return &SMTPMailer{host: host, port: port, username: username, password: password, from: from, logger: logger}
}

func (m *SMTPMailer) SendOTP(_ context.Context, to, otp string) error {
	if m.host == "" {
		if m.logger != nil {
			m.logger.Info("email otp generated", "email", to, "otp", otp)
		}
		return nil
	}
	addr := m.host + ":" + m.port
	var auth smtp.Auth
	if m.username != "" {
		auth = smtp.PlainAuth("", m.username, m.password, m.host)
	}
	subject := "CodeLife Study email verification"
	body := fmt.Sprintf("Ma OTP xac thuc email cua ban la: %s\nOTP co hieu luc trong vai phut.", otp)
	message := strings.Join([]string{
		"From: " + m.from,
		"To: " + to,
		"Subject: " + subject,
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=UTF-8",
		"",
		body,
	}, "\r\n")
	return smtp.SendMail(addr, auth, m.from, []string{to}, []byte(message))
}
