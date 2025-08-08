package mail

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
)

// VerificationTemplateData holds data for verification email template
type VerificationTemplateData struct {
	Email           string
	Token           string
	VerificationURL string
}

// LoadVerificationTemplate loads and renders the verification email template
func LoadVerificationTemplate(email, token, verificationURL string) (string, error) {
	templatePath := filepath.Join("globals", "mail", "templates", "verification", "verification.html")
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to read template file: %w", err)
	}

	tmpl, err := template.New("verification").Parse(string(templateContent))
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	data := VerificationTemplateData{
		Email:           email,
		Token:           token,
		VerificationURL: verificationURL,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}
