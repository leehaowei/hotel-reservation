package db

import (
	"context"
	"fmt"
	"os"

	"github.com/leehaowei/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	userColl = "users"
)

type Map map[string]any

type Dropper interface {
	Drop(context.Context) error
}

// UserStore act as an abstraction layer of the db (dao/store)
type UserStore interface {
	Dropper

	GetUserByEmail(context.Context, string) (*types.User, error)
	GetUserByID(context.Context, string) (*types.User, error)
	GetUsers(context.Context) ([]*types.User, error)
	InsertUsers(context.Context, *types.User) (*types.User, error)
	DeleteUsers(context.Context, string) error
	UpdateUsers(ctx context.Context, filter Map, params types.UpdateUserParams) error
}

type MongoUserStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoUserStore(client *mongo.Client) *MongoUserStore {
	dbname := os.Getenv(MongoDBNameEnvName)
	return &MongoUserStore{
		client: client,
		coll:   client.Database(dbname).Collection(userColl),
	}
}

func (s *MongoUserStore) Drop(ctx context.Context) error {
	fmt.Println("--- dropping user collection")
	return s.coll.Drop(ctx)
}

func (s *MongoUserStore) UpdateUsers(ctx context.Context, filter Map, params types.UpdateUserParams) error {
	oid, err := primitive.ObjectIDFromHex(filter["_id"].(string))
	if err != nil {
		return err
	}
	filter["_id"] = oid
	update := bson.M{"$set": params.ToBSON()}
	_, err = s.coll.UpdateOne(ctx, filter, update) // filter: who, update: what
	if err != nil {
		return err
	}
	return nil
}

func (s *MongoUserStore) DeleteUsers(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	//TODO: Maybe its a good idea to handle if we did not delete any user.
	// maybe log it or do something
	_, err = s.coll.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return err
	}
	return nil
}

func (s *MongoUserStore) InsertUsers(ctx context.Context, user *types.User) (*types.User, error) {
	res, err := s.coll.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}
	user.ID = res.InsertedID.(primitive.ObjectID)
	return user, nil
}

func (s *MongoUserStore) GetUsers(ctx context.Context) ([]*types.User, error) {
	cur, err := s.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	var users []*types.User
	if err := cur.All(ctx, &users); err != nil {
		return []*types.User{}, nil
	}
	return users, nil
}

func (s *MongoUserStore) GetUserByEmail(ctx context.Context, email string) (*types.User, error) {
	var user types.User
	if err := s.coll.FindOne(ctx, bson.M{"email": email}).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *MongoUserStore) GetUserByID(ctx context.Context, id string) (*types.User, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var user types.User
	if err := s.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}
