package model

import (
	"encoding/json"
	"time"

	"github.com/katana/back-end/orcafacil-go/internal/config/logger"
	"github.com/katana/back-end/orcafacil-go/internal/dto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Produto struct {
	ID            primitive.ObjectID      `bson:"id" json:"id"`
	Nome          string                  `bson:"nome" json:"nome"`
	ValorUnitario float64                 `bson:"valor_unitario" json:"valor_unitario"`
	Enabled       bool                    `bson:"enabled" json:"enabled"`
	CreatedAt     string                  `bson:"created_at" json:"created_at,omitempty"`
	UpdatedAt     string                  `bson:"updated_at" json:"updated_at,omitempty"`
	Fornecedores  []dto.FornecedoresEmPrd `json:"fornecedores"`
}

func (s Produto) ProdutoConvet() string {
	data, err := json.Marshal(s)

	if err != nil {
		logger.Error("error to convert Client to JSON", err)

		return ""
	}

	return string(data)
}

type FilterProduto struct {
	Nome    string `json:"nome"`
	Enabled string `json:"enabled"`
}

func NewProduto(client_request Produto) *Produto {
	return &Produto{
		ID: primitive.NewObjectID(),

		Nome:      client_request.Nome,
		Enabled:   true,
		CreatedAt: time.Now().String(),
	}
}
