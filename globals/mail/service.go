package mail

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/wneessen/go-mail"
)

// MailService represents the email service
type MailService struct {
	client *mail.Client
	config *MailConfig
}

// MailConfig holds email configuration
type MailConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	From     string `yaml:"from"`
	TLSMode  string `yaml:"tls_mode"` // "mandatory" (STARTTLS), "opportunistic" (STARTTLS with fallback), "ssl" (SSL), "none" (NoTLS)
}

// NewMailService creates a new mail service
func NewMailService(config *MailConfig) (*MailService, error) {
	// - "mandatory": Uses STARTTLS (requires TLS encryption)
	// - "opportunistic": Uses STARTTLS with fallback to plain text if TLS fails
	// - "ssl": Uses SSL (direct SSL connection, typically port 465)
	// - "none": Uses NoTLS (no encryption, plain text)
	var client *mail.Client
	var err error

	switch config.TLSMode {
	case "mandatory":
		client, err = mail.NewClient(
			config.Host,
			mail.WithPort(config.Port),
			mail.WithUsername(config.Username),
			mail.WithPassword(config.Password),
			mail.WithTLSPolicy(mail.TLSMandatory),
			mail.WithSMTPAuth(mail.SMTPAuthPlain),
		)
	case "opportunistic":
		client, err = mail.NewClient(
			config.Host,
			mail.WithPort(config.Port),
			mail.WithUsername(config.Username),
			mail.WithPassword(config.Password),
			mail.WithTLSPolicy(mail.TLSOpportunistic),
			mail.WithSMTPAuth(mail.SMTPAuthPlain),
		)
	case "ssl":
		client, err = mail.NewClient(
			config.Host,
			mail.WithPort(465),
			mail.WithUsername(config.Username),
			mail.WithPassword(config.Password),
			mail.WithSSL(),
			mail.WithSMTPAuth(mail.SMTPAuthPlain),
		)
	case "none":
		client, err = mail.NewClient(
			config.Host,
			mail.WithPort(config.Port),
			mail.WithUsername(config.Username),
			mail.WithPassword(config.Password),
			mail.WithTLSPolicy(mail.NoTLS),
			mail.WithSMTPAuth(mail.SMTPAuthPlain),
		)
	default:
		client, err = mail.NewClient(
			config.Host,
			mail.WithPort(config.Port),
			mail.WithUsername(config.Username),
			mail.WithPassword(config.Password),
			mail.WithTLSPolicy(mail.TLSMandatory),
			mail.WithSMTPAuth(mail.SMTPAuthPlain),
		)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create mail client: %w", err)
	}

	return &MailService{
		client: client,
		config: config,
	}, nil
}

func (s *MailService) SendVerificationEmail(to, token, verificationURL string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	msg := mail.NewMsg()
	if err := msg.From(s.config.From); err != nil {
		return fmt.Errorf("failed to set from address: %w", err)
	}
	if err := msg.To(to); err != nil {
		return fmt.Errorf("failed to set to address: %w", err)
	}
	msg.Subject("Circle Center Email Verification")

	htmlContent, err := LoadVerificationTemplate(to, token, verificationURL)
	if err != nil {
		return fmt.Errorf("failed to load verification template: %w", err)
	}

	msg.SetBodyString(mail.TypeTextHTML, htmlContent)

	if err := s.client.DialAndSendWithContext(ctx, msg); err != nil {
		slog.Error("Failed to send verification email", "to", to, "error", err)
		return fmt.Errorf("failed to send email: %w", err)
	}

	slog.Info("Verification email sent successfully", "to", to)
	return nil
}

func (s *MailService) SendTextEmail(to, subject, body string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	msg := mail.NewMsg()
	if err := msg.From(s.config.From); err != nil {
		return fmt.Errorf("failed to set from address: %w", err)
	}
	if err := msg.To(to); err != nil {
		return fmt.Errorf("failed to set to address: %w", err)
	}
	msg.Subject(subject)

	msg.SetBodyString(mail.TypeTextPlain, body)

	if err := s.client.DialAndSendWithContext(ctx, msg); err != nil {
		slog.Error("Failed to send text email", "to", to, "error", err)
		return fmt.Errorf("failed to send email: %w", err)
	}

	slog.Info("Text email sent successfully", "to", to)
	return nil
}

func (s *MailService) Close() error {
	if s.client != nil {
		return s.client.Close()
	}
	return nil
}
