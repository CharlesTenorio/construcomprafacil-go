package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type GrupoProduto struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Grupo     string             `bson:"grupo" json:"grupo"`
	SubGrupos []SubGrupo         `bson:"sub_grupos" json:"sub_grupos"`
}
