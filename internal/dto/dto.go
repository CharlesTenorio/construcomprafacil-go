package dto

import (
	"time"

	"github.com/katana/back-end/orcafacil-go/pkg/service/validation"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GetJwtInput struct {
	Email string `json:"email"`
	Senha string `json:"senha"`
	Role  string `json:"role"`
}

type GetJWTOutput struct {
	AccessToken string `json:"access_token"`
}

type FornecedoresEmPrd struct {
	ID         string  `json:"id"`
	PrecoVenda float64 `json:"preco_venda"`
}

type ProdutosEmFornecedor struct {
	ID          string `json:"id"`
	DataType    string `bson:"data_type" json:"-"`
	IDCategoria string `json:"id_categoria"`
	Nome        string `json:"nome"`
	Enabled     bool   `json:"enabled"`
	CreatedAt   string `bson:"created_at" json:"created_at,omitempty"`
	UpdatedAt   string `bson:"updated_at" json:"updated_at,omitempty"`
}

type ProdutosEmCategorias struct {
	ID      primitive.ObjectID `bson:"id" json:"id"`
	Nome    string             `bson:"nome" json:"nome"`
	Enabled bool               `bson:"enabled" json:"enabled"`
}

type ProdutosPayload struct {
	Produtos []ProdutosEmFornecedor `json:"produtos"`
}

type FornecedorPaylaod struct {
	Fornecedores []FornecedoresEmPrd `json:"fornecedores"`
}

func NewProdutosEmFornecedor(prdFornec_request ProdutosEmFornecedor) *ProdutosEmFornecedor {
	return &ProdutosEmFornecedor{
		ID:          prdFornec_request.ID,
		DataType:    "produto_fornecedor",
		IDCategoria: prdFornec_request.IDCategoria,
		Nome:        validation.CareString(prdFornec_request.Nome),
		Enabled:     true,
		CreatedAt:   time.Now().String(),
	}
}

func NewProdutoEmCategoria(prdCat_request ProdutosEmCategorias) *ProdutosEmCategorias {
	return &ProdutosEmCategorias{
		ID:      primitive.NewObjectID(),
		Nome:    validation.CareString(prdCat_request.Nome),
		Enabled: true,
	}
}
