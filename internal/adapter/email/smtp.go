package email

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"mime"
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
	subject := mime.BEncoding.Encode("utf-8", "Xác thực Email - CodeLife Academy")
	body := fmt.Sprintf(otpEmailTemplate, otp)
	return []byte(strings.Join([]string{
		"From: " + m.fromHeader,
		"To: " + to,
		"Subject: " + subject,
		"MIME-Version: 1.0",
		"Content-Type: text/html; charset=UTF-8",
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

const otpEmailTemplate = `<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Xác thực Email - CodeLife Academy</title>
</head>
<body style="margin: 0; padding: 0; background-color: #f6f9fc; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif; -webkit-font-smoothing: antialiased; -moz-osx-font-smoothing: grayscale;">
    <table border="0" cellpadding="0" cellspacing="0" width="100%%" style="background-color: #f6f9fc; padding: 40px 0;">
        <tr>
            <td align="center">
                <table border="0" cellpadding="0" cellspacing="0" width="100%%" style="max-width: 500px; background-color: #ffffff; border-radius: 12px; overflow: hidden; box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05); border: 1px solid #eef2f5;">
                    <!-- Header with Gradient -->
                    <tr>
                        <td align="center" style="background: linear-gradient(135deg, #6366f1 0%%, #a855f7 100%%); padding: 32px 24px;">
                            <h1 style="color: #ffffff; margin: 0; font-size: 24px; font-weight: 700; letter-spacing: -0.5px;">CodeLife Academy</h1>
                        </td>
                    </tr>
                    <!-- Body Content -->
                    <tr>
                        <td style="padding: 40px 32px; color: #334155; font-size: 15px; line-height: 1.6;">
                            <p style="margin: 0 0 16px 0; font-weight: 600; font-size: 18px; color: #0f172a;">Xác thực tài khoản của bạn</p>
                            <p style="margin: 0 0 24px 0; color: #64748b;">Cảm ơn bạn đã đồng hành cùng CodeLife Academy. Vui lòng sử dụng mã OTP dưới đây để hoàn tất xác thực địa chỉ email:</p>

                            <!-- OTP Box -->
                            <table border="0" cellpadding="0" cellspacing="0" width="100%%" style="margin-bottom: 24px;">
                                <tr>
                                    <td align="center" style="background-color: #f8fafc; border-radius: 8px; border: 1px dashed #cbd5e1; padding: 20px 24px;">
                                        <div style="font-size: 32px; font-weight: 700; letter-spacing: 6px; color: #6366f1; font-family: 'Courier New', Courier, monospace; margin: 0;">%s</div>
                                    </td>
                                </tr>
                            </table>

                            <p style="margin: 0 0 24px 0; font-size: 13px; color: #94a3b8; line-height: 1.5;">
                                * Mã OTP này có hiệu lực trong vòng <strong>5 phút</strong>. Vì lý do bảo mật, vui lòng không chia sẻ mã này với bất kỳ ai khác.
                            </p>

                            <hr style="border: 0; border-top: 1px solid #f1f5f9; margin: 24px 0;" />

                            <p style="margin: 0; font-size: 12px; color: #94a3b8; text-align: center;">
                                Nếu bạn không yêu cầu thực hiện hành động này, bạn có thể an tâm bỏ qua email.
                            </p>
                        </td>
                    </tr>
                    <!-- Footer -->
                    <tr>
                        <td style="background-color: #f8fafc; padding: 20px 32px; text-align: center; border-top: 1px solid #f1f5f9;">
                            <p style="margin: 0; font-size: 12px; color: #94a3b8;">
                                &copy; 2026 CodeLife Academy. All rights reserved.
                            </p>
                        </td>
                    </tr>
                </table>
            </td>
        </tr>
    </table>
</body>
</html>`
