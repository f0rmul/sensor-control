package usecase

import (
	"context"

	"github.com/alash3al/go-pubsub"
	"github.com/f0rmul/sensor-control/internal/models"
	"github.com/f0rmul/sensor-control/pkg/logger"
	"github.com/google/uuid"
)

const (
	brokerTopic = "snapshots"
)

type SnapshotRepository interface {
	Create(context.Context, *models.Snapshot) (*models.Snapshot, error)
	GetByID(context.Context, string) (*models.Snapshot, error)
}

type snapshotUsecase struct {
	snaphostRepo SnapshotRepository
	broker       *pubsub.Broker
	logger       logger.Logger
}

func NewSnapshotUsecase(repo SnapshotRepository, broker *pubsub.Broker, logger logger.Logger) *snapshotUsecase {
	return &snapshotUsecase{snaphostRepo: repo, broker: broker, logger: logger}
}

func (s *snapshotUsecase) PushAndSave(ctx context.Context, item *models.Snapshot) (*models.Snapshot, error) {
	s.broker.Broadcast(item, brokerTopic)

	item.ID = uuid.New().String()
	s.logger.Infof("Snapshot with ID: %s was published", item.ID)

	return s.snaphostRepo.Create(ctx, item)
}

func (s *snapshotUsecase) GetByID(ctx context.Context, snapshotID string) (*models.Snapshot, error) {
	s.logger.Infof("Fetching snapshot with ID: %s", snapshotID)
	return s.snaphostRepo.GetByID(ctx, snapshotID)
}

func (s *snapshotUsecase) Broker() *pubsub.Broker {
	return s.broker
}
