package smtp

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strings"
	"time"

	"github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/entities"
)

type SMTPSender struct {
	host     string
	port     int
	username string
	password string
	from     string
}

func NewSMTPSender(host string, port int, username, password, from string) *SMTPSender {
	return &SMTPSender{
		host:     host,
		port:     port,
		username: username,
		password: password,
		from:     from,
	}
}

func (s *SMTPSender) Send(msg entities.Message) error {
	addr := fmt.Sprintf("%s:%d", s.host, s.port)

	auth := smtp.PlainAuth("", s.username, s.password, s.host)

	// Prepare email headers and body
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("From: %s\r\n", s.from))
	sb.WriteString(fmt.Sprintf("To: %s\r\n", msg.To))
	sb.WriteString(fmt.Sprintf("Subject: %s\r\n", msg.Subject))
	sb.WriteString(fmt.Sprintf("Date: %s\r\n", time.Now().UTC().Format(time.RFC1123Z)))
	sb.WriteString("MIME-Version: 1.0\r\n")
	sb.WriteString("Content-Type: text/html; charset=\"UTF-8\"\r\n")
	sb.WriteString("\r\n") // empty line before body
	sb.WriteString(msg.Body)

	// TLS config
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         s.host,
	}

	// Connect to SMTP server
	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return fmt.Errorf("smtp dial tls failed: %w", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, s.host)
	if err != nil {
		return fmt.Errorf("smtp new client failed: %w", err)
	}
	defer client.Quit()

	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("smtp auth failed: %w", err)
	}

	if err = client.Mail(s.from); err != nil {
		return fmt.Errorf("smtp mail failed: %w", err)
	}

	if err = client.Rcpt(msg.To); err != nil {
		return fmt.Errorf("smtp rcpt failed: %w", err)
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("smtp data open failed: %w", err)
	}

	_, err = w.Write([]byte(sb.String()))
	if err != nil {
		return fmt.Errorf("smtp write failed: %w", err)
	}

	err = w.Close()
	if err != nil {
		return fmt.Errorf("smtp close failed: %w", err)
	}

	return nil
}
