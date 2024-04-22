package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/leehaowei/hotel-reservation/api"
	"github.com/leehaowei/hotel-reservation/db"
	"github.com/leehaowei/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client       *mongo.Client
	roomStore    db.RoomStore
	hotelStore   db.HotelStore
	userStore    db.UserStore
	bookingStore db.BookingStore
	ctx          = context.Background()
)

func seedUser(isAdmin bool, fname, lname, email, password string) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		Email:     email,
		FirstName: fname,
		LastName:  lname,
		Password:  password,
	})
	if err != nil {
		log.Fatal(err)
	}
	user.IsAdmin = isAdmin
	insertedUser, err := userStore.InsertUsers(ctx, user)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s -> %s\n", user.Email, api.CreateTokenFromUser(user))
	return insertedUser
}

func seedRoom(size string, ss bool, price float64, hotelID primitive.ObjectID) *types.Room {
	room := &types.Room{
		Size:    size,
		Seaside: ss,
		Price:   price,
		HotelID: hotelID,
	}
	insertedRoom, err := roomStore.InsertRoom(ctx, room)
	if err != nil {
		log.Fatal(err)
	}
	return insertedRoom
}

func seedBooking(userID, roomID primitive.ObjectID, from, till time.Time) {
	booking := &types.Booking{
		UserID:   userID,
		RoomID:   roomID,
		FromDate: from,
		TillDate: till,
	}
	resp, err := bookingStore.InsertBooking(ctx, booking)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("booking:", resp.ID)
}

func seedHotel(name string, location string, rating int) *types.Hotel {
	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []primitive.ObjectID{},
		Rating:   rating,
	}
	insertedHotel, err := hotelStore.InsertHotel(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}
	return insertedHotel
}

func main() {
	james := seedUser(false, "james", "foo", "james@foo.com", "supersecurepassword")
	seedUser(true, "admin", "admin", "admin@admin.com", "adminpassword123")
	seedHotel("Bellucia", "France", 3)
	seedHotel("The cozy hotel", "The Netherlands", 4)
	hotel := seedHotel("Dont die in your sleep", "London", 1)
	seedRoom("small", true, 89.99, hotel.ID)
	seedRoom("medium", false, 189.99, hotel.ID)
	room := seedRoom("large", true, 289.99, hotel.ID)
	seedBooking(james.ID, room.ID, time.Now(), time.Now().AddDate(0, 0, 2))
}

func init() {
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
	bookingStore = db.NewMongoBookingStore(client)
}
