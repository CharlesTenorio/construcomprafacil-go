package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type ProdutoForcedor struct {
	FornecedorID    primitive.ObjectID `bson:"fornecedor_id" json:"fornecedor_id"`
	CodigoPrdNaLoja string             `bson:"codigo_produto_na_loja" json:"codigo_produto_na_loja"`
	ProdutoID       primitive.ObjectID `bson:"produto_id" json:"produto_id"`
	UnidadeMedida   string             `bson:"unidade_medida" json:"unidade_medida"`
	Descricao       string             `bson:"Descricao" json:"Descricao"`
}
