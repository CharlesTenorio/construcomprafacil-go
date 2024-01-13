package dto

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ClienteGrupoDto struct {
	ClienteID primitive.ObjectID `bson:"Cliente_id" json:"Cliente_id"`
}

type FornecedorDto struct {
	FornecedorID primitive.ObjectID `bson:"fornecedor_id" json:"fornecedor_id"`
	Produtos     []ProdutoDto       `bson:"produto" json:"produto"`
}

type ProdutoDto struct {
	ProdutoID         primitive.ObjectID `bson:"produto_id" json:"produto_id"`
	Nome              string             `bson:"nome_produto" json:"nome_produto"`
	CompraID          []CompraDto        `bson:"compra_id" json:"compra_id"`
	Quantidade        int                `bson:"quantidade" json:"quantidade"`
	Valor             float64            `bson:"valor" json:"valor"`
	Desconto          float64            `bson:"desconto" json:"desconto"`
	PrazoEntrega      int                `bson:"prazoEntrega" json:"prazoEntrega"`
	DataEnvio         time.Time          `bson:"dataEnvio" json:"dataEnvio"`
	EstimativaEntrega time.Time          `bson:"estimativaEntrega" json:"estimativaEntrega"`
	DataEntrega       time.Time          `bson:"dataEntrega" json:"dataEntrega"`
	RespondeuCliente  bool               `bson:"respondeuCliente" json:"respondeuCliente"`
	FornecedorRecusou bool               `bson:"fornecedorRecusou" json:"fornecedorRecusou"`
}

type CompraDto struct {
	CompraID primitive.ObjectID `bson:"compra_id" json:"compra_id"`
}

type MeioPagamentoDto struct {
	MeioPagamentoID primitive.ObjectID `bson:"meioPagamento_id" json:"meioPagamento_id"`
}

type OrcamentoFilaPrdFornecedor struct {
	IdOrcamento primitive.ObjectID `bson:"id_orcamento" json:"id_orcamento"`
	Fornecedor  FornecedorDto      `bson:"prd_enviado_fornecedor" json:"prd_enviado_fornecedor"`
	StatusFila  string             `bson:"status" json:"status"`
}

type ProdutoEnvidadosParaContacaoDTO struct {
	ProdutoID  primitive.ObjectID `bson:"produto_id" json:"produto_id"`
	Nome       string             `bson:"nome_produto" json:"nome_produto"`
	Quantidade int                `bson:"quantidade" json:"quantidade"`
}

type ProdutoEnviadoParaFilaDeOrcamentoDTO struct {
	IdOrcamento primitive.ObjectID                `bson:"id_orcamento" json:"id_orcamento"`
	Produtos    []ProdutoEnvidadosParaContacaoDTO `bson:"produtos" json:"produtos"`
	DataEnvio   time.Time                         `bson:"dataEnvio" json:"dataEnvio"`
}
