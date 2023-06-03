package main

import (
	"context"
	"log"

	"github.com/f0rmul/sensor-control/config"
	"github.com/f0rmul/sensor-control/internal/server"
	"github.com/f0rmul/sensor-control/pkg/logger"
	"github.com/f0rmul/sensor-control/pkg/mongodb"
	"github.com/f0rmul/sensor-control/pkg/rabbitmq"
)

func main() {
	log.Println("[+] Starting server")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.NewConfig("config/config.yml")
	if err != nil {
		log.Fatalf("Loading config: %v", err)
	}

	logger := logger.NewApiLogger(cfg)
	logger.InitLogger()
	logger.Infof("App version: %s, Log level: %s, encoding: %s",
		cfg.Http.AppVersion,
		cfg.Logger.Level,
		cfg.Logger.Encoding)

	amqpConn, err := rabbitmq.NewRabbitMQConn(cfg)

	if err != nil {
		logger.Fatal(err)
	}

	defer amqpConn.Close()

	logger.Info("[+] Connected to message broker")
	mongo, err := mongodb.NewMongoDBConn(ctx, cfg)

	if err != nil {
		logger.Fatal(err)
	}

	defer mongo.Disconnect(ctx)
	logger.Info("[+] Connected to  storage")
	server := server.NewServer(amqpConn, mongo, cfg, logger)

	logger.Fatal(server.Run())
}
