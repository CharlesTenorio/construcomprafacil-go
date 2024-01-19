package model

type Endereco struct {
	DataType    string      `bson:"data_type" json:"-"`
	CEP         string      `bson:"cep" json:"cep"`
	Rua         string      `bson:"rua" json:"rua"`
	Bairro      string      `bson:"bairro" json:"bairro"`
	Cidade      string      `bson:"cidade" json:"cidade"`
	Estado      string      `bson:"uf" json:"uf"`
	Coordenadas Localizacao `bson:"coordenadas" json:"coordenadas"`
}

type Localizacao struct {
	Latitude  float64 `bson:"latitude" json:"latitude"`
	Longitude float64 `bson:"longitude" json:"longitude"`
}
