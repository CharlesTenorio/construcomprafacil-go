package model

import (
	"encoding/json"
	"time"

	"github.com/katana/back-end/orcafacil-go/internal/config/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Produto struct {
	ID        primitive.ObjectID `bson:"produto_id" json:"produto_id"`
	Nome      string             `bson:"nome" json:"nome"`
	Enabled   bool               `bson:"enabled" json:"enabled"`
	Categoria Categoria          `bson:"categoria" json:"categoria"`
	CreatedAt string             `bson:"created_at" json:"created_at,omitempty"`
	UpdatedAt string             `bson:"updated_at" json:"updated_at,omitempty"`
}

func (p Produto) ProtutoTOString() string {
	data, err := json.Marshal(p)

	if err != nil {
		logger.Error("error to convert Produto to JSON", err)

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
