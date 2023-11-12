package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Pedido struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	DataPedido     time.Time          `bson:"data_pedido" json:"data_pedido"`
	ClienteID      primitive.ObjectID `bson:"cliente_id" json:"cliente_id"`
	StatusPedido   string             `bson:"status_pedido" json:"status_pedido"`
	DetalhesPedido []DetalhePedido    `bson:"detalhes_pedido" json:"detalhes_pedido"`
	CustoFrete     float64            `bson:"custo_frete" json:"custo_frete"`
	CriadoEm       string             `bson:"criado_em" json:"criado_em,omitempty"`
	AtualizadoEm   string             `bson:"atualizado_em" json:"atualizado_em,omitempty"`
}

type DetalhePedido struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Produto      Produto            `bson:"produto" json:"produto"`
	Quantidade   int                `bson:"quantidade" json:"quantidade"`
	ValorVenda   float64            `bson:"valor_venda" json:"valor_venda"`
	Subtotal     float64            `bson:"subtotal" json:"subtotal"`
	AtualizadoEm string             `bson:"atualizado_em" json:"atualizado_em,omitempty"`
}

type FiltroPedido struct {
	StatusPedido string    `json:"status_pedido"`
	DataPedido   time.Time `bson:"data_pedido" json:"data_pedido"`
}
