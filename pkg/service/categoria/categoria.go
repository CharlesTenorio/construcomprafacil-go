package categoria

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

type CategoriaServiceInterface interface {
	Create(ctx context.Context, categoria model.Categoria) (*model.Categoria, error)
	Update(ctx context.Context, ID string, meioToChange *model.Categoria) (bool, error)
	GetByID(ctx context.Context, ID string) (*model.Categoria, error)
	GetAll(ctx context.Context, filters model.FilterCategoria, limit, page int64) (*model.Paginate, error)
	ListPrd(ctx context.Context, ID string, limit, page int64) (*model.Paginate, error)
	ListSubcategoria(ctx context.Context, ID string, limit, page int64) (*model.Paginate, error)
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
	collection := cat.mdb.GetCollection("categorias")

	dt := time.Now().Format(time.RFC3339)

	categoria.Enabled = true
	categoria.CreatedAt = dt
	categoria.UpdatedAt = dt

	result, err := collection.InsertOne(ctx, categoria)
	if err != nil {
		logger.Error("erro salvar  categoria", err)
		return &categoria, err
	}

	categoria.ID = result.InsertedID.(primitive.ObjectID)

	return &categoria, nil
}

func (cat *CategoriaDataService) Update(ctx context.Context, ID string, categoria *model.Categoria) (bool, error) {
	collection := cat.mdb.GetCollection("categorias")

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

	collection := cat.mdb.GetCollection("categorias")

	categoria := &model.Categoria{}

	objectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {

		logger.Error("Error to parse ObjectIDFromHex", err)
		return nil, err
	}

	filter := bson.D{
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
	collection := cat.mdb.GetCollection("categorias")

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

func (cat *CategoriaDataService) ListPrd(ctx context.Context, categoriaID string, limit, page int64) (*model.Paginate, error) {
	collection := cat.mdb.GetCollection("produtos")

	categoriaObjectID, err := primitive.ObjectIDFromHex(categoriaID)
	if err != nil {
		logger.Error("Error parsing ObjectIDFromHex for categoria", err)
		return nil, err
	}

	// Consulta produtos na categoria especificada
	query := bson.M{"categoria._id": categoriaObjectID}

	curr, err := collection.Find(ctx, query)
	if err != nil {
		logger.Error("Error while querying produtos", err)
		return nil, err
	}
	defer curr.Close(ctx)

	/*var produtos []*model.Produto
	for curr.Next(ctx) {
		produto := &model.Produto{}
		if err := curr.Decode(produto); err != nil {
			logger.Error("Error decoding produto", err)
			return nil, err
		}
		produtos = append(produtos, produto)
	}

	return produtos, nil*/

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

func (cat *CategoriaDataService) ListSubcategoria(ctx context.Context, categoriaID string, limit, page int64) (*model.Paginate, error) {
	collection := cat.mdb.GetCollection("categorias")

	categoriaObjectID, err := primitive.ObjectIDFromHex(categoriaID)
	if err != nil {
		logger.Error("Error parsing ObjectIDFromHex for categoria", err)
		return nil, err
	}

	// Consulta a categoria especificada
	filter := bson.D{{Key: "_id", Value: categoriaObjectID}}
	projection := bson.D{{Key: "subcategorias", Value: 1}}

	var categoria model.Categoria
	err = collection.FindOne(ctx, filter, options.FindOne().SetProjection(projection)).Decode(&categoria)
	if err != nil {
		logger.Error("Error while querying categoria", err)
		return nil, err
	}

	// Filtra subcategorias com campo Enabled igual a true
	subcategorias := make([]model.Subcategoria, 0)
	for _, subcategoria := range categoria.Subcategorias {
		if subcategoria.Enabled {
			// Remova os campos que você não deseja retornar
			subcategoria.Enabled = false
			subcategoria.CreatedAt = ""
			subcategoria.UpdatedAt = ""
			subcategorias = append(subcategorias, subcategoria)
		}
	}

	// Paginação
	count := int64(len(subcategorias))
	pagination := model.NewPaginate(limit, page, count)

	pagination.Paginate(subcategorias)

	return pagination, nil
}
