package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Fornecdores struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	NomeForncedor string             `bson:"company_name" json:"company_name"`
	Name          string             `bson:"name" json:"name"`
	CNPJ          string             `bson:"cnpj" json:"cnpj"`
	Endereco      string             `bson:"address" json:"endereco"`
	Complemento   string             `bson:"complement" json:"complemento"`
	Cidade        string             `bson:"city" json:"cidade"`
	Estado        string             `bson:"state" json:"estado"`
	CEP           string             `bson:"cep" json:"cep"`
	Email         string             `bson:"email" json:"email"`
	Ativo         bool               `bson:"enabled" json:"ativo"`
	CriadoEm      string             `bson:"created_at" json:"criado_em,omitempty"`
	AtualizadoEm  string             `bson:"updated_at" json:"atualizado_em,omitempty"`
}
