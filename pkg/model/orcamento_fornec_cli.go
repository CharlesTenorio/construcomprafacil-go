package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrcamentoFornecCli struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	OrcamentoID     primitive.ObjectID `bson:"cliente_id" json:"fornecedor_id"`
	DetalheForCli   []DetalheForCli    `bson:"detalhes,omitempty" json:"detalhes,omitempty"`
	Validade        time.Time          `bson:"validade" json:"valido"`
	StatusOrcamento string             `bson:"status_orcamento" json:"status_orcamento"`
	AtualizadoEm    string             `bson:"atualizado_em" json:"atualizado_em,omitempty"`
	CriadoEm        string             `bson:"criado_em" json:"criado_em,omitempty"`
	StatusOrsamento string             `bson:"status_orsamento" json:"status_orsamento"`
}

type DetalheForCli struct {
	ID                  primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FornecedorID        primitive.ObjectID `bson:"fornecedor_id" json:"fornecedor_id"`
	Quantidade          int                `bson:"quantidade" json:"quantidade"`
	ProdutoID           primitive.ObjectID `bson:"produto_id" json:"produto_id"`
	ImgUrl              string             `bson:"img_url" json:"img_url"`
	Descricao           string             `bson:"descricao" json:"descricao"`
	ValorUnitario       float64            `bson:"ValorUnitario" json:"ValorUnitario"`
	UnidadeMedida       string             `bson:"unidade medida" json:"unidade medida"`
	QuantidadeAtenidida int                `bson:"quantidade atendida" json:"quantidade atendida"`
}
