package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserLogged struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Username      string             `bson:"username" json:"username"`
	RemoteAddress string             `bson:"remote_address" json:"remote_address"`
	AccessToken   string             `bson:"access_token" json:"access_token"`
	RefreshToken  string             `bson:"refresh_token" json:"refresh_token"`
	CreatedAt     string             `bson:"created_at" json:"created_at,omitempty"`
}
