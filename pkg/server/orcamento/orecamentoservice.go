package orcamento

import (
	"context"

	"github.com/katana/back-end/orcafacil-go/pkg/adapter/mongodb"
	"github.com/katana/back-end/orcafacil-go/pkg/adapter/rabbitmq"
	"github.com/katana/back-end/orcafacil-go/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrcamentoServiceInterface interface {
	GetAll(ctx context.Context, filters model.FilterOrcamento, limit, page int64) (*model.Paginate, error)
	GetOrcamentoByID(id primitive.ObjectID) (*model.Orcamento, error)
	CountPeidoPorClient(cliente_id, start_date, end_date string) ([]bson.M, error)
	GetByID(ctx context.Context, ID string) (*model.User, error)
	Create(ctx context.Context, Orcamento *model.Orcamento) (map[string]interface{}, error)
	Update(ctx context.Context, ID string, peidodoToChange *model.Orcamento) (bool, error)
	Delete(ctx context.Context, ID string) (bool, error)
}

type OrcamentoDataService struct {
	mdb       mongodb.MongoDBInterface
	rabbit_mq rabbitmq.RabbitInterface
}

func NewPeididoService(rabbit_connection rabbitmq.RabbitInterface, mongo_connection mongodb.MongoDBInterface) *OrcamentoDataService {
	return &OrcamentoDataService{
		rabbit_mq: rabbit_connection,
		mdb:       mongo_connection,
	}
}
