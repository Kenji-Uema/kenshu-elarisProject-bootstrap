package infra

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/Kenji-Uema/bootstrap/internal/config"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

type MongoDB struct {
	client   *mongo.Client
	database *mongo.Database
}

func NewMongoDB(ctx context.Context, mongoConfig config.MongoConfig) (*MongoDB, error) {
	uri := fmt.Sprintf("mongodb://%s:%s@%s", string(mongoConfig.Username), string(mongoConfig.Password), mongoConfig.Host)

	ctx, connectionCancel := context.WithTimeout(ctx, 10*time.Second)
	defer connectionCancel()

	client, err := mongo.Connect(options.Client().ApplyURI(uri))

	if err != nil {
		return nil, err
	}

	databaseContext, databaseCancel := context.WithTimeout(ctx, 5*time.Second)
	defer databaseCancel()

	if err := client.Ping(databaseContext, readpref.Primary()); err != nil {
		disconnectCtx, disconnectCancel := context.WithTimeout(ctx, 5*time.Second)
		defer disconnectCancel()
		_ = client.Disconnect(disconnectCtx)
		return nil, fmt.Errorf("mongo ping failed for URI: %s, error: %w", uri, err)
	}

	return &MongoDB{client: client, database: client.Database(mongoConfig.Database)}, nil
}

func (db *MongoDB) NewCollection(name string) *mongo.Collection {
	return db.Collection(name)
}

func SetIndex(ctx context.Context, collection *mongo.Collection, fieldName string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	idx := mongo.IndexModel{
		Keys:    bson.D{{Key: fieldName, Value: 1}},
		Options: options.Index().SetUnique(true).SetName(fieldName),
	}

	_, err := collection.Indexes().CreateOne(ctx, idx)
	if err != nil {
		return fmt.Errorf("failed to create index %q: %w", fieldName, err)
	}
	return nil
}

func (db *MongoDB) Close(ctx context.Context) error {
	return db.client.Disconnect(ctx)
}

func (db *MongoDB) Collection(name string) *mongo.Collection {
	return db.database.Collection(name)
}

func (db *MongoDB) DropAll(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	err := db.database.Drop(ctx)

	if err != nil {
		return fmt.Errorf("failed to drop database: %w", err)
	}
	slog.Info("database dropped")
	return nil
}

func Seed[I interface{}](ctx context.Context, collection *mongo.Collection, filepath string) error {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("failed to read seed file %q: %w", filepath, err)
	}

	var items []I
	if err := json.Unmarshal(data, &items); err != nil {
		return fmt.Errorf("failed to unmarshal seed data %q: %w", filepath, err)
	}

	docs := make([]interface{}, len(items))
	for i, g := range items {
		docs[i] = g
	}

	seedCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if _, err := collection.InsertMany(seedCtx, docs); err != nil {
		return fmt.Errorf("failed to insert seed data %q: %w", filepath, err)
	}

	return nil
}
