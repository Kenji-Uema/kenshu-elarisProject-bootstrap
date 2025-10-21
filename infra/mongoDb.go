package infra

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoDb struct {
	client   *mongo.Client
	database *mongo.Database
}

func NewMongoDb(connectionContext context.Context, uri, dbName string) (*MongoDb, error) {
	connectionContext, connectionCancel := context.WithTimeout(connectionContext, 10*time.Second)
	defer connectionCancel()

	client, err := mongo.Connect(connectionContext, options.Client().ApplyURI(uri))

	if err != nil {
		return nil, err
	}

	databaseContext, databaseCancel := context.WithTimeout(connectionContext, 5*time.Second)
	defer databaseCancel()

	if err := client.Ping(databaseContext, readpref.Primary()); err != nil {
		_ = client.Disconnect(context.Background())
		return nil, fmt.Errorf("mongo ping failed for URI: %s, error: %w", uri, err)
	}

	return &MongoDb{client: client, database: client.Database(dbName)}, nil
}

func (db *MongoDb) NewCollection(name string) *mongo.Collection {
	return db.Collection(name)
}

func SetIndex(collection *mongo.Collection, fieldName string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	idx := mongo.IndexModel{
		Keys:    bson.D{{Key: fieldName, Value: 1}},
		Options: options.Index().SetUnique(true).SetName(fieldName),
	}

	_, err := collection.Indexes().CreateOne(ctx, idx)
	if err != nil {
		log.Fatal(err)
	}
}

func (db *MongoDb) Close(ctx context.Context) error {
	return db.client.Disconnect(ctx)
}

func (db *MongoDb) Collection(name string) *mongo.Collection {
	return db.database.Collection(name)
}

func (db *MongoDb) DropAll(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	err := db.database.Drop(ctx)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Database dropped")
}

func Seed[I interface{}](collection *mongo.Collection, filepath string) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatal(err)
	}

	var items []I
	if err := json.Unmarshal(data, &items); err != nil {
		log.Fatal(err)
	}

	docs := make([]interface{}, len(items))
	for i, g := range items {
		docs[i] = g
	}
	if _, err := collection.InsertMany(context.Background(), docs); err != nil {
		log.Fatal(err)
	}
}
