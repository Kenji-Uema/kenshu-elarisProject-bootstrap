package main

import (
	"context"
	"initMongoDB/domain"
	"initMongoDB/infra"
	"log"
)

func main() {
	ctx := context.Background()

	mongoDb, err := infra.NewMongoDb(ctx, "mongodb://admin:admin123@localhost:32017", "CeladonLuxuryCottage")
	if err != nil {
		log.Fatal(err)
	}

	mongoDb.DropAll(ctx)

	cottageCollection := mongoDb.NewCollection("Cottage")
	guestCollection := mongoDb.NewCollection("Guest")
	mongoDb.NewCollection("Booking")

	infra.SetIndex(cottageCollection, "name")
	infra.SetIndex(guestCollection, "email")

	infra.Seed[domain.Cottage](cottageCollection, "resources/cottages.json")
}
