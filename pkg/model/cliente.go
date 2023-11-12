package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Cliente struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	TipoDado     string             `bson:"data_type" json:"tipo_cliente"`
	Nome         string             `bson:"name" json:"nome"`
	CPF          string             `bson:"cnpj" json:"cnpj"`
	Endereco     string             `bson:"address" json:"endereco"`
	Complemento  string             `bson:"complement" json:"complemento"`
	Cidade       string             `bson:"city" json:"cidade"`
	Estado       string             `bson:"state" json:"estado"`
	CEP          string             `bson:"cep" json:"cep"`
	Email        string             `bson:"email" json:"email"`
	Ativo        bool               `bson:"enabled" json:"ativo"`
	CriadoEm     string             `bson:"created_at" json:"criado_em,omitempty"`
	AtualizadoEm string             `bson:"updated_at" json:"atualizado_em,omitempty"`
}

type EnderecoEntrega struct {
	CEP         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	UF          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	DDD         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

type FiltroCliente struct {
	Nome  string `json:"nome"`
	Ativo string `json:"ativo"`
}
