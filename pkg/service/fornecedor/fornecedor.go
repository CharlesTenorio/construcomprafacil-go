package fornecedor

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/katana/back-end/orcafacil-go/internal/config/logger"
	"github.com/katana/back-end/orcafacil-go/internal/dto"
	"github.com/katana/back-end/orcafacil-go/pkg/adapter/mongodb"
	"github.com/katana/back-end/orcafacil-go/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type FornecedorServiceInterface interface {
	Create(ctx context.Context, Fornecedor model.Fornecedor) (*model.Fornecedor, error)
	Update(ctx context.Context, ID string, meioToChange *model.Fornecedor) (bool, error)
	GetByID(ctx context.Context, ID string) (*model.Fornecedor, error)
	GetAll(ctx context.Context, filters model.FilterFornecedor, limit, page int64) (*model.Paginate, error)
	GetByCnpj(ctx context.Context, Cnpj string) bool
	AddProdutos(ctx context.Context, ID string, prds []dto.ProdutosEmFornecedor) (bool, error)
	UpdFornecedorParaPrd(ctx context.Context, idPrd string, produto *model.Produto) (bool, error)
}

type FornecedorDataService struct {
	mdb mongodb.MongoDBInterface
}

func NewFornecedorervice(mongo_connection mongodb.MongoDBInterface) *FornecedorDataService {
	return &FornecedorDataService{
		mdb: mongo_connection,
	}
}

func (fornec *FornecedorDataService) Create(ctx context.Context, Fornecedor model.Fornecedor) (*model.Fornecedor, error) {
	collection := fornec.mdb.GetCollection("fornecedores")

	dt := time.Now().Format(time.RFC3339)

	Fornecedor.Enabled = true
	Fornecedor.CreatedAt = dt
	Fornecedor.UpdatedAt = dt
	Fornecedor.ID = primitive.NewObjectID()

	result, err := collection.InsertOne(ctx, Fornecedor)
	if err != nil {
		logger.Error("erro salvar  Fornecedor", err)
		return &Fornecedor, err
	}

	Fornecedor.ID = result.InsertedID.(primitive.ObjectID)

	return &Fornecedor, nil
}

func (fornec *FornecedorDataService) Update(ctx context.Context, ID string, Fornecedor *model.Fornecedor) (bool, error) {
	collection := fornec.mdb.GetCollection("fornecedores")

	opts := options.Update().SetUpsert(true)

	objectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {

		logger.Error("Error to parse ObjectIDFromHex", err)
		return false, err
	}

	filter := bson.D{

		{Key: "_id", Value: objectID},
	}

	update := bson.D{{Key: "$set",
		Value: bson.D{
			{Key: "nome", Value: Fornecedor.Nome},
			{Key: "enabled", Value: Fornecedor.Enabled},
			{Key: "updated_at", Value: time.Now().Format(time.RFC3339)},
		},
	}}

	_, err = collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		logger.Error("Error while updating data", err)

		return false, err
	}

	return true, nil
}

func (fornec *FornecedorDataService) GetByID(ctx context.Context, ID string) (*model.Fornecedor, error) {

	collection := fornec.mdb.GetCollection("fornecedores")

	Fornecedor := &model.Fornecedor{}

	objectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {

		logger.Error("Error to parse ObjectIDFromHex", err)
		return nil, err
	}

	filter := bson.D{
		{Key: "_id", Value: objectID},
	}

	err = collection.FindOne(ctx, filter).Decode(Fornecedor)
	if err != nil {
		logger.Error("erro ao consultar Fornecedor", err)
		return nil, err
	}

	return Fornecedor, nil
}

func (fornec *FornecedorDataService) GetAll(ctx context.Context, filters model.FilterFornecedor, limit, page int64) (*model.Paginate, error) {
	collection := fornec.mdb.GetCollection("fornecedores")

	query := bson.M{}

	if filters.Nome != "" || filters.Enabled != "" {
		if filters.Nome != "" {
			query["nome"] = bson.M{"$regex": fmt.Sprintf(".*%s.*", filters.Nome), "$options": "i"}
		}
		if filters.Enabled != "" {
			enable, err := strconv.ParseBool(filters.Enabled)
			if err != nil {
				logger.Error("erro converter campo enabled", err)
				return nil, err
			}
			query["enabled"] = enable
		}
	}
	count, err := collection.CountDocuments(ctx, query, &options.CountOptions{})

	if err != nil {
		logger.Error("erro ao consultar todas as Fornecedors", err)
		return nil, err
	}

	pagination := model.NewPaginate(limit, page, count)

	curr, err := collection.Find(ctx, query, pagination.GetPaginatedOpts())
	if err != nil {
		return nil, err
	}

	result := make([]*model.Fornecedor, 0)
	for curr.Next(ctx) {
		fornec := &model.Fornecedor{}
		if err := curr.Decode(fornec); err != nil {
			logger.Error("erro ao consulta todas as Fornecedors", err)
		}
		result = append(result, fornec)
	}

	pagination.Paginate(result)

	return pagination, nil
}

func (fornec *FornecedorDataService) GetByCnpj(ctx context.Context, Cnpj string) bool {

	collection := fornec.mdb.GetCollection("fornecedores")

	// Utilizando o método CountDocuments para verificar a existência
	filter := bson.D{{Key: "cnpj", Value: Cnpj}}
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		logger.Error("erro ao consultar Fornecedor", err)
		return false
	}

	// Se count for maior que zero, o fornecedor existe
	return count > 0
}

func (fornec *FornecedorDataService) AddProdutos(ctx context.Context, ID string, prds []dto.ProdutosEmFornecedor) (bool, error) {
	collection := fornec.mdb.GetCollection("fornecedores")

	fornecedorID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		logger.Error("Erro ao converter ID para ObjectID", err)
		return false, err
	}
	// push
	update := bson.D{{Key: "$addToSet",
		Value: bson.D{
			{Key: "produtos", Value: bson.D{
				{Key: "$each", Value: prds},
			}},
		},
	}}

	filter := bson.D{
		{Key: "_id", Value: fornecedorID},
	}

	// Atualize os produtos do fornecedor
	_, err = collection.UpdateOne(ctx, filter, update)
	if err != nil {
		logger.Error("Erro ao adicionar produtos ao fornecedor", err)
		return false, err
	}

	// Atualize a lista de fornecedores associados a cada produto
	for _, prd := range prds {
		produto := &model.Produto{
			Fornecedores: []dto.FornecedoresEmPrd{{ID: ID}}, // Adicione o ID do fornecedor à lista de fornecedores do produto
		}

		_, err := fornec.UpdFornecedorParaPrd(ctx, prd.ID, produto)
		if err != nil {
			logger.Error("Erro ao atualizar lista de fornecedores do produto", err)
			return false, err
		}
	}

	return true, nil
}

func (fornec *FornecedorDataService) UpdFornecedorParaPrd(ctx context.Context, idPrd string, produto *model.Produto) (bool, error) {
	collection := fornec.mdb.GetCollection("produtos")

	opts := options.Update().SetUpsert(true)

	objectID, err := primitive.ObjectIDFromHex(idPrd)
	if err != nil {

		logger.Error("Error to parse ObjectIDFromHex", err)
		return false, err
	}

	filter := bson.D{

		{Key: "_id", Value: objectID},
	}

	update := bson.D{{Key: "$set",
		Value: bson.D{

			{Key: "fornecedores", Value: produto.Fornecedores},
			{Key: "updated_at", Value: time.Now().Format(time.RFC3339)},
		},
	}}

	_, err = collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		logger.Error("Erro ao atuilziar Produto", err)

		return false, err
	}

	return true, nil

}
