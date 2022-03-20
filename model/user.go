package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id" json:"id"`
	Username string             `json:"username"`
	Email    string             `json:"email"`
	Hash     string             `json:"-"`
}
