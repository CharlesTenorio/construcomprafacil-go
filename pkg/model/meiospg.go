package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type MeioPagamento struct {
	ID primitive.ObjectID `json:"meio_pagamento_id"`
}
