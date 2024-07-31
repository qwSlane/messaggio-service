package main

import (
	"github.com/sirupsen/logrus"
	"log"
	"mesaggio-test/config"
	"mesaggio-test/internal/server"
	"mesaggio-test/pkg/postgres"
)

// @title Messaggio REST API
// @version 1.0
// @description Test task for messaggio
// @contact.name Siarhei Vasileuski
// @contact.url https://t.me/kataomione
// @contact.email sergej.vasilewsckij@yandex.ru

func main() {
	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatal(err)
	}

	appLogger := logrus.New()

	psqlDB, err := postgres.NewPsqlDB(cfg)
	if err != nil {
		appLogger.Fatalf("postgresql init: %s", err)
	} else {
		appLogger.Infof("postgres connected, Status: %#v", psqlDB.Stats())
	}
	defer psqlDB.Close()

	s := server.NewServer(cfg, appLogger, psqlDB)
	appLogger.Fatal(s.Run())
}
