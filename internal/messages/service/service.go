package service

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"mesaggio-test/config"
	"mesaggio-test/internal/messages"
	"mesaggio-test/internal/models"
	kafkaClient "mesaggio-test/pkg/kafka"
	"time"
)

type messageService struct {
	cfg           *config.Config
	logger        *logrus.Logger
	messageRepo   messages.Repository
	kafkaProducer kafkaClient.Producer
}

func NewMessagesService(log *logrus.Logger, cfg *config.Config, messageRepository messages.Repository, kafkaProducer kafkaClient.Producer) messages.Service {
	return &messageService{
		cfg:           cfg,
		logger:        log,
		messageRepo:   messageRepository,
		kafkaProducer: kafkaProducer,
	}
}

func (m messageService) ReceiveMessage(ctx context.Context, message *models.Message) error {
	savedMessage, err := m.messageRepo.SaveMessage(ctx, message)
	if err != nil {
		return err
	}
	m.logger.Info("Successfully created")

	msgBytes, err := json.Marshal(savedMessage)
	if err != nil {
		return err
	}

	kafkaMessage := kafka.Message{
		Topic: m.cfg.KafkaTopics.MessageSaved.TopicName,
		Value: msgBytes,
		Time:  time.Now().UTC(),
	}

	return m.kafkaProducer.PublishMessage(ctx, kafkaMessage)
}

func (m messageService) UpdateMessage(ctx context.Context, message *models.Message) (*models.Message, error) {
	return m.messageRepo.UpdateMessage(ctx, message)
}

func (m messageService) GetMessagesStatistics(ctx context.Context) (*models.Statistics, error) {
	return m.messageRepo.GetMessagesStatistics(ctx)
}
