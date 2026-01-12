package queue

import (
	"log/slog"

	"github.com/PRPO-skupina-02/common/messaging"
	"github.com/PRPO-skupina-02/obvestila/services"
)

type EmailConsumer struct {
	consumer     *messaging.Consumer
	emailService *services.EmailService
}

func NewEmailConsumer(rabbitmqURL string, emailService *services.EmailService) (*EmailConsumer, error) {
	ec := &EmailConsumer{
		emailService: emailService,
	}

	consumer, err := messaging.NewConsumer(rabbitmqURL, ec.handleEmailMessage)
	if err != nil {
		return nil, err
	}

	ec.consumer = consumer
	return ec, nil
}

func (ec *EmailConsumer) Start() error {
	return ec.consumer.Start()
}

func (ec *EmailConsumer) Close() error {
	if ec.consumer != nil {
		return ec.consumer.Close()
	}
	return nil
}

func (ec *EmailConsumer) WaitForever() {
	ec.consumer.WaitForever()
}

func (ec *EmailConsumer) handleEmailMessage(msg *messaging.EmailMessage) error {
	slog.Info("Processing email message",
		"to", msg.To,
		"template", msg.Template,
	)

	return ec.emailService.SendTemplatedEmail(msg.To, msg.Template, msg.TemplateData)
}
