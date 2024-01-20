package categoria

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/katana/back-end/orcafacil-go/internal/config/logger"
	"github.com/katana/back-end/orcafacil-go/internal/dto"
	"github.com/katana/back-end/orcafacil-go/pkg/adapter/mongodb"
	"github.com/katana/back-end/orcafacil-go/pkg/model"
	"github.com/katana/back-end/orcafacil-go/pkg/service/validation"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CategoriaServiceInterface interface {
	Create(ctx context.Context, categoria model.Categoria) (*model.Categoria, error)
	Update(ctx context.Context, ID string, categoriaToChange *model.Categoria) (bool, error)
	GetByID(ctx context.Context, ID string) (*model.Categoria, error)
	GetAll(ctx context.Context, filters model.FilterCategoria, limit, page int64) (*model.Paginate, error)
	ListProduto(ctx context.Context, ID string, limit, page int64) (*model.Paginate, error)
}

type CategoriaDataService struct {
	mdb mongodb.MongoDBInterface
}

func NewCategoriaervice(mongo_connection mongodb.MongoDBInterface) *CategoriaDataService {
	return &CategoriaDataService{
		mdb: mongo_connection,
	}
}

func (cat *CategoriaDataService) Create(ctx context.Context, categoria model.Categoria) (*model.Categoria, error) {
	collection := cat.mdb.GetCollection("cfStore")
	cate := model.NewCategoria(categoria)

	for i := range cate.Produtos {
		cate.Produtos[i].ID = primitive.NewObjectID()
		cate.Produtos[i].Enabled = true
	}

	result, err := collection.InsertOne(ctx, cate)
	if err != nil {
		logger.Error("erro salvar  categoria", err)
		return &categoria, err
	}

	categoria.ID = result.InsertedID.(primitive.ObjectID)

	return &categoria, nil
}

func (cat *CategoriaDataService) Update(ctx context.Context, ID string, categoria *model.Categoria) (bool, error) {
	collection := cat.mdb.GetCollection("cfStore")

	opts := options.Update().SetUpsert(true)

	objectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {

		logger.Error("Error to parse ObjectIDFromHex", err)
		return false, err
	}

	filter := bson.D{
		{Key: "data_type", Value: "categoria"},
		{Key: "_id", Value: objectID},
	}
	categoria.Nome = validation.CareString(categoria.Nome)
	update := bson.D{{Key: "$set",
		Value: bson.D{
			{Key: "nome", Value: categoria.Nome},
			{Key: "enabled", Value: categoria.Enabled},
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

func (cat *CategoriaDataService) GetByID(ctx context.Context, ID string) (*model.Categoria, error) {

	collection := cat.mdb.GetCollection("cfStore")

	categoria := &model.Categoria{}

	objectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {

		logger.Error("Error to parse ObjectIDFromHex", err)
		return nil, err
	}

	filter := bson.D{
		{Key: "data_type", Value: "categoria"},
		{Key: "_id", Value: objectID},
	}

	err = collection.FindOne(ctx, filter).Decode(categoria)
	if err != nil {
		logger.Error("erro ao consultar categoria", err)
		return nil, err
	}

	return categoria, nil
}

func (cat *CategoriaDataService) GetAll(ctx context.Context, filters model.FilterCategoria, limit, page int64) (*model.Paginate, error) {
	collection := cat.mdb.GetCollection("cfStore")

	query := bson.M{
		"data_type": "categoria",
	}

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
		logger.Error("erro ao consultar todas as categorias", err)
		return nil, err
	}

	pagination := model.NewPaginate(limit, page, count)

	curr, err := collection.Find(ctx, query, pagination.GetPaginatedOpts())
	if err != nil {
		return nil, err
	}

	result := make([]*model.Categoria, 0)
	for curr.Next(ctx) {
		cat := &model.Categoria{}
		if err := curr.Decode(cat); err != nil {
			logger.Error("erro ao consulta todas as categorias", err)
		}
		result = append(result, cat)
	}

	pagination.Paginate(result)

	return pagination, nil
}

func (cat *CategoriaDataService) ListProduto(ctx context.Context, categoriaID string, limit, page int64) (*model.Paginate, error) {
	collection := cat.mdb.GetCollection("cfStore")

	categoriaObjectID, err := primitive.ObjectIDFromHex(categoriaID)
	if err != nil {
		logger.Error("Error parsing ObjectIDFromHex for categoria", err)
		return nil, err
	}

	// Consulta a categoria especificada
	filter := bson.D{
		{Key: "data_type", Value: "categoria"},
		{Key: "_id", Value: categoriaObjectID}}
	projection := bson.D{
		{Key: "Produtos", Value: 1},
	}

	var categoria model.Categoria
	err = collection.FindOne(ctx, filter, options.FindOne().SetProjection(projection)).Decode(&categoria)
	if err != nil {
		logger.Error("Error while querying categoria", err)
		return nil, err
	}

	// Filtra Produtos com campo Enabled igual a true
	Produtos := make([]dto.ProdutosEmCategorias, 0)
	for _, Produto := range categoria.Produtos {
		if Produto.Enabled {
			// Remova os campos que você não deseja retornar
			Produto.Enabled = false

			Produtos = append(Produtos, Produto)
		}
	}

	// Paginação
	count := int64(len(Produtos))
	pagination := model.NewPaginate(limit, page, count)

	pagination.Paginate(Produtos)

	return pagination, nil
}
