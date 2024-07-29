package service

import (
	"context"
	"github.com/sirupsen/logrus"
	"mesaggio-test/config"
	"mesaggio-test/internal/messages"
	"mesaggio-test/internal/models"
	"mesaggio-test/pkg/kafka"
)

type messageService struct {
	cfg           *config.Config
	logger        *logrus.Logger
	messageRepo   messages.Repository
	kafkaProducer kafka.Producer
}

func NewMessagesService(log *logrus.Logger, cfg *config.Config, messageRepository messages.Repository, kafkaProducer kafka.Producer) messages.Service {
	return &messageService{
		cfg:           cfg,
		logger:        log,
		messageRepo:   messageRepository,
		kafkaProducer: kafkaProducer,
	}
}

func (m messageService) ReceiveMessage(ctx context.Context, message *models.Message) error {
	_, err := m.messageRepo.SaveMessage(ctx, message)
	if err != nil {
		return err
	}

	m.logger.Info("Successfully created")

	return nil
}

func (m messageService) UpdateMessage(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}
