package repository

import (
	"context"

	"github.com/f0rmul/sensor-control/internal/models"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	snapshotsDB         = "snapshots"
	snapshotsCollection = "snapshots"
)

type snapshotMongoRepository struct {
	mongoDB *mongo.Client
}

func NewSnapshotRepository(client *mongo.Client) *snapshotMongoRepository {
	return &snapshotMongoRepository{mongoDB: client}
}

func (r *snapshotMongoRepository) Create(ctx context.Context, item *models.Snapshot) (*models.Snapshot, error) {

	collection := r.mongoDB.Database(snapshotsDB).Collection(snapshotsCollection)

	result, err := collection.InsertOne(ctx, item, &options.InsertOneOptions{})

	if err != nil {
		return nil, errors.Wrap(err, "collection.InsertOne()")
	}

	_, ok := result.InsertedID.(primitive.ObjectID)

	if !ok {
		return nil, errors.New("result.InsertedID type conversion") // BLOB
	}
	return item, nil
}

func (r *snapshotMongoRepository) GetByID(ctx context.Context, snapshotID string) (*models.Snapshot, error) {

	collection := r.mongoDB.Database(snapshotsDB).Collection(snapshotsCollection)

	var snapshot models.Snapshot

	if err := collection.FindOne(ctx, bson.M{"snapshot_id": snapshotID}).Decode(&snapshot); err != nil {
		return nil, errors.Wrap(err, "collection.FindOne()")
	}
	return &snapshot, nil
}
