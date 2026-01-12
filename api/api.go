package api

import (
	"net/http"

	_ "github.com/PRPO-skupina-02/obvestila/api/docs"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//	@title			Obvestila API
//	@version		1.0
//	@description	Email notification service for the PRPO project. Messages are consumed from RabbitMQ queue.

//	@host		localhost:8080
//	@BasePath	/api/v1/obvestila

// Register sets up the API routes
func Register(router *gin.Engine) {
	// Healthcheck
	router.GET("/healthcheck", healthcheck)

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API endpoints (for documentation and potential future direct endpoints)
	v1 := router.Group("/api/v1/obvestila")

	// Info endpoint - describes how to use the service
	v1.GET("/info", getInfo)
}

// healthcheck godoc
//
//	@Summary		Health check
//	@Description	Returns OK if the service is healthy
//	@Tags			system
//	@Produce		plain
//	@Success		200	{string}	string	"OK"
//	@Router			/healthcheck [get]
func healthcheck(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}

// EmailMessageSchema represents the schema for RabbitMQ email messages
type EmailMessageSchema struct {
	To           string                 `json:"to" binding:"required,email" example:"user@example.com"`
	Template     string                 `json:"template" binding:"required" example:"password_reset"`
	TemplateData map[string]interface{} `json:"template_data,omitempty"`
} // @name EmailMessage

// ServiceInfo represents information about the email service
type ServiceInfo struct {
	ServiceName      string   `json:"service_name" example:"obvestila"`
	Description      string   `json:"description" example:"Email notification service using RabbitMQ and Resend"`
	QueueName        string   `json:"queue_name" example:"emails"`
	AvailableTemplates []string `json:"available_templates"`
	Usage            string   `json:"usage" example:"Publish JSON messages to the 'emails' queue"`
} // @name ServiceInfo

// getInfo godoc
//
//	@Summary		Get service information
//	@Description	Returns information about the email service, available templates, and usage instructions
//	@Tags			info
//	@Produce		json
//	@Success		200	{object}	ServiceInfo
//	@Router			/api/v1/obvestila/info [get]
func getInfo(c *gin.Context) {
	info := ServiceInfo{
		ServiceName:      "obvestila",
		Description:      "Email notification service using RabbitMQ and Resend. All emails are sent from CineCore with template-based content.",
		QueueName:        "emails",
		AvailableTemplates: []string{
			"password_reset",
			"movie_suggestion",
			"welcome",
		},
		Usage: "Publish EmailMessage JSON to the 'emails' RabbitMQ queue. See schema in Swagger documentation.",
	}

	c.JSON(http.StatusOK, info)
}
