package mdb

import (
	"context"
	"errors"
	"fmt"

	"github.com/Kenji-Uema/bootstrap/internal/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func createBookingCollection(ctx context.Context, db *mongo.Database, collectionName string) error {
	validator := bson.D{
		{Key: "$jsonSchema", Value: bson.D{
			{Key: "bsonType", Value: "object"},
			{Key: "required", Value: bson.A{"status"}},
			{Key: "properties", Value: bson.D{
				{Key: "status", Value: bson.D{
					{Key: "bsonType", Value: "string"},
					{Key: "enum", Value: bson.A{
						domain.Pending,
						domain.Confirmed,
						domain.Cancelled,
						domain.Past,
					}},
				}},
			}},
		}},
	}

	if err := ensureCollectionValidator(ctx, db, collectionName, validator); err != nil {
		return fmt.Errorf("ensure booking collection: %w", err)
	}

	return nil
}

func createGuestCollection(ctx context.Context, db *mongo.Database, collectionName string) error {
	if err := ensureCollectionExists(ctx, db, collectionName); err != nil {
		return fmt.Errorf("ensure guest collection: %w", err)
	}

	return nil
}

func createInvoiceCollection(ctx context.Context, db *mongo.Database, collectionName string) error {
	validator := bson.D{
		{Key: "$jsonSchema", Value: bson.D{
			{Key: "bsonType", Value: "object"},
			{Key: "required", Value: bson.A{"status"}},
			{Key: "properties", Value: bson.D{
				{Key: "status", Value: bson.D{
					{Key: "bsonType", Value: "string"},
					{Key: "enum", Value: bson.A{
						domain.InvoiceStatusPending,
						domain.InvoiceStatusPaid,
						domain.InvoiceStatusCancelled,
					}},
				}},
			}},
		}},
	}

	if err := ensureCollectionValidator(ctx, db, collectionName, validator); err != nil {
		return fmt.Errorf("ensure invoice collection: %w", err)
	}

	return nil
}

func createReceiptCollection(ctx context.Context, db *mongo.Database, collectionName string) error {
	if err := ensureCollectionExists(ctx, db, collectionName); err != nil {
		return fmt.Errorf("ensure receipt collection: %w", err)
	}

	return nil
}

func createStockCollection(ctx context.Context, db *mongo.Database, collectionName string) error {
	if err := ensureCollectionExists(ctx, db, collectionName); err != nil {
		return fmt.Errorf("ensure stock collection: %w", err)
	}

	return nil
}

func createCottageCollection(ctx context.Context, db *mongo.Database, collectionName string) error {
	validator := bson.D{
		{Key: "$jsonSchema", Value: bson.D{
			{Key: "bsonType", Value: "object"},
			{Key: "required", Value: bson.A{"cleaning_status", "key"}},
			{Key: "properties", Value: bson.D{
				{Key: "cleaning_status", Value: bson.D{
					{Key: "bsonType", Value: "string"},
					{Key: "enum", Value: bson.A{
						domain.PreparedForGuest,
						domain.DailyCleaned,
						domain.PreparedForSleep,
						domain.FullyCleaned,
					}},
				}},
				{Key: "key", Value: bson.D{
					{Key: "bsonType", Value: "object"},
					{Key: "required", Value: bson.A{"holder"}},
					{Key: "properties", Value: bson.D{
						{Key: "holder", Value: bson.D{
							{Key: "bsonType", Value: "string"},
							{Key: "enum", Value: bson.A{
								domain.KeyHolderGuest,
								domain.KeyHolderCottage,
							}},
						}},
					}},
				}},
			}},
		}},
	}

	if err := ensureCollectionValidator(ctx, db, collectionName, validator); err != nil {
		return fmt.Errorf("ensure cottage collection: %w", err)
	}

	return nil
}

func ensureCollectionValidator(ctx context.Context, db *mongo.Database, collectionName string, validator bson.D) error {
	err := db.RunCommand(ctx, bson.D{
		{Key: "create", Value: collectionName},
		{Key: "validator", Value: validator},
		{Key: "validationLevel", Value: "strict"},
		{Key: "validationAction", Value: "error"},
	}).Err()
	if err == nil {
		return nil
	}

	var cmdErr mongo.CommandError
	ok := errors.As(err, &cmdErr)
	if ok && cmdErr.Code == 48 {
		return db.RunCommand(ctx, bson.D{
			{Key: "collMod", Value: collectionName},
			{Key: "validator", Value: validator},
			{Key: "validationLevel", Value: "strict"},
			{Key: "validationAction", Value: "error"},
		}).Err()
	}

	return err
}

func ensureCollectionExists(ctx context.Context, db *mongo.Database, collectionName string) error {
	err := db.RunCommand(ctx, bson.D{
		{Key: "create", Value: collectionName},
	}).Err()
	if err == nil {
		return nil
	}

	var cmdErr mongo.CommandError
	ok := errors.As(err, &cmdErr)
	if ok && cmdErr.Code == 48 {
		return nil
	}

	return err
}
