package model

import (
	"encoding/json"
	"time"

	"github.com/katana/back-end/orcafacil-go/internal/config/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Categoria struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Nome      string             `bson:"nome" json:"nome"`
	Enabled   bool               `bson:"enabled" json:"enabled"`
	CreatedAt string             `bson:"created_at" json:"created_at,omitempty"`
	UpdatedAt string             `bson:"updated_at" json:"updated_at,omitempty"`
}

func (c Categoria) CategoriaConvet() string {
	data, err := json.Marshal(c)

	if err != nil {
		logger.Error("error to convert Client to JSON", err)

		return ""
	}

	return string(data)
}

type FilterCategoria struct {
	Nome    string `json:"nome"`
	Enabled string `json:"enabled"`
}

func NewCategoria(client_request Categoria) *Categoria {
	return &Categoria{
		ID: primitive.NewObjectID(),

		Nome:      client_request.Nome,
		Enabled:   true,
		CreatedAt: time.Now().String(),
	}
}
