package server

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
	"mesaggio-test/config"
	_ "mesaggio-test/docs"
	"mesaggio-test/internal/messages"
	httpHandler "mesaggio-test/internal/messages/delivery/http"
	kafkaHandler "mesaggio-test/internal/messages/delivery/kafka"
	messagesRepository "mesaggio-test/internal/messages/repository"
	"mesaggio-test/internal/messages/service"
	"mesaggio-test/pkg/kafka"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	ctxTimeout = 5
)

type Server struct {
	router *http.ServeMux
	log    *logrus.Logger
	cfg    *config.Config
	db     *sqlx.DB
}

func NewServer(cfg *config.Config, log *logrus.Logger, db *sqlx.DB) *Server {
	return &Server{router: http.NewServeMux(), log: log, cfg: cfg, db: db}
}

func (s *Server) Run() error {
	srv := &http.Server{
		Addr:    s.cfg.HttpPort,
		Handler: s.router,
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	err := s.initKafkaTopics(ctx)
	if err != nil {
		return errors.Wrap(err, "s.connectKafkaBrokers")
	}

	messageRepository := messagesRepository.NewMessagesRepository(s.db)
	kafkaProducer := kafka.NewProducer(s.log, s.cfg.Kafka.Brokers)
	msgService := service.NewMessagesService(s.log, s.cfg, messageRepository, kafkaProducer)

	h := httpHandler.NewMessageHandler(s.log, msgService)
	s.MapHandlers(h)

	readerMessageProcessor := kafkaHandler.NewReaderMessageProcessor(s.log, s.cfg, msgService)
	cg := kafka.NewConsumerGroup(s.cfg.Kafka.Brokers, s.cfg.Kafka.GroupID, s.log)
	go cg.ConsumeTopic(ctx, s.getConsumerGroupTopics(), kafka.PoolSize, readerMessageProcessor.ProcessMessages)

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			s.log.Errorf("http Server listen err: %v", err)
		}
	}()
	s.log.Infof("http Server listen port: %v", s.cfg.HttpPort)

	<-ctx.Done()
	ctxShutdown, shutdown := context.WithTimeout(context.Background(), ctxTimeout*time.Second)
	defer shutdown()

	s.log.Info("server exited properly")
	return srv.Shutdown(ctxShutdown)
}

func (s *Server) MapHandlers(h messages.Handler) {
	s.router.Handle("POST /msg", h.ReceiveMessage())
	s.router.Handle("GET /stats", h.GetStatistics())
	s.router.HandleFunc("GET /swagger/*", httpSwagger.Handler())
}
