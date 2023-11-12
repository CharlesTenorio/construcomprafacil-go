package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Entrega struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Pedido     Pedido             `bson:"pedido" json:"pedido"`
	EndEngrega EnderecoEntrega    `bson:"edereco_entrega" json:"edereco_entrega"`
	Status     string             `bson:"status" json:"status"`
}
