package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Produto struct {
	ID        primitive.ObjectID `bson:"produto_id" json:"produto_id"`
	Nome      string             `bson:"nome" json:"nome"`
	Categoria string             `bson:"categoria" json:"categoria"`
}
