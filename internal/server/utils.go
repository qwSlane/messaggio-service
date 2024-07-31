package server

import (
	"context"
	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
	"net"
	"strconv"
)

func (s *Server) initKafkaTopics(ctx context.Context) error {
	kafkaConn, err := kafka.DialContext(ctx, "tcp", s.cfg.Kafka.Brokers[0])
	if err != nil {
		return errors.Wrap(err, "kafka.NewKafkaCon")
	}

	brokers, err := kafkaConn.Brokers()
	if err != nil {
		return errors.Wrap(err, "kafkaConn.Brokers")
	}
	s.log.Infof("kafka connected to brokers: %+v", brokers)

	controller, err := kafkaConn.Controller()
	if err != nil {
		s.log.Warn("kafkaConn.Controller", err)
		return err
	}

	controllerURI := net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port))
	s.log.Infof("kafka controller uri: %s", controllerURI)

	conn, err := kafka.DialContext(ctx, "tcp", controllerURI)
	if err != nil {
		s.log.Warn("initKafkaTopics.DialContext", err)
		return err
	}
	defer conn.Close()

	messageSavedTopic := kafka.TopicConfig{
		Topic:             s.cfg.KafkaTopics.MessageSaved.TopicName,
		NumPartitions:     s.cfg.KafkaTopics.MessageSaved.Partitions,
		ReplicationFactor: s.cfg.KafkaTopics.MessageSaved.ReplicationFactor,
	}

	if err := conn.CreateTopics(messageSavedTopic); err != nil {
		s.log.Warn("kafkaConn.CreateTopics", err)
		return err
	}

	s.log.Infof("kafka topics created or already exists: %+v", messageSavedTopic)

	return nil
}

func (s *Server) getConsumerGroupTopics() []string {
	return []string{
		s.cfg.KafkaTopics.MessageSaved.TopicName,
	}
}
