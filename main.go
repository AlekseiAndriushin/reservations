package main

import (
	"context"
	"flag"
	"log"

	"github.com/AlexeyAndryushin/reservations/api"
	"github.com/AlexeyAndryushin/reservations/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var config = fiber.Config{
	ErrorHandler: api.ErrorHandler,
}

func main() {
	listenAddr := flag.String("listenAddr", ":3000", "The listend address of the API server")
	flag.Parse()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	//hanlers initialization
	var (
	hotelStore = db.NewMongoHotelStore(client)
	roomStore = db.NewMongoRoomStore(client, hotelStore)
	userStore = db.NewMongoUserStore(client)
	bookingStore = db.NewMongoBookingStore(client)
	store = &db.Store{
		Hotel: hotelStore,
		Room: roomStore,
		User: userStore,
		Booking: bookingStore,
	}
	userHandler = api.NewUserHandler(db.NewMongoUserStore(client))
	hotelHandler = api.NewHotelHandler(store)
	authHandler = api.NewAuthHandler(userStore)
	roomHandler = api.NewRoomHandler(store)
	bookingHandler = api.NewBookingHandler(store)

  app = fiber.New(config)
	auth = app.Group("/api")
	apiv1 = app.Group("/api/v1", api.JWTAuthentication(userStore))
	admin = apiv1.Group("/admin", api.AdminAuth)
	)

	//auth
	auth.Post("/auth", authHandler.HandleAuthenticate)

	//Versioned API Routes
	//user hanlders
	apiv1.Put("/user/:id", userHandler.HandlePutUser)
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiv1.Post("/user", userHandler.HandlePostUser)
	apiv1.Get("/user", userHandler.HandleGetUsers)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)

	//hotel handlers
	apiv1.Get("/hotel",hotelHandler.HanldeGetHotels)
	apiv1.Get("/hotel/:id",hotelHandler.HanldeGetHotel)
	apiv1.Get("/hotel/:id/rooms",hotelHandler.HandleGetRooms)

	//rooms handlers
	apiv1.Get("/room", roomHandler.HanlderGetRooms)
	apiv1.Post("/room/:id/book", roomHandler.HandleBookRoom)

	//booking handlers
	apiv1.Get("/booking/:id", bookingHandler.HandleGetBooking)
	apiv1.Get("/booking/:id/cancel", bookingHandler.HandleCancelBooking)
	
	//admin routes
	admin.Get("/booking", bookingHandler.HandleGetBookings)

  app.Listen(*listenAddr)
}
