package server

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"mesaggio-test/config"
	httpHandler "mesaggio-test/internal/messages/delivery/http"
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

	s.MapHandlers()
	//kafka producer

	//kafka consumer

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			s.log.Errorf("http Server listen err: %v", err)
		}
	}()
	s.log.Infof("http Server listen port: %v", s.cfg.HttpPort)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	ctx, shutdown := context.WithTimeout(context.Background(), ctxTimeout*time.Second)
	defer shutdown()

	s.log.Info("Server Exited Properly")
	return srv.Shutdown(ctx)
}

func (s *Server) MapHandlers() {
	messageRepository := messagesRepository.NewMessagesRepository(s.db)
	kafkaProducer := kafka.NewProducer(s.log, s.cfg.Kafka.Brokers)

	msgService := service.NewMessagesService(s.log, s.cfg, messageRepository, kafkaProducer)

	h := httpHandler.NewMessageHandler(s.log, msgService)

	s.router.Handle("POST /msg", h.ReceiveMessage())
}
