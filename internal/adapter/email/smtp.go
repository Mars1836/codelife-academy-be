package email

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/mail"
	"net/smtp"
	"strings"
	"time"
)

type SMTPMailer struct {
	host, port, username, password string
	fromAddress, fromHeader        string
	secure                         bool
}

func NewSMTPMailer(host, port, username, password, from string, secure bool) (*SMTPMailer, error) {
	if host == "" {
		return nil, fmt.Errorf("MAIL_HOST is required")
	}
	if port == "" {
		return nil, fmt.Errorf("MAIL_PORT is required when MAIL_HOST is configured")
	}
	parsedFrom, err := mail.ParseAddress(from)
	if err != nil || parsedFrom.Address == "" {
		return nil, fmt.Errorf("MAIL_FROM must contain a valid email address")
	}
	return &SMTPMailer{
		host:        host,
		port:        port,
		username:    username,
		password:    password,
		fromAddress: parsedFrom.Address,
		fromHeader:  parsedFrom.String(),
		secure:      secure,
	}, nil
}

func (m *SMTPMailer) SendOTP(ctx context.Context, to, otp string) error {
	parsedTo, err := mail.ParseAddress(to)
	if err != nil || parsedTo.Address == "" {
		return fmt.Errorf("recipient email is invalid")
	}

	message := m.otpMessage(parsedTo.Address, otp)
	if m.secure {
		return m.sendImplicitTLS(ctx, parsedTo.Address, message)
	}

	var auth smtp.Auth
	if m.username != "" {
		auth = smtp.PlainAuth("", m.username, m.password, m.host)
	}
	return smtp.SendMail(net.JoinHostPort(m.host, m.port), auth, m.fromAddress, []string{parsedTo.Address}, message)
}

func (m *SMTPMailer) otpMessage(to, otp string) []byte {
	subject := "CodeLife Study email verification"
	body := fmt.Sprintf("Ma OTP xac thuc email cua ban la: %s\nOTP co hieu luc trong vai phut.", otp)
	return []byte(strings.Join([]string{
		"From: " + m.fromHeader,
		"To: " + to,
		"Subject: " + subject,
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=UTF-8",
		"",
		body,
	}, "\r\n"))
}

func (m *SMTPMailer) sendImplicitTLS(ctx context.Context, to string, message []byte) error {
	dialer := &net.Dialer{Timeout: 10 * time.Second}
	connection, err := dialer.DialContext(ctx, "tcp", net.JoinHostPort(m.host, m.port))
	if err != nil {
		return err
	}
	tlsConnection := tls.Client(connection, &tls.Config{ServerName: m.host, MinVersion: tls.VersionTLS12})
	if err := tlsConnection.HandshakeContext(ctx); err != nil {
		_ = connection.Close()
		return err
	}
	client, err := smtp.NewClient(tlsConnection, m.host)
	if err != nil {
		_ = tlsConnection.Close()
		return err
	}
	defer client.Close()

	if m.username != "" {
		if err := client.Auth(smtp.PlainAuth("", m.username, m.password, m.host)); err != nil {
			return err
		}
	}
	if err := client.Mail(m.fromAddress); err != nil {
		return err
	}
	if err := client.Rcpt(to); err != nil {
		return err
	}
	wc, err := client.Data()
	if err != nil {
		return err
	}
	if _, err := wc.Write(message); err != nil {
		_ = wc.Close()
		return err
	}
	if err := wc.Close(); err != nil {
		return err
	}
	if err := client.Quit(); err != nil && err != io.EOF {
		return err
	}
	return nil
}
