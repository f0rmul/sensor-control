package server

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/alash3al/go-pubsub"
	"github.com/f0rmul/sensor-control/config"
	rabbit "github.com/f0rmul/sensor-control/internal/delivery/amqp"
	"github.com/f0rmul/sensor-control/internal/delivery/ws"
	"github.com/f0rmul/sensor-control/internal/repository"
	"github.com/f0rmul/sensor-control/internal/usecase"
	"github.com/f0rmul/sensor-control/pkg/httpserver"
	"github.com/f0rmul/sensor-control/pkg/logger"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/mongo"
)

type Server struct {
	amqpConn    *amqp.Connection
	mongoClient *mongo.Client
	cfg         *config.Config
	logger      logger.Logger
}

func NewServer(amqpConn *amqp.Connection, mongoClient *mongo.Client, cfg *config.Config, logger logger.Logger) *Server {
	return &Server{amqpConn: amqpConn, mongoClient: mongoClient, cfg: cfg, logger: logger}
}

func (s *Server) Run() error {

	s.logger.Info("Initializing components")

	snapshotRepo := repository.NewSnapshotRepository(s.mongoClient)

	broker := pubsub.NewBroker()
	snapshotUsecase := usecase.NewSnapshotUsecase(snapshotRepo, broker, s.logger)
	amqpConsumer := rabbit.NewSnapshotConsumer(s.amqpConn, snapshotUsecase, s.logger)
	notifier := ws.NewPushNotifier(snapshotUsecase, s.logger)

	mux := http.NewServeMux()
	mux.HandleFunc("/stream", notifier.HandlePush)

	httpServer := httpserver.New(mux, httpserver.Addr(s.cfg.Http.Host, s.cfg.Http.Port))

	go func() {
		err := amqpConsumer.StartConsumer(
			s.cfg.RabbitMQ.WorkerPool,
			s.cfg.RabbitMQ.Exchange,
			s.cfg.RabbitMQ.Queue,
			s.cfg.RabbitMQ.RoutingKey,
			s.cfg.RabbitMQ.ConsumerTag,
		)
		if err != nil {
			s.logger.Errorf("StartConsumer(): %v", err)
			return
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case sig := <-quit:
		s.logger.Infof("app - Run - signal: %s", sig.String())
	case err := <-httpServer.Notify():
		s.logger.Errorf("app - Run - httpServer.Notify: %w", err)
	}

	if err := httpServer.Shutdown(); err != nil {
		return err
	}

	s.logger.Info("App was quited successfuly")
	return nil
}
