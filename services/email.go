package services

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"log/slog"
	"os"

	"github.com/resend/resend-go/v2"
)

//go:embed templates/*.html
var templatesFS embed.FS

type EmailService struct {
	client    *resend.Client
	fromEmail string
	fromName  string
	templates map[string]*template.Template
}

func NewEmailService() (*EmailService, error) {
	apiKey := os.Getenv("RESEND_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("RESEND_API_KEY environment variable is required")
	}

	fromEmail := os.Getenv("RESEND_FROM_EMAIL")
	if fromEmail == "" {
		fromEmail = "noreply@yourdomain.com"
	}

	client := resend.NewClient(apiKey)

	service := &EmailService{
		client:    client,
		fromEmail: fromEmail,
		fromName:  "CineCore",
		templates: make(map[string]*template.Template),
	}

	if err := service.loadTemplates(); err != nil {
		return nil, fmt.Errorf("failed to load templates: %w", err)
	}

	slog.Info("Email service initialized",
		"from_email", fromEmail,
		"from_name", "CineCore",
	)

	return service, nil
}

func (s *EmailService) loadTemplates() error {
	templates := []string{
		"password_reset",
		"movie_suggestion",
		"welcome",
		"recommendation",
	}

	for _, name := range templates {
		path := fmt.Sprintf("templates/%s.html", name)
		content, err := templatesFS.ReadFile(path)
		if err != nil {
			slog.Warn("Template not found, skipping", "template", name, "error", err)
			continue
		}

		tmpl, err := template.New(name).Parse(string(content))
		if err != nil {
			return fmt.Errorf("failed to parse template %s: %w", name, err)
		}

		s.templates[name] = tmpl
		slog.Info("Loaded email template", "template", name)
	}

	return nil
}

func (s *EmailService) SendEmail(to, subject, htmlBody string) error {
	params := &resend.SendEmailRequest{
		From:    fmt.Sprintf("%s <%s>", s.fromName, s.fromEmail),
		To:      []string{to},
		Subject: subject,
		Html:    htmlBody,
	}

	sent, err := s.client.Emails.Send(params)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	slog.Info("Email sent successfully",
		"to", to,
		"subject", subject,
		"message_id", sent.Id,
	)

	return nil
}

func (s *EmailService) SendTemplatedEmail(to, templateName string, data map[string]interface{}) error {
	tmpl, ok := s.templates[templateName]
	if !ok {
		return fmt.Errorf("template %s not found", templateName)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}

	subject, ok := data["Subject"].(string)
	if !ok {
		subject = fmt.Sprintf("Email from %s", s.fromName)
	}

	return s.SendEmail(to, subject, buf.String())
}
