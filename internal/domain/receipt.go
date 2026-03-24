package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Receipt struct {
	Id            bson.ObjectID `bson:"_id,omitempty"`
	ReceiptNumber string        `bson:"receipt_number"`
	InvoiceNumber string        `bson:"invoice_number"`
	Card          CardSummary   `bson:"card"`
	ProcessedAt   time.Time     `bson:"processed_at"`
}

type CardSummary struct {
	Brand string `bson:"brand"`
	Last4 string `bson:"last4"`
}
