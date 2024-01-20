package catalagoprd

import (
	"context"
	"encoding/json"

	"github.com/katana/back-end/orcafacil-go/internal/config/logger"
	"github.com/katana/back-end/orcafacil-go/pkg/adapter/mongodb"
	rd "github.com/katana/back-end/orcafacil-go/pkg/adapter/redis"
	"github.com/katana/back-end/orcafacil-go/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CatalagoServiceInterface interface {
	SalveRedis(ctx context.Context) (bool, error)
	LerRedis(ctx context.Context, limit, page int64) (*model.Paginate, error)
	ListProduto(ctx context.Context) (*[]model.Categoria, error)
}

type CatalagoDataService struct {
	mdb mongodb.MongoDBInterface
	rds rd.RedisClientInterface
}

func NewCatalagoService(mongo_connection mongodb.MongoDBInterface, redis_connection rd.RedisClientInterface) *CatalagoDataService {
	return &CatalagoDataService{
		mdb: mongo_connection,
		rds: redis_connection,
	}
}

func (cat *CatalagoDataService) SalveRedis(ctx context.Context) (bool, error) {

	categorias, err := cat.ListProduto(ctx)
	if err != nil {
		logger.Error("Erro ao listar as  categorias", err)
		return false, err
	}
	categoriasBytes, err := json.Marshal(categorias)
	if err != nil {
		logger.Error("Erro ao converter categorias para bytes", err)
		return false, err
	}
	cat.rds.SaveData(ctx, "catalogo", categoriasBytes, 0)

	return true, nil
}

func (cat *CatalagoDataService) LerRedis(ctx context.Context, limit, page int64) (*model.Paginate, error) {
	// Lendo dados do Redis
	categoriasBytes, err := cat.rds.ReadData(ctx, "catalogo")
	if err != nil {
		logger.Error("Erro ao ler dados do Redis", err)
		return nil, err
	}

	// Convertendo bytes para estrutura de dados
	var categorias []model.Categoria
	if err := json.Unmarshal(categoriasBytes, &categorias); err != nil {
		logger.Error("Erro ao decodificar dados do Redis", err)
		return nil, err
	}

	// Aplicando paginação
	pagination := model.NewPaginate(limit, page, int64(len(categorias)))
	pagination.Paginate(categorias)

	return pagination, nil
}

func (cat *CatalagoDataService) ListProduto(ctx context.Context) ([]model.Categoria, error) {
	collection := cat.mdb.GetCollection("cfStore")

	// Consulta todas as categorias
	query := bson.M{
		"data_type": "categoria",
	}
	projection := bson.D{
		{Key: "_id", Value: 1},
		{Key: "nome", Value: 1},
		{Key: "enabled", Value: 1},
		{Key: "produtos", Value: 1}, // Adicionando o campo "produtos" à projeção
	}

	var categorias []model.Categoria
	cursor, err := collection.Find(ctx, query, options.Find().SetProjection(projection).SetSort(bson.D{{Key: "nome", Value: 1}}))
	if err != nil {
		logger.Error("Erro ao consultar categorias:", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	// Iterar sobre o cursor e decodificar as categorias
	for cursor.Next(ctx) {
		var categoria model.Categoria
		if err := cursor.Decode(&categoria); err != nil {
			logger.Error("Erro ao decodificar categoria:", err)
			return nil, err
		}
		categorias = append(categorias, categoria)
	}

	if err := cursor.Err(); err != nil {
		logger.Error("Erro durante a iteração do cursor:", err)
		return nil, err
	}

	return categorias, nil
}
