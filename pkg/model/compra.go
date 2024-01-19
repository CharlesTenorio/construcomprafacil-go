package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Compra struct {
	ID          primitive.ObjectID `bson:"_id" json:"_id"`
	DataType    string             `bson:"data_type" json:"-"`
	Valor       string             `bson:"valor" json:"valor"`
	Data        time.Time          `bson:"data" json:"data"`
	TituloPagar []struct {
		ID            primitive.ObjectID `bson:"_id" json:"_id"`
		Valor         string             `bson:"valor" json:"valor"`
		Data          time.Time          `bson:"data" json:"data"`
		Fornecedor    string             `bson:"fornecedor" json:"fornecedor"`
		TituloQuitado string             `bson:"titulo_quitado" json:"titulo_quitado"`
	} `bson:"tituloPagar" json:"tituloPagar"`
}
