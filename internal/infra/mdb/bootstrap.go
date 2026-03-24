package mdb

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/Kenji-Uema/bootstrap/internal/config"
	"github.com/Kenji-Uema/bootstrap/internal/domain"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func BootstrapMongodb(baseCtx context.Context, configs config.Configs) (func(), error) {
	mongoDb, err := NewMongoDB(baseCtx, configs.MongoConfig)
	if err != nil {
		return nil, err
	}
	shouldCloseOnError := true
	defer func() {
		if !shouldCloseOnError {
			return
		}

		closeCtx, cancel := context.WithTimeout(baseCtx, 5*time.Second)
		defer cancel()
		if err := mongoDb.Close(closeCtx); err != nil {
			slog.Error("failed to close mongo connection", "error", err)
		}
	}()

	if err := mongoDb.DropAll(baseCtx); err != nil {
		return nil, err
	}

	if err := CreateCollections(baseCtx, mongoDb.Database(), configs.MongoConfig); err != nil {
		return nil, err
	}

	if err := SetIndexes(baseCtx, mongoDb, configs); err != nil {
		return nil, err
	}

	if err := SeedCollections(baseCtx, mongoDb, configs); err != nil {
		return nil, err
	}
	shouldCloseOnError = false

	return func() {
		closeCtx, cancel := context.WithTimeout(baseCtx, 5*time.Second)
		defer cancel()
		if err := mongoDb.Close(closeCtx); err != nil {
			slog.Error("failed to close mongo connection", "error", err)
		}
	}, err
}

func CreateCollections(ctx context.Context, db *mongo.Database, mongoConfig config.MongoConfig) error {
	return errors.Join(
		createCottageCollection(ctx, db, mongoConfig.Collections.CottageCollection),
		createGuestCollection(ctx, db, mongoConfig.Collections.GuestCollection),
		createBookingCollection(ctx, db, mongoConfig.Collections.BookingCollection),
		createInvoiceCollection(ctx, db, mongoConfig.Collections.InvoiceCollection),
		createReceiptCollection(ctx, db, mongoConfig.Collections.ReceiptCollection),
		createStockCollection(ctx, db, mongoConfig.Collections.StockCollection),
	)
}

func SetIndexes(ctx context.Context, mongoDb *MongoDB, configs config.Configs) error {
	cottageCollection := mongoDb.GetCollection(configs.MongoConfig.Collections.CottageCollection)
	guestCollection := mongoDb.GetCollection(configs.MongoConfig.Collections.GuestCollection)
	bookingCollection := mongoDb.GetCollection(configs.MongoConfig.Collections.BookingCollection)
	invoiceCollection := mongoDb.GetCollection(configs.MongoConfig.Collections.InvoiceCollection)
	receiptCollection := mongoDb.GetCollection(configs.MongoConfig.Collections.ReceiptCollection)

	return errors.Join(
		SetIndex(ctx, cottageCollection, "name"),
		SetIndex(ctx, guestCollection, "document_id"),
		SetIndex(ctx, guestCollection, "email"),
		SetNonUniqueIndex(ctx, bookingCollection, "main_guest"),
		SetIndex(ctx, invoiceCollection, "invoice_number"),
		SetIndex(ctx, invoiceCollection, "idempotency_id"),
		SetIndex(ctx, receiptCollection, "receipt_number"),
		SetIndex(ctx, receiptCollection, "invoice_number"),
	)
}

func SeedCollections(ctx context.Context, mongoDb *MongoDB, configs config.Configs) error {
	cottageCollection := mongoDb.GetCollection(configs.MongoConfig.Collections.CottageCollection)
	stockCollection := mongoDb.GetCollection(configs.MongoConfig.Collections.StockCollection)

	return errors.Join(
		Seed[domain.Cottage](ctx, cottageCollection, "resources/cottages.json"),
		Seed[domain.Stock](ctx, stockCollection, "resources/stocks.json"),
	)
}
