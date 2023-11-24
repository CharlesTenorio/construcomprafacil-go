package pedido

import (
	"context"

	"github.com/katana/back-end/orcafacil-go/pkg/adapter/mongodb"
	"github.com/katana/back-end/orcafacil-go/pkg/adapter/rabbitmq"
	"github.com/katana/back-end/orcafacil-go/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PeididoServiceInterface interface {
	GetAll(ctx context.Context, filters model.FilterPedido, limit, page int64) (*model.Paginate, error)
	GetPedidoByID(id primitive.ObjectID) (*model.Pedido, error)
	CountPeidoPorClient(cliente_id, start_date, end_date string) ([]bson.M, error)
	GetByID(ctx context.Context, ID string) (*model.User, error)
	Create(ctx context.Context, pedido *model.Pedido) (*model.Pedido, error)
	Update(ctx context.Context, ID string, peidodoToChange *model.Pedido) (bool, error)
	Delete(ctx context.Context, ID string) (bool, error)
}

type PeididoDataService struct {
	mdb       mongodb.MongoDBInterface
	rabbit_mq rabbitmq.RabbitInterface
}

func NewPeididoService(rabbit_connection rabbitmq.RabbitInterface, mongo_connection mongodb.MongoDBInterface) *PeididoDataService {
	return &PeididoDataService{
		rabbit_mq: rabbit_connection,
		mdb:       mongo_connection,
	}
}
