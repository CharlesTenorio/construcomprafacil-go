package fornecedor

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/katana/back-end/orcafacil-go/internal/config/logger"
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
	ListPrdFornecedor(ctx context.Context, ID string, limit, page int64) (*model.Paginate, error)
}

type FornecedorDataService struct {
	mdb mongodb.MongoDBInterface
}

func NewFornecedorervice(mongo_connection mongodb.MongoDBInterface) *FornecedorDataService {
	return &FornecedorDataService{
		mdb: mongo_connection,
	}
}

func (cat *FornecedorDataService) Create(ctx context.Context, Fornecedor model.Fornecedor) (*model.Fornecedor, error) {
	collection := cat.mdb.GetCollection("fornecedores")

	dt := time.Now().Format(time.RFC3339)

	Fornecedor.Enabled = true
	Fornecedor.CreatedAt = dt
	Fornecedor.UpdatedAt = dt

	result, err := collection.InsertOne(ctx, Fornecedor)
	if err != nil {
		logger.Error("erro salvar  Fornecedor", err)
		return &Fornecedor, err
	}

	Fornecedor.ID = result.InsertedID.(primitive.ObjectID)

	return &Fornecedor, nil
}

func (cat *FornecedorDataService) Update(ctx context.Context, ID string, Fornecedor *model.Fornecedor) (bool, error) {
	collection := cat.mdb.GetCollection("fornecedores")

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

func (cat *FornecedorDataService) GetByID(ctx context.Context, ID string) (*model.Fornecedor, error) {

	collection := cat.mdb.GetCollection("fornecedores")

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

func (cat *FornecedorDataService) GetAll(ctx context.Context, filters model.FilterFornecedor, limit, page int64) (*model.Paginate, error) {
	collection := cat.mdb.GetCollection("fornecedores")

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
		cat := &model.Fornecedor{}
		if err := curr.Decode(cat); err != nil {
			logger.Error("erro ao consulta todas as Fornecedors", err)
		}
		result = append(result, cat)
	}

	pagination.Paginate(result)

	return pagination, nil
}

func (cat *FornecedorDataService) ListPrdFornecedor(ctx context.Context, fornecedorID string, limit, page int64) (*model.Paginate, error) {

	collection := cat.mdb.GetCollection("fornecedores")

	objectID, err := primitive.ObjectIDFromHex(fornecedorID)
	if err != nil {
		return nil, err
	}

	filter := bson.D{{Key: "_id", Value: objectID}}

	// Executa a consulta
	var fornecedor model.Fornecedor
	err = collection.FindOne(ctx, filter).Decode(&fornecedor)
	if err != nil {
		return nil, err
	}

	// Filtra os produtos habilitados
	var produtosHabilitados []model.Produto
	for _, produto := range fornecedor.Produto {
		if produto.Enabled {
			produtosHabilitados = append(produtosHabilitados, produto)
		}
	}

	// Paginação dos produtos habilitados associados ao fornecedor
	paginate := model.NewPaginate(limit, page, int64(len(produtosHabilitados)))

	// Calcula os índices de início e fim para a fatia de produtos
	start := (paginate.Page - 1) * paginate.Limit
	end := start + paginate.Limit

	// Evita índices fora do alcance
	if start > paginate.Total {
		start = paginate.Total
	}
	if end > paginate.Total {
		end = paginate.Total
	}

	// Obtem a fatia de produtos para a página atual
	produtosPaginados := produtosHabilitados[start:end]

	// Atualiza a estrutura de Paginate com os dados paginados
	paginate.Paginate(produtosPaginados)

	return paginate, nil
}
