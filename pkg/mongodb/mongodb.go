package mongodb

import (
	"context"
	"fmt"
	"time"

	"github.com/f0rmul/sensor-control/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	connectionTimeout = 30 * time.Second
	maxConnIdleTime   = 3 * time.Second
	minPoolSize       = 20
	maxPoolSize       = 300
)

func NewMongoDBConn(ctx context.Context, cfg *config.Config) (*mongo.Client, error) {

	client, err := mongo.NewClient(
		options.Client().
			ApplyURI(fmt.Sprintf("mongodb://%s:%s", cfg.MongoDB.Host, cfg.MongoDB.Port)).
			SetAuth(
				options.Credential{
					Username: cfg.MongoDB.User,
					Password: cfg.MongoDB.Password,
				}).
			SetConnectTimeout(connectionTimeout).
			SetMaxConnIdleTime(maxConnIdleTime).
			SetMinPoolSize(minPoolSize).
			SetMaxPoolSize(maxPoolSize))

	if err != nil {
		return nil, err
	}

	if err := client.Connect(ctx); err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return client, nil
}
