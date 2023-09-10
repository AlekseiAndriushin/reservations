package main

import (
	"context"
	"fmt"
	"log"

	"github.com/AlexeyAndryushin/reservations/api"
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
	userStore db.UserStore
	ctx = context.Background()
)

func seedUser(fname, lname, email string) {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		Email: email,
		FirstName: fname,
		LastName: lname,
		Password: "superpasswordhehe123",
	})
	if err != nil {
		log.Fatal(err)
	}

	_, err = userStore.InsertUser(context.TODO(), user)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s -> %s\n", user.Email, api.CreateTokenFromUser(user))

}

func seedHotel(name string, location string, rating int) {
	hotel := types.Hotel{
		Name: name,
		Location: location,
		Rooms: []primitive.ObjectID{},
		Rating: rating,
	}
	rooms := []types.Room{
		{
			Size: "small",
			BasePrice: 88.9,
		},
		{
			Size: "normal",
			BasePrice: 500,
		},
		{
			Size: "big",
			BasePrice: 2500,
		},
	}


	insertedHotel, err := hotelStore.InsertHotel(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}
	for _, room := range rooms {
		room.HotelID = insertedHotel.ID
		insertedRoom, err := roomStore.InsertRoom(ctx, &room)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("room ->%s\n", insertedRoom.ID)
	}
}

func main() {
	seedHotel("Grand hotel", "Azors", 5)
	seedHotel("The advanced hotel", "Madeira", 2)
	seedHotel("Don't drink in your sleep", "London",1 )
	seedUser("alex", "alex", "alex@alex.com")
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
	userStore = db.NewMongoUserStore(client)
}