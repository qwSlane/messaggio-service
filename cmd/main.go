package main

import (
	"github.com/sirupsen/logrus"
	"log"
	"mesaggio-test/config"
	"mesaggio-test/internal/server"
	"mesaggio-test/pkg/postgres"
)

func main() {
	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatal(err)
	}

	appLogger := logrus.New()

	psqlDB, err := postgres.NewPsqlDB(cfg)
	if err != nil {
		appLogger.Fatalf("Postgresql init: %s", err)
	} else {
		appLogger.Infof("Postgres connected, Status: %#v", psqlDB.Stats())
	}
	defer psqlDB.Close()

	s := server.NewServer(cfg, appLogger, psqlDB)
	appLogger.Fatal(s.Run())
}
