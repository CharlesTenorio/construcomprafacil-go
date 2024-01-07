package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

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
	ID         string  `json:"id"`
	Descricao  string  `json:"descricao"`
	PrecoVenda float64 `json:"preco_venda"`
	Enabled    bool    `json:"enabled"`
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
