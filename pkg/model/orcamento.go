package model

import (
	"time"

	"github.com/katana/back-end/orcafacil-go/internal/dto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Orcamento struct {
	ID                   primitive.ObjectID                    `bson:"_id" json:"_id"`
	DataType             string                                `bson:"data_type" json:"-"`
	ClienteID            primitive.ObjectID                    `bson:"cliente_id" json:"cliente_id"`
	Descricao            string                                `bson:"descricao" json:"descricao"`
	DataSolicitacao      time.Time                             `bson:"data_solicitacao" json:"data_solicitacao"`
	PrazoRespostaFor     time.Time                             `bson:"dataPrazoFor" json:"dataPrazoFor"`
	PrazoRespostaCli     time.Time                             `bson:"dataPrazoCli" json:"dataPrazoCli"`
	SugestaoPrazoEntrega string                                `bson:"sugestaoprazoEntrega" json:"sugestaoprazoEntrega"`
	Finalizado           bool                                  `bson:"finalizado" json:"finalizado"`
	PegarEstabelecimento bool                                  `bson:"pegarEstabelecimento" json:"pegarEstabelecimento"`
	GrupoDeCliente       []dto.ClienteGrupoDto                 `bson:"listaCliente" json:"listaCliente"`
	EnderecoCliente      []Endereco                            `bson:"enderecoCliente" json:"enderecoCliente"`
	Fornecedores         []dto.FornecedorDto                   `bson:"fornecedor" json:"fornecedor"`
	ProdutosContacao     []dto.ProdutoEnvidadosParaContacaoDTO `bson:"produtosContacao" json:"produtosContacao"`
	Enabled              bool                                  `bson:"enabled" json:"enabled"`
	CreatedAt            string                                `bson:"created_at" json:"created_at,omitempty"`
	UpdatedAt            string                                `bson:"updated_at" json:"updated_at,omitempty"`
	StatusOrecamento     string                                `bson:"status" json:"status"`
}

type FilterOrcamento struct {
	OrcamentoID       string
	DataType          string
	ClienteID         string
	FornecedorID      string
	DataEnvio         time.Time
	EstimativaEntrega time.Time
	DataEntrega       time.Time
	Enabled           string
	DataInical        string
	DataFinal         string
}

type OrcamentoResult struct {
	OrcamentoID   string
	MeioPagamento string
	Total         float64
}

func NewOrcamento(orcamento_request Orcamento) *Orcamento {
	return &Orcamento{
		ID:               primitive.NewObjectID(),
		DataType:         "orcamento",
		ClienteID:        orcamento_request.ClienteID,
		StatusOrecamento: "enviado para contacao",
		Enabled:          true,
		CreatedAt:        time.Now().String(),
	}
}
