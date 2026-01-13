package main

import (
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/PRPO-skupina-02/common/logging"
	"github.com/PRPO-skupina-02/obvestila/api"
	"github.com/PRPO-skupina-02/obvestila/queue"
	"github.com/PRPO-skupina-02/obvestila/services"
	"github.com/gin-gonic/gin"
)

func main() {
	err := run()

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func run() error {
	slog.Info("Starting obvestila service")

	logger := logging.GetDefaultLogger()
	slog.SetDefault(logger)

	// Initialize email service
	emailService, err := services.NewEmailService()
	if err != nil {
		return err
	}
	slog.Info("Email service initialized")

	// Get RabbitMQ URL from environment
	rabbitmqURL := os.Getenv("RABBITMQ_URL")
	if rabbitmqURL == "" {
		rabbitmqURL = "amqp://guest:guest@localhost:5672/"
	}

	// Initialize RabbitMQ consumer
	consumer, err := queue.NewEmailConsumer(rabbitmqURL, emailService)
	if err != nil {
		return err
	}
	defer consumer.Close()

	// Start consuming messages
	err = consumer.Start()
	if err != nil {
		return err
	}
	slog.Info("RabbitMQ consumer started")

	// Setup HTTP server
	router := gin.Default()

	// Add CORS middleware
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	api.Register(router)

	// Start HTTP server in a goroutine
	go func() {
		slog.Info("Starting HTTP server on :8080")
		if err := router.Run(":8080"); err != nil {
			slog.Error("HTTP server error", "error", err)
		}
	}()

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	slog.Info("Shutting down gracefully...")
	return nil
}
