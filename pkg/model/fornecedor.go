package model

import (
	"encoding/json"
	"time"

	"github.com/katana/back-end/orcafacil-go/internal/config/logger"
	"github.com/katana/back-end/orcafacil-go/internal/dto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Fornecedor struct {
	ID            primitive.ObjectID         `bson:"_id" json:"_id"`
	DataType      string                     `bson:"data_type" json:"-"`
	IDUsuario     primitive.ObjectID         `bson:"user_id " json:"id_usr"`
	Nome          string                     `bson:"nome" json:"nome"`
	Telefone      string                     `bson:"telefone" json:"telefone"`
	CNPJ          string                     `bson:"cnpj" json:"cnpj"`
	Raio          string                     `bson:"raio" json:"raio"`
	Excluido      string                     `bson:"excluido" json:"excluido"`
	Endereco      []Endereco                 `bson:"endereco" json:"endereco"`
	MeioPagamento []MeioPagamento            `bson:"meio_pagamento" json:"meio_pagamento"`
	Produtos      []dto.ProdutosEmFornecedor `json:"produtos"`
	Enabled       bool                       `bson:"enabled" json:"enabled"`
	CreatedAt     string                     `bson:"created_at" json:"created_at,omitempty"`
	UpdatedAt     string                     `bson:"updated_at" json:"updated_at,omitempty"`
}

func (f Fornecedor) FornecedorConvet() string {
	data, err := json.Marshal(f)

	if err != nil {
		logger.Error("error to convert Client to JSON", err)

		return ""
	}

	return string(data)
}

type FilterFornecedor struct {
	Nome      string             `json:"nome"`
	IDUsuario primitive.ObjectID `bson:"user_id " json:"id_usr"`
	CNPJ      string             `bson:"cnpj" json:"cnpj"`
	Enabled   string             `json:"enabled"`
}

func NewFornecedor(fonecedor_request Fornecedor) *Fornecedor {
	return &Fornecedor{
		ID:        primitive.NewObjectID(),
		DataType:  "fornecedor",
		IDUsuario: fonecedor_request.IDUsuario,
		Nome:      fonecedor_request.Nome,
		Enabled:   true,
		CreatedAt: time.Now().String(),
	}
}
