package main

import (
	"context"
	"log"

	"github.com/AlexeyAndryushin/reservations/db"
	"github.com/AlexeyAndryushin/reservations/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client *mongo.Client
	roomStore db.RoomStore
	hotelStore db.HotelStore
	ctx = context.Background()
)

func seedHotel(name, location string) {
	hotel := types.Hotel{
		Name: name,
		Location: location,
		Rooms: []primitive.ObjectID{},
	}
	rooms := []types.Room{
		{
			Type: types.SingleRoomType,
			BasePrice: 88.9,
		},
		{
			Type: types.SeaSideRoomType,
			BasePrice: 500,
		},
		{
			Type: types.DeluxeRoomType,
			BasePrice: 2500,
		},
	}


	insertedHotel, err := hotelStore.InsertHotel(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}
	for _, room := range rooms {
		room.HotelID = insertedHotel.ID
		_, err := roomStore.Insert(ctx, &room)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	seedHotel("Grand hotel", "Azors")
	seedHotel("The advanced hotel", "Madeira")
	seedHotel("Don't drink in your sleep", "London")
}

func init () {
	var err error
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}
	
	hotelStore = db.NewMongoHotelStore(client)
	roomStore = db.NewMongoRoomStore(client, hotelStore)
}