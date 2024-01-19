package produto

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

type ProdutoServiceInterface interface {
	Create(ctx context.Context, Produto model.Produto) (*model.Produto, error)
	Update(ctx context.Context, ID string, subToChange *model.Produto) (bool, error)
	GetByID(ctx context.Context, ID string) (*model.Produto, error)
	GetAll(ctx context.Context, filters model.FilterProduto, limit, page int64) (*model.Paginate, error)
	AddFornecedroes(ctx context.Context, ID string, fornec *[]dto.FornecedoresEmPrd) (bool, error)
}

type ProdutoDataService struct {
	mdb mongodb.MongoDBInterface
}

func NewProdutoervice(mongo_connection mongodb.MongoDBInterface) *ProdutoDataService {
	return &ProdutoDataService{
		mdb: mongo_connection,
	}
}

func (prd *ProdutoDataService) Create(ctx context.Context, produto model.Produto) (*model.Produto, error) {
	collection := prd.mdb.GetCollection("cfStore")
	prod := model.NewProduto(produto)

	result, err := collection.InsertOne(ctx, prod)
	if err != nil {
		logger.Error("erro salvar  Produto", err)
		return &produto, err
	}

	produto.ID = result.InsertedID.(primitive.ObjectID)

	return &produto, nil
}

func (prd *ProdutoDataService) Update(ctx context.Context, ID string, Produto *model.Produto) (bool, error) {
	collection := prd.mdb.GetCollection("cfStore")

	opts := options.Update().SetUpsert(true)

	objectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {

		logger.Error("Error to parse ObjectIDFromHex", err)
		return false, err
	}

	filter := bson.D{
		{Key: "data_type", Value: "produto"},
		{Key: "_id", Value: objectID},
	}

	update := bson.D{{Key: "$set",
		Value: bson.D{
			{Key: "nome", Value: Produto.Nome},
			{Key: "enabled", Value: Produto.Enabled},
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

func (prd *ProdutoDataService) GetByID(ctx context.Context, ID string) (*model.Produto, error) {

	collection := prd.mdb.GetCollection("cfStore")

	Produto := &model.Produto{}

	objectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {

		logger.Error("Error to parse ObjectIDFromHex", err)
		return nil, err
	}

	filter := bson.D{
		{Key: "_id", Value: objectID},
	}

	err = collection.FindOne(ctx, filter).Decode(Produto)
	if err != nil {
		logger.Error("erro ao consultar produto", err)
		return nil, err
	}

	return Produto, nil
}

func (prd *ProdutoDataService) GetAll(ctx context.Context, filters model.FilterProduto, limit, page int64) (*model.Paginate, error) {
	collection := prd.mdb.GetCollection("cfStore")

	query := bson.M{"data_type": "produto"}

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
		logger.Error("erro ao consultar todas as Produtos", err)
		return nil, err
	}

	pagination := model.NewPaginate(limit, page, count)

	curr, err := collection.Find(ctx, query, pagination.GetPaginatedOpts())
	if err != nil {
		return nil, err
	}

	result := make([]*model.Produto, 0)
	for curr.Next(ctx) {
		prd := &model.Produto{}
		if err := curr.Decode(prd); err != nil {
			logger.Error("erro ao consulta todas as Produtos", err)
		}
		result = append(result, prd)
	}

	pagination.Paginate(result)

	return pagination, nil
}

func (prd *ProdutoDataService) AddFornecedroes(ctx context.Context, ID string, fornecedores *[]dto.FornecedoresEmPrd) (bool, error) {
	collection := prd.mdb.GetCollection("cfStore")

	produtoID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		logger.Error("Erro ao converter ID para ObjectID", err)
		return false, err
	}

	update := bson.D{{Key: "$push",
		Value: bson.D{
			{Key: "fornecedores", Value: *fornecedores},
		},
	}}

	filter := bson.D{
		{Key: "_id", Value: produtoID},
	}

	_, err = collection.UpdateOne(ctx, filter, update)
	if err != nil {
		logger.Error("Erro ao adicionar fornecedores ao produto", err)
		return false, err
	}

	return true, nil
}
