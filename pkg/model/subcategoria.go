package model

import (
	"encoding/json"
	"time"

	"github.com/katana/back-end/orcafacil-go/internal/config/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Subcategoria struct {
	ID            primitive.ObjectID `bson:"id" json:"id"`
	Nome          string             `bson:"nome" json:"nome"`
	ValorUnitario float64            `bson:"valor_unitario" json:"valor_unitario"`
	Enabled       bool               `bson:"enabled" json:"enabled"`
	CreatedAt     string             `bson:"created_at" json:"created_at,omitempty"`
	UpdatedAt     string             `bson:"updated_at" json:"updated_at,omitempty"`
}

func (s Subcategoria) SubCategoriaConvet() string {
	data, err := json.Marshal(s)

	if err != nil {
		logger.Error("error to convert Client to JSON", err)

		return ""
	}

	return string(data)
}

type FilterSubcategoria struct {
	Nome    string `json:"nome"`
	Enabled string `json:"enabled"`
}

func NewSubCategoria(client_request Subcategoria) *Subcategoria {
	return &Subcategoria{
		ID: primitive.NewObjectID(),

		Nome:      client_request.Nome,
		Enabled:   true,
		CreatedAt: time.Now().String(),
	}
}
