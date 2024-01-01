package produto

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

type ProdutoServiceInterface interface {
	CreateProduto(ctx context.Context, produto model.Produto) (*model.Produto, error)
	UpdateProduto(ctx context.Context, ID string, produto *model.Produto) (bool, error)
	GetProdutoByID(ctx context.Context, ID string) (*model.Produto, error)
	GetAllProdutos(ctx context.Context, filters model.FilterProduto, limit, page int64) (*model.Paginate, error)
}

type ProdutoDataService struct {
	mdb mongodb.MongoDBInterface
}

func NewProdutoervice(mongo_connection mongodb.MongoDBInterface) *ProdutoDataService {
	return &ProdutoDataService{
		mdb: mongo_connection,
	}
}

func (p *ProdutoDataService) CreateProduto(ctx context.Context, Produto model.Produto) (*model.Produto, error) {

	collection := p.mdb.GetCollection("produtos")

	dt := time.Now().Format(time.RFC3339)

	Produto.Enabled = true
	Produto.CreatedAt = dt
	Produto.UpdatedAt = dt
	Produto.ID = primitive.NewObjectID()

	result, err := collection.InsertOne(ctx, Produto)
	if err != nil {
		logger.Error("erro salvar  Produto", err)
		return &Produto, err
	}

	Produto.ID = result.InsertedID.(primitive.ObjectID)

	return &Produto, nil
}

func (p *ProdutoDataService) UpdateProduto(ctx context.Context, ID string, Produto *model.Produto) (bool, error) {
	collection := p.mdb.GetCollection("produtos")

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
			{Key: "nome", Value: Produto.Nome},
			{Key: "enabled", Value: Produto.Enabled},
			{Key: "categoria", Value: Produto.Categoria},
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

func (p *ProdutoDataService) GetProdutoByID(ctx context.Context, ID string) (*model.Produto, error) {

	collection := p.mdb.GetCollection("produtos")

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
		logger.Error("erro ao consultar Produto", err)
		return nil, err
	}

	return Produto, nil
}

func (p *ProdutoDataService) GetAllProdutos(ctx context.Context, filters model.FilterProduto, limit, page int64) (*model.Paginate, error) {
	collection := p.mdb.GetCollection("produtos")

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
		p := &model.Produto{}
		if err := curr.Decode(p); err != nil {
			logger.Error("erro ao consulta todas as Produtos", err)
		}
		result = append(result, p)
	}

	pagination.Paginate(result)

	return pagination, nil
}
