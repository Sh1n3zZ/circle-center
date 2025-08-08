package mail

import (
	"context"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/wneessen/go-mail"
)

func TestMailService_SendVerificationEmail(t *testing.T) {
	host := os.Getenv("TEST_SMTP_HOST")
	port := os.Getenv("TEST_SMTP_PORT")
	username := os.Getenv("TEST_SMTP_USERNAME")
	password := os.Getenv("TEST_SMTP_PASSWORD")
	from := os.Getenv("TEST_SMTP_FROM")

	if host == "" || port == "" || username == "" || password == "" || from == "" {
		t.Skip("SMTP configuration not available, skipping test")
	}

	portInt, err := strconv.Atoi(port)
	if err != nil {
		t.Fatalf("Failed to convert port to int: %v", err)
	}

	// - "mandatory": Uses STARTTLS (requires TLS encryption)
	// - "opportunistic": Uses STARTTLS with fallback to plain text if TLS fails
	// - "ssl": Uses SSL (direct SSL connection, typically port 465)
	// - "none": Uses NoTLS (no encryption, plain text)
	tlsMode := os.Getenv("TEST_SMTP_TLS_MODE")
	if tlsMode == "" {
		tlsMode = "opportunistic"
	}

	config := &MailConfig{
		Host:     host,
		Port:     portInt,
		Username: username,
		Password: password,
		From:     from,
		TLSMode:  tlsMode,
	}

	service, err := NewMailService(config)
	if err != nil {
		t.Fatalf("Failed to create mail service: %v", err)
	}
	defer service.Close()

	testEmail := os.Getenv("TEST_SMTP_TO")
	if testEmail == "" {
		t.Skip("TEST_SMTP_TO environment variable not set, skipping test")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	msg := mail.NewMsg()
	if err := msg.From(from); err != nil {
		t.Fatalf("Failed to set from address: %v", err)
	}
	if err := msg.To(testEmail); err != nil {
		t.Fatalf("Failed to set to address: %v", err)
	}
	msg.Subject("Circle Center Test Email")

	htmlContent := `
		<!DOCTYPE html>
		<html>
		<head>
			<meta charset="UTF-8">
			<title>Circle Center Test Email</title>
		</head>
		<body>
			<h1>Circle Center Test Email</h1>
			<p>This is a test email to verify that the email service is working correctly.</p>
			<p>Sent at: ` + time.Now().Format("2006-01-02 15:04:05") + `</p>
			<p>If you receive this email, it means the email service is configured correctly.</p>
		</body>
		</html>
	`

	msg.SetBodyString(mail.TypeTextHTML, htmlContent)

	if err := service.client.DialAndSendWithContext(ctx, msg); err != nil {
		t.Fatalf("Failed to send test email: %v", err)
	}

	t.Log("Test email sent successfully")
}

func TestMailService_SendTextEmail(t *testing.T) {
	host := os.Getenv("TEST_SMTP_HOST")
	port := os.Getenv("TEST_SMTP_PORT")
	username := os.Getenv("TEST_SMTP_USERNAME")
	password := os.Getenv("TEST_SMTP_PASSWORD")
	from := os.Getenv("TEST_SMTP_FROM")

	if host == "" || port == "" || username == "" || password == "" || from == "" {
		t.Skip("SMTP configuration not available, skipping test")
	}

	portInt, err := strconv.Atoi(port)
	if err != nil {
		t.Fatalf("Failed to convert port to int: %v", err)
	}

	tlsMode := os.Getenv("TEST_SMTP_TLS_MODE")
	if tlsMode == "" {
		tlsMode = "opportunistic"
	}

	config := &MailConfig{
		Host:     host,
		Port:     portInt,
		Username: username,
		Password: password,
		From:     from,
		TLSMode:  tlsMode,
	}

	service, err := NewMailService(config)
	if err != nil {
		t.Fatalf("Failed to create mail service: %v", err)
	}
	defer service.Close()

	testEmail := os.Getenv("TEST_SMTP_TO")
	if testEmail == "" {
		t.Skip("TEST_SMTP_TO environment variable not set, skipping test")
	}

	subject := "Test Plain Text Email"
	body := "This is a test plain text email from Circle Center.\n\nSent at: " + time.Now().Format("2006-01-02 15:04:05") + "\n\nBest regards,\nCircle Center Team"

	err = service.SendTextEmail(testEmail, subject, body)
	if err != nil {
		t.Fatalf("Failed to send text email: %v", err)
	}

	t.Log("Text email sent successfully")
}
