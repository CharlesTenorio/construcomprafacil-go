package catalgo

import (
	"github.com/katana/back-end/orcafacil-go/internal/dto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CategoriaDto struct {
	ID       primitive.ObjectID         `bson:"_id,omitempty" json:"_id,omitempty"`
	Nome     string                     `bson:"nome" json:"nome"`
	Enabled  bool                       `bson:"enabled" json:"enabled"`
	Produtos []dto.ProdutosEmCategorias `bson:"produtos" json:"produtos"`
}
