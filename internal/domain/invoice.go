package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type InvoiceStatus string

const (
	InvoiceStatusPending   InvoiceStatus = "pending"
	InvoiceStatusPaid      InvoiceStatus = "paid"
	InvoiceStatusCancelled InvoiceStatus = "cancelled"
)

type Invoice struct {
	Id            bson.ObjectID `bson:"_id,omitempty"`
	InvoiceNumber string        `bson:"invoice_number"`
	Status        InvoiceStatus `bson:"status"`

	IssuedAt time.Time `bson:"issued_at"`
	DueAt    time.Time `bson:"due_at"`

	IdempotencyId string `bson:"idempotency_id"`

	BookingId string `bson:"booking_id"`
	PayerId   string `bson:"payer_id"`
	Payer     Payer  `bson:"payer"`

	Booking       BookingSnapshot `bson:"booking"`
	Total         Money           `bson:"total"`
	TaxTotal      Money           `bson:"tax_total"`
	DiscountTotal Money           `bson:"discount_total"`

	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

type Money struct {
	Amount   int64  `bson:"amount"`
	Currency string `bson:"currency"`
}

type Payer struct {
	Name           string `bson:"name"`
	Email          string `bson:"email"`
	DocumentNumber string `bson:"document_number"`
	BillingAddress string `bson:"billing_address"`
}

type BookingSnapshot struct {
	CottageName    string `bson:"cottage_name"`
	Nights         int32  `bson:"nights"`
	NumberOfGuests int32  `bson:"number_of_guests"`
	ValuePerNight  Money  `bson:"value_per_night"`
}
