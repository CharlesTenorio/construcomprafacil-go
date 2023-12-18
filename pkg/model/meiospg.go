package model

import (
	"encoding/json"
	"time"

	"github.com/katana/back-end/orcafacil-go/internal/config/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MeioPagamento struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Meiopg    string             `bson:"meio_pg" json:"meio_pg"`
	Enabled   bool               `bson:"enabled" json:"enabled"`
	CreatedAt string             `bson:"created_at" json:"created_at,omitempty"`
	UpdatedAt string             `bson:"updated_at" json:"updated_at,omitempty"`
}

func (b MeioPagamento) MeioPG() string {
	data, err := json.Marshal(b)

	if err != nil {
		logger.Error("error to convert Client to JSON", err)

		return ""
	}

	return string(data)
}

type FilterMeioPg struct {
	Meiopg  string `json:"meio_pg"`
	Enabled string `json:"enabled"`
}

func NewMeioPG(client_request MeioPagamento) *MeioPagamento {
	return &MeioPagamento{
		ID: primitive.NewObjectID(),

		Meiopg:    client_request.Meiopg,
		Enabled:   true,
		CreatedAt: time.Now().String(),
	}
}
