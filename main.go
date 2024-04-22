package main

import (
	"context"
	"flag"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/leehaowei/hotel-reservation/api"
	"github.com/leehaowei/hotel-reservation/db"
	"github.com/leehaowei/hotel-reservation/middleware"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var config = fiber.Config{
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		return c.JSON(map[string]string{"error": err.Error()})
	},
}

func main() {
	listenAddr := flag.String("listenAddr", ":5000", "The listen address of the API server")
	flag.Parse()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
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
		apiv1          = app.Group("api/v1", middleware.JWTAuthentication(userStore))
		admin          = apiv1.Group("/admin", middleware.AdminAuth)
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
	// TODO: cancel a booking

	// bookings handlers
	apiv1.Get("/booking/:id", BookingHandler.HandleGetBooking)
	apiv1.Get("/booking/:id/cancel", BookingHandler.HandleCancelBookings)

	// admin handlers
	admin.Get("/booking", BookingHandler.HandleGetBookings)

	app.Listen(*listenAddr)
}
