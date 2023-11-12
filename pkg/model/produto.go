package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Produto struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Nome         string             `bson:"name" json:"nome"`
	Ean13        string             `bson:"ean13" json:"ean13"`
	Descrição    string             `bson:"description" json:"descricao"`
	IDGrupo      primitive.ObjectID `bson:"id_grupo" json:"id_grupo"`
	IDSubGrupo   primitive.ObjectID `bson:"id_sub_grupo" json:"id_sub_grupo"`
	UrlImagem    []string           `bson:"imgurl" json:"url_imagem"`
	Ativo        bool               `bson:"enabled" json:"ativo"`
	CriadoEm     string             `bson:"created_at" json:"criado_em,omitempty"`
	AtualizadoEm string             `bson:"updated_at" json:"atualizado_em,omitempty"`
}
