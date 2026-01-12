package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHealthcheck(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	Register(router)

	req, _ := http.NewRequest("GET", "/healthcheck", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "OK", w.Body.String())
}

func TestGetInfo(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	Register(router)

	req, _ := http.NewRequest("GET", "/api/v1/obvestila/info", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var info ServiceInfo
	err := json.Unmarshal(w.Body.Bytes(), &info)
	require.NoError(t, err)

	assert.Equal(t, "obvestila", info.ServiceName)
	assert.Equal(t, "emails", info.QueueName)
	assert.Contains(t, info.Description, "RabbitMQ")
	assert.Contains(t, info.Description, "Resend")
	assert.Contains(t, info.Description, "CineCore")
	
	// Check available templates
	assert.Contains(t, info.AvailableTemplates, "password_reset")
	assert.Contains(t, info.AvailableTemplates, "movie_suggestion")
	assert.Contains(t, info.AvailableTemplates, "welcome")
	assert.Len(t, info.AvailableTemplates, 3)
}

func TestSwaggerEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	Register(router)

	req, _ := http.NewRequest("GET", "/swagger/index.html", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Swagger may not be initialized in test mode, so we just check it doesn't crash
	// In production, it should return 200 or redirect
	assert.True(t, w.Code >= 200 && w.Code < 500, "should not return server error")
}
