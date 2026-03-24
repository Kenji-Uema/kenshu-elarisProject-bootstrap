package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type BookingStatus string

const (
	Pending   BookingStatus = "pending"
	Confirmed BookingStatus = "confirmed"
	Cancelled BookingStatus = "cancelled"
	Past      BookingStatus = "past"
)

type Period struct {
	CheckIn  time.Time `bson:"checkin"`
	checkOut time.Time `bson:"checkout"`
}

type Booking struct {
	Id             bson.ObjectID `bson:"_id,omitempty"`
	MainGuest      bson.ObjectID `bson:"main_guest"`
	NumberOfGuests int           `bson:"number_of_guests"`
	StayPeriod     Period        `bson:"stay_period"`
	CottageName    string        `bson:"cottage_name"`
	Status         BookingStatus `bson:"status"`
}
