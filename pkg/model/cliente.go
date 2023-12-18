package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Cliente struct {
	ID            primitive.ObjectID `bson:"_id" json:"_id"`
	Nome          string             `bson:"nome" json:"nome"`
	Email         string             `bson:"email" json:"email"`
	Sexo          string             `bson:"sexo" json:"sexo"`
	Telefone      string             `bson:"telefone" json:"telefone"`
	Tipo          string             `bson:"tipo" json:"tipo"`
	CPFCNPJ       string             `bson:"cpf_cnpj" json:"cpf_cnpj"`
	DataCadastro  time.Time          `bson:"data" json:"data"`
	Senha         string             `bson:"senha" json:"senha"`
	Excluido      string             `bson:"excluido" json:"excluido"`
	MeioPagamento []MeioPagamento    `bson:"meioPagamento" json:"meioPagamento"`
}
