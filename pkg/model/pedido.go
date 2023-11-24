package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Pedido struct {
	ID                   primitive.ObjectID `bson:"_id" json:"_id"`
	ClienteID            primitive.ObjectID `bson:"cliente_id" json:"cliente_id"`
	Descricao            string             `bson:"descricao" json:"descricao"`
	DataSolicitacao      time.Time          `bson:"data" json:"data"`
	PrazoRespostaFor     time.Time          `bson:"dataPrazoFor" json:"dataPrazoFor"`
	PrazoRespostaCli     time.Time          `bson:"dataPrazoCli" json:"dataPrazoCli"`
	SugestaoPrazoEntrega string             `bson:"sugestaoprazoEntrega" json:"sugestaoprazoEntrega"`
	Finalizado           bool               `bson:"finalizado" json:"finalizado"`
	PegarEstabelecimento bool               `bson:"pegarEstabelecimento" json:"pegarEstabelecimento"`
	ListaCliente         []struct {
		ClienteID primitive.ObjectID `bson:"Cliente_id" json:"Cliente_id"`
	} `bson:"listaCliente" json:"listaCliente"`
	EnderecoCliente []Endereco `bson:"enderecoCliente" json:"enderecoCliente"`
	Fornecedor      []struct {
		FornecedorID primitive.ObjectID `bson:"fornecedor_id" json:"fornecedor_id"`
		Produto      []struct {
			ProdutoID primitive.ObjectID `bson:"produto_id" json:"produto_id"`
			CompraID  []struct {
				CompraID primitive.ObjectID `bson:"compra_id" json:"compra_id"`
			} `bson:"compra_id" json:"compra_id"`
			Quantidade        int       `bson:"quantidade" json:"quantidade"`
			Valor             float64   `bson:"valor" json:"valor"`
			Desconto          string    `bson:"desconto" json:"desconto"`
			PrazoEntrega      int       `bson:"prazoEntrega" json:"prazoEntrega"`
			DataEnvio         time.Time `bson:"dataEnvio" json:"dataEnvio"`
			EstimativaEntrega time.Time `bson:"estimativaEntrega" json:"estimativaEntrega"`
			DataEntrega       time.Time `bson:"dataEntrega" json:"dataEntrega"`
			RespondeuCliente  bool      `bson:"respondeuCliente" json:"respondeuCliente"`
			FornecedorRecusou bool      `bson:"fornecedorRecusou" json:"fornecedorRecusou"`
		} `bson:"produto" json:"produto"`
		MeioPagamento []struct {
			MeioPagamentoID primitive.ObjectID `bson:"meioPagamento_id" json:"meioPagamento_id"`
		} `bson:"meioPagamento" json:"meioPagamento"`
	} `bson:"fornecedor" json:"fornecedor"`
}

type FilterPedido struct {
	ClienteID         string
	FornecedorID      string
	DataEnvio         time.Time
	EstimativaEntrega time.Time `bson:"estimativaEntrega" json:"estimativaEntrega"`
	DataEntrega       time.Time
}
