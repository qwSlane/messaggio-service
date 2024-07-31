package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

type Producer interface {
	PublishMessage(ctx context.Context, msgs ...kafka.Message) error
	Close() error
}

type producer struct {
	log     *logrus.Logger
	brokers []string
	w       *kafka.Writer
}

// NewProducer create new kafka producer
func NewProducer(log *logrus.Logger, brokers []string) Producer {
	return &producer{log: log, brokers: brokers, w: NewWriter(brokers, kafka.LoggerFunc(log.Errorf))}
}

func (p *producer) PublishMessage(ctx context.Context, msgs ...kafka.Message) error {
	return p.w.WriteMessages(ctx, msgs...)
}

func (p *producer) Close() error {
	return p.w.Close()
}
