package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Categoria struct {
	ID   primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Nome string             `bson:"nome" json:"nome"`
}
