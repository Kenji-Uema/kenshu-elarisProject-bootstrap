package domain

import "go.mongodb.org/mongo-driver/v2/bson"

type Stock struct {
	Id                bson.ObjectID `bson:"_id,omitempty"`
	CleaningItems     StockItem     `bson:"cleaning_items"`
	Soap              StockItem     `bson:"soap"`
	BathroomAmenities StockItem     `bson:"bathroom_amenities"`
	AromaCandles      StockItem     `bson:"aroma_candles"`
	WaterBottle       StockItem     `bson:"water_bottle"`
	WineBottle        StockItem     `bson:"wine_bottle"`
	TeaBags           StockItem     `bson:"tea_bags"`
	Sweets            StockItem     `bson:"sweets"`
	Chips             StockItem     `bson:"chips"`
}

type StockItem struct {
	Name     string `bson:"name"`
	Quantity int    `bson:"quantity"`
}
