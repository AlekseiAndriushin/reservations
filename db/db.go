package db

const (
	DBNAME = "hotel-reservation"
	TestDBNAME = "hotel-reservation-test"
 	DBURI = "mongodb://localhost:27017"
	MongoDBNameEnvName = "MONGO_DB_NAME"
)


type Store struct {
	User UserStore
	Hotel HotelStore
	Room RoomStore
	Booking BookingStore
}