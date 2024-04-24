package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/leehaowei/hotel-reservation/api"
	"github.com/leehaowei/hotel-reservation/db"
	"github.com/leehaowei/hotel-reservation/db/fixtures"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}
	hotelStore := db.NewMongoHotelStore(client)

	store := &db.Store{
		User:    db.NewMongoUserStore(client),
		Booking: db.NewMongoBookingStore(client),
		Room:    db.NewMongoRoomStore(client, hotelStore),
		Hotel:   hotelStore,
	}

	user := fixtures.AddUser(store, "james", "foo", false)
	admin := fixtures.AddUser(store, "admin", "admin", true)
	hotel := fixtures.AddHotel(store, "some hotel", "bermuda", 5, nil)
	room := fixtures.AddRoom(store, "large", true, 88.44, hotel.ID)
	booking := fixtures.AddBooking(store, user.ID, room.ID, time.Now(), time.Now().AddDate(0, 0, 5))

	fmt.Println("james ->", api.CreateTokenFromUser(user))
	fmt.Println("admin ->", api.CreateTokenFromUser(admin))
	fmt.Println("booking ->", booking.ID)

	for i := 0; i < 100; i++ {
		name := fmt.Sprintf("random_hotel name %d", i)
		location := fmt.Sprintf("location %d", i)
		fixtures.AddHotel(store, name, location, rand.Intn(5)+1, nil)
	}
}
