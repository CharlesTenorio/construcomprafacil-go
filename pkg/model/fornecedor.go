package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Fornecedor struct {
	ID            primitive.ObjectID `bson:"_id" json:"_id"`
	Nome          string             `bson:"nome" json:"nome"`
	Email         string             `bson:"email" json:"email"`
	Telefone      string             `bson:"telefone" json:"telefone"`
	CNPJ          string             `bson:"cnpj" json:"cnpj"`
	Raio          string             `bson:"raio" json:"raio"`
	DataCadastro  time.Time          `bson:"data" json:"data"`
	Senha         string             `bson:"senha" json:"senha"`
	Excluido      string             `bson:"excluido" json:"excluido"`
	Endereco      []Endereco         `bson:"endereco" json:"endereco"`
	MeioPagamento []MeioPagamento    `bson:"meio_pagamento" json:"meio_pagamento"`
	Produto       []Produto          `bson:"produto" json:"produto"`
}
