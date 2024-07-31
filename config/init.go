package config

import (
	"flag"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"mesaggio-test/pkg/constants"
	"mesaggio-test/pkg/kafka"
	"os"
	"strings"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config.go", "", "API Gateway microservice config.go path")
}

type Config struct {
	ServiceName string
	HttpPort    string
	Postgres    PostgresConfig
	Kafka       *kafka.Config
	KafkaTopics KafkaTopics
}

type PostgresConfig struct {
	PostgresqlHost     string
	PostgresqlPort     string
	PostgresqlUser     string
	PostgresqlPassword string
	PostgresqlDbname   string
	PostgresqlSSLMode  bool
	PgDriver           string
}

type KafkaTopics struct {
	MessageSaved kafka.TopicConfig
}

func InitConfig() (*Config, error) {
	if configPath == "" {
		configPathFromEnv := os.Getenv(constants.ConfigPath)
		if configPathFromEnv != "" {
			configPath = configPathFromEnv
		} else {
			getwd, err := os.Getwd()
			if err != nil {
				return nil, errors.Wrap(err, "os.Getwd")
			}
			configPath = fmt.Sprintf("%s/config/config.yml", getwd)
		}
	}

	cfg := &Config{}

	viper.SetConfigType(constants.Yaml)
	viper.SetConfigFile(configPath)

	if err := viper.ReadInConfig(); err != nil {
		return nil, errors.Wrap(err, "viper.ReadInConfig")
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return nil, errors.Wrap(err, "viper.Unmarshal")
	}

	httpPort := os.Getenv(constants.HttpPort)
	if httpPort != "" {
		cfg.HttpPort = httpPort
	}

	kafkaBrokers := os.Getenv(constants.KafkaBroker)
	if kafkaBrokers != "" {
		cfg.Kafka.Brokers = strings.Split(kafkaBrokers, ",")
	}

	return cfg, nil
}
