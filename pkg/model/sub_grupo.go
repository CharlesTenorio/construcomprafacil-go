package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type SubGrupo struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	SubGrupo string             `bson:"sub_grupo" json:"sub_grupo"`
	IDGrupo  primitive.ObjectID `bson:"id_grupo" json:"id_grupo"`
	Produtos []Produto          `bson:"produtos" json:"produtos"`
}
