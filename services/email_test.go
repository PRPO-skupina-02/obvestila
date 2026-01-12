package services

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadTemplates(t *testing.T) {
	// Set minimal env vars for initialization
	t.Setenv("RESEND_API_KEY", "test-key")
	t.Setenv("RESEND_FROM_EMAIL", "test@example.com")

	service, err := NewEmailService()
	require.NoError(t, err)
	require.NotNil(t, service)

	// Check that expected templates are loaded
	expectedTemplates := []string{
		"password_reset",
		"movie_suggestion",
		"welcome",
	}

	for _, name := range expectedTemplates {
		tmpl, ok := service.templates[name]
		assert.True(t, ok, "template %s should be loaded", name)
		assert.NotNil(t, tmpl, "template %s should not be nil", name)
	}
}

func TestNewEmailService_MissingAPIKey(t *testing.T) {
	// Explicitly unset RESEND_API_KEY by setting it to empty string
	t.Setenv("RESEND_API_KEY", "")
	t.Setenv("RESEND_FROM_EMAIL", "test@example.com")

	service, err := NewEmailService()
	assert.Error(t, err)
	assert.Nil(t, service)
	assert.Contains(t, err.Error(), "RESEND_API_KEY")
}

func TestNewEmailService_DefaultFromEmail(t *testing.T) {
	t.Setenv("RESEND_API_KEY", "test-key")
	// Explicitly unset RESEND_FROM_EMAIL to test default
	t.Setenv("RESEND_FROM_EMAIL", "")

	service, err := NewEmailService()
	require.NoError(t, err)
	assert.Equal(t, "noreply@yourdomain.com", service.fromEmail)
	assert.Equal(t, "CineCore", service.fromName)
}

func TestSendTemplatedEmail_TemplateNotFound(t *testing.T) {
	t.Setenv("RESEND_API_KEY", "test-key")
	t.Setenv("RESEND_FROM_EMAIL", "test@example.com")

	service, err := NewEmailService()
	require.NoError(t, err)

	// Try to send with non-existent template
	err = service.SendTemplatedEmail("test@example.com", "nonexistent", nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "template nonexistent not found")
}

func TestRenderTemplate_PasswordReset(t *testing.T) {
	t.Setenv("RESEND_API_KEY", "test-key")
	t.Setenv("RESEND_FROM_EMAIL", "test@example.com")

	service, err := NewEmailService()
	require.NoError(t, err)

	// Test template rendering without actually sending
	tmpl, ok := service.templates["password_reset"]
	require.True(t, ok, "password_reset template should exist")

	data := map[string]interface{}{
		"Subject":   "Reset Your Password",
		"UserName":  "John Doe",
		"ResetLink": "https://example.com/reset?token=abc123",
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	require.NoError(t, err)

	html := buf.String()
	assert.Contains(t, html, "John Doe")
	assert.Contains(t, html, "https://example.com/reset?token=abc123")
	assert.Contains(t, html, "CineCore")
}

func TestRenderTemplate_MovieSuggestion(t *testing.T) {
	t.Setenv("RESEND_API_KEY", "test-key")
	t.Setenv("RESEND_FROM_EMAIL", "test@example.com")

	service, err := NewEmailService()
	require.NoError(t, err)

	tmpl, ok := service.templates["movie_suggestion"]
	require.True(t, ok, "movie_suggestion template should exist")

	data := map[string]interface{}{
		"Subject":      "New Movie Available!",
		"UserName":     "Jane Smith",
		"MovieTitle":   "Inception",
		"MovieGenre":   "Sci-Fi",
		"MovieDuration": "148 minutes",
		"ShowingDate":  "January 15, 2026",
		"BookingLink":  "https://example.com/book/inception",
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	require.NoError(t, err)

	html := buf.String()
	assert.Contains(t, html, "Jane Smith")
	assert.Contains(t, html, "Inception")
	assert.Contains(t, html, "Sci-Fi")
	assert.Contains(t, html, "148 minutes")
	assert.Contains(t, html, "January 15, 2026")
	assert.Contains(t, html, "CineCore")
}

func TestRenderTemplate_Welcome(t *testing.T) {
	t.Setenv("RESEND_API_KEY", "test-key")
	t.Setenv("RESEND_FROM_EMAIL", "test@example.com")

	service, err := NewEmailService()
	require.NoError(t, err)

	tmpl, ok := service.templates["welcome"]
	require.True(t, ok, "welcome template should exist")

	data := map[string]interface{}{
		"Subject":  "Welcome to CineCore!",
		"UserName": "Alice Johnson",
		"AppLink":  "https://example.com/login",
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	require.NoError(t, err)

	html := buf.String()
	assert.Contains(t, html, "Alice Johnson")
	assert.Contains(t, html, "Welcome to CineCore!")
	assert.Contains(t, html, "https://example.com/login")
	assert.Contains(t, html, "CineCore Team")
}

func TestTemplateData_MissingSubject(t *testing.T) {
	t.Setenv("RESEND_API_KEY", "test-key")
	t.Setenv("RESEND_FROM_EMAIL", "test@example.com")

	service, err := NewEmailService()
	require.NoError(t, err)

	tmpl, ok := service.templates["welcome"]
	require.True(t, ok)

	// Render template without Subject in data
	data := map[string]interface{}{
		"UserName": "Test User",
		"AppLink":  "https://example.com",
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	require.NoError(t, err)

	// Template should still render successfully
	html := buf.String()
	assert.Contains(t, html, "Test User")
}
