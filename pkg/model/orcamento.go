package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Orcamento struct {
	ID                 primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	ClienteID          primitive.ObjectID `bson:"cliente_id" json:"fornecedor_id"`
	DetalheOrecamentos []DetalheOrcamento `bson:"detalhes,omitempty" json:"detalhes,omitempty"`
	Validade           time.Time          `bson:"validade" json:"valido"`
	StatusOrcamento    string             `bson:"status_orcamento" json:"status_orcamento"`
	AtualizadoEm       string             `bson:"atualizado_em" json:"atualizado_em,omitempty"`
	CriadoEm           string             `bson:"criado_em" json:"criado_em,omitempty"`
	StatusOrsamento    string             `bson:"status_orsamento" json:"status_orsamento"`
}

type DetalheOrcamento struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Quantidade int                `bson:"quantidade" json:"quantidade"`
	ProdutoID  primitive.ObjectID `bson:"produto_id" json:"produto_id"`
}
