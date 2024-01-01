package subSubcategoria

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

type SubcategoriaServiceInterface interface {
	Create(ctx context.Context, subcategoria model.Subcategoria) (*model.Subcategoria, error)
	Update(ctx context.Context, ID string, subToChange *model.Subcategoria) (bool, error)
	GetByID(ctx context.Context, ID string) (*model.Subcategoria, error)
	GetAll(ctx context.Context, filters model.FilterSubcategoria, limit, page int64) (*model.Paginate, error)
	ListPrd(ctx context.Context, ID string, limit, page int64) (*model.Paginate, error)
}

type SubcategoriaDataService struct {
	mdb mongodb.MongoDBInterface
}

func NewSubcategoriaervice(mongo_connection mongodb.MongoDBInterface) *SubcategoriaDataService {
	return &SubcategoriaDataService{
		mdb: mongo_connection,
	}
}

func (cat *SubcategoriaDataService) Create(ctx context.Context, subcategoria model.Subcategoria) (*model.Subcategoria, error) {
	collection := cat.mdb.GetCollection("subcategorias")

	dt := time.Now().Format(time.RFC3339)

	subcategoria.Enabled = true
	subcategoria.CreatedAt = dt
	subcategoria.UpdatedAt = dt

	result, err := collection.InsertOne(ctx, subcategoria)
	if err != nil {
		logger.Error("erro salvar  Subcategoria", err)
		return &subcategoria, err
	}

	subcategoria.ID = result.InsertedID.(primitive.ObjectID)

	return &subcategoria, nil
}

func (cat *SubcategoriaDataService) Update(ctx context.Context, ID string, subcategoria *model.Subcategoria) (bool, error) {
	collection := cat.mdb.GetCollection("subcategorias")

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
			{Key: "nome", Value: subcategoria.Nome},
			{Key: "enabled", Value: subcategoria.Enabled},
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

func (cat *SubcategoriaDataService) GetByID(ctx context.Context, ID string) (*model.Subcategoria, error) {

	collection := cat.mdb.GetCollection("subcategorias")

	subcategoria := &model.Subcategoria{}

	objectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {

		logger.Error("Error to parse ObjectIDFromHex", err)
		return nil, err
	}

	filter := bson.D{
		{Key: "_id", Value: objectID},
	}

	err = collection.FindOne(ctx, filter).Decode(subcategoria)
	if err != nil {
		logger.Error("erro ao consultar Subcategoria", err)
		return nil, err
	}

	return subcategoria, nil
}

func (cat *SubcategoriaDataService) GetAll(ctx context.Context, filters model.FilterSubcategoria, limit, page int64) (*model.Paginate, error) {
	collection := cat.mdb.GetCollection("subcategorias")

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
		logger.Error("erro ao consultar todas as Subcategorias", err)
		return nil, err
	}

	pagination := model.NewPaginate(limit, page, count)

	curr, err := collection.Find(ctx, query, pagination.GetPaginatedOpts())
	if err != nil {
		return nil, err
	}

	result := make([]*model.Subcategoria, 0)
	for curr.Next(ctx) {
		cat := &model.Subcategoria{}
		if err := curr.Decode(cat); err != nil {
			logger.Error("erro ao consulta todas as Subcategorias", err)
		}
		result = append(result, cat)
	}

	pagination.Paginate(result)

	return pagination, nil
}

func (cat *SubcategoriaDataService) ListPrd(ctx context.Context, subcategoriaID string, limit, page int64) (*model.Paginate, error) {
	collection := cat.mdb.GetCollection("produtos")

	subcategoriaObjectID, err := primitive.ObjectIDFromHex(subcategoriaID)
	if err != nil {
		logger.Error("Error parsing ObjectIDFromHex for Subcategoria", err)
		return nil, err
	}

	// Consulta produtos na Subcategoria especificada
	query := bson.M{"subcategoria._id": subcategoriaObjectID}

	curr, err := collection.Find(ctx, query)
	if err != nil {
		logger.Error("Error while querying produtos", err)
		return nil, err
	}
	defer curr.Close(ctx)

	count, err := collection.CountDocuments(ctx, query, &options.CountOptions{})

	if err != nil {
		logger.Error("erro ao consultar todas as Produtos", err)
		return nil, err
	}

	pagination := model.NewPaginate(limit, page, count)

	curr, err = collection.Find(ctx, query, pagination.GetPaginatedOpts())
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
