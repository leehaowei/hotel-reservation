package main

import (
	"context"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/leehaowei/hotel-reservation/api"
	"github.com/leehaowei/hotel-reservation/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Configuration
// 1. MongoDB endpoint
// 2. ListenAddress of HTTP server
// 3. JWT secret
// 4. MongoDBName

var config = fiber.Config{
	ErrorHandler: api.ErrorHandler,
}

func main() {
	mongoEndPoint := os.Getenv("MONGO_DB_URL")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoEndPoint))
	if err != nil {
		log.Fatal(err)
	}

	// handlers initialization
	var (
		hotelStore   = db.NewMongoHotelStore(client)
		roomStore    = db.NewMongoRoomStore(client, hotelStore)
		userStore    = db.NewMongoUserStore(client)
		bookingStore = db.NewMongoBookingStore(client)
		store        = &db.Store{
			Hotel:   hotelStore,
			Room:    roomStore,
			User:    userStore,
			Booking: bookingStore,
		}
		userHandler    = api.NewUserHandler(userStore)
		hotelHandler   = api.NewHotelHandler(store)
		authHandler    = api.NewAuthHandler(userStore)
		roomHandler    = api.NewRoomHandler(store)
		BookingHandler = api.NewBookingHandler(store)
		app            = fiber.New(config)
		auth           = app.Group("api")
		apiv1          = app.Group("api/v1", api.JWTAuthentication(userStore))
		admin          = apiv1.Group("/admin", api.AdminAuth)
	)

	// auth
	auth.Post("/auth", authHandler.HandleAuthenticate)

	// Versioned API routes
	// user handlers
	apiv1.Put("/user/:id", userHandler.HandlePutUser)
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiv1.Post("/user", userHandler.HandlePostUser)
	apiv1.Get("/user", userHandler.HandleGetUsers)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)

	//hotel handlers
	apiv1.Get("/hotel", hotelHandler.HandleGetHotels)
	apiv1.Get("/hotel/:id", hotelHandler.HandleGetHotel)
	apiv1.Get("/hotel/:id/rooms", hotelHandler.HandleGetRooms)

	// rooms handlers
	apiv1.Get("/room", roomHandler.HandleGetRooms)
	apiv1.Post("/room/:id/book", roomHandler.HandleBookRoom)

	// bookings handlers
	apiv1.Get("/booking/:id", BookingHandler.HandleGetBooking)
	apiv1.Get("/booking/:id/cancel", BookingHandler.HandleCancelBookings)

	// admin handlers
	admin.Get("/booking", BookingHandler.HandleGetBookings)

	listenAddr := os.Getenv("HTTP_LISTEN_ADDRESS")
	app.Listen(listenAddr)
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal()
	}
}
