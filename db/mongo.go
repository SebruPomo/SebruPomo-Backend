package db

import (
	"context"
	"time"

	"github.com/labstack/gommon/log"
	"github.com/sebrupomo/sebrupomo-backend/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var database *mongo.Database
var ctx context.Context

func GetConnection() {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Error("Failed to create MongoDB client: %v", err)
	}

	database = client.Database("SebruPomo")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	err = client.Connect(ctx)
	if err != nil {
		log.Error("Failed to create MongoDB client: %v", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Error("Failed to create MongoDB client: %v", err)
	}
}

func CreateUser(user *model.User) (*model.User, error) {
	user.ID = primitive.NewObjectID()
	_, err := database.Collection("users").InsertOne(ctx, user)
	return user, err
}

func ExistsUserByEmail(email string) bool {
	err := database.Collection("users").FindOne(ctx, bson.D{{Key: "email", Value: email}}).Decode(new(bson.M))
	return err == nil
}

func ExistsUserByUsername(username string) bool {
	err := database.Collection("users").FindOne(ctx, bson.D{{Key: "username", Value: username}}).Decode(new(bson.M))
	return err == nil
}

func FindUserByUsername(username string) (*model.User, error) {
	result := new(model.User)
	err := database.Collection("users").FindOne(ctx, bson.D{{Key: "username", Value: username}}).Decode(result)
	return result, err
}

func FindUserByID(id primitive.ObjectID) (*model.User, error) {
	result := new(model.User)
	err := database.Collection("users").FindOne(ctx, bson.D{{Key: "_id", Value: id}}).Decode(result)
	return result, err
}
