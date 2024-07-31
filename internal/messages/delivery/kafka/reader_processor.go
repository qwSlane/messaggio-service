package kafka

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"mesaggio-test/config"
	"mesaggio-test/internal/messages"
	"mesaggio-test/internal/models"
	"mesaggio-test/pkg/constants"
	"sync"
)

type readerMessageProcessor struct {
	log     *logrus.Logger
	cfg     *config.Config
	service messages.Service
}

func NewReaderMessageProcessor(log *logrus.Logger, cfg *config.Config, service messages.Service) *readerMessageProcessor {
	return &readerMessageProcessor{log: log, cfg: cfg, service: service}
}

func (s *readerMessageProcessor) ProcessMessages(ctx context.Context, r *kafka.Reader, wg *sync.WaitGroup, workerID int) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		m, err := r.FetchMessage(ctx)
		if err != nil {
			s.log.Warnf("workerID: %v, err: %v", workerID, err)
			continue
		}

		switch m.Topic {
		case s.cfg.KafkaTopics.MessageSaved.TopicName:
			s.processMessageSaved(ctx, r, m)
		}
	}
}

func (s *readerMessageProcessor) processMessageSaved(ctx context.Context, r *kafka.Reader, m kafka.Message) {
	msg := &models.Message{}
	err := json.Unmarshal(m.Value, msg)
	if err != nil {
		s.log.Warn("json.Unmarshal", err)
		return
	}

	message, err := s.service.UpdateMessage(ctx, msg)
	if err != nil {
		s.log.Warn("UpdateMessage", err)
		return
	}

	s.log.WithFields(logrus.Fields{
		constants.MessageID: message.MessageID.String(),
		constants.Topic:     m.Topic,
		constants.Partition: m.Partition,
		constants.Offset:    m.Offset,
	}).Info("Committed Kafka message")

	if err := r.CommitMessages(ctx, m); err != nil {
		s.log.Warn("commitMessage", err)
	}
}
