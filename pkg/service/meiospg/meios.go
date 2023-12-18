package meiospg

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

type MeiosServiceInterface interface {
	Create(ctx context.Context, meiopg model.MeioPagamento) (*model.MeioPagamento, error)
	Update(ctx context.Context, ID string, meioToChange *model.MeioPagamento) (bool, error)
	GetByID(ctx context.Context, ID string) (*model.MeioPagamento, error)
	GetAll(ctx context.Context, filters model.FilterMeioPg, limit, page int64) (*model.Paginate, error)
}

type MeioPgDataService struct {
	mdb mongodb.MongoDBInterface
}

func NewPeididoService(mongo_connection mongodb.MongoDBInterface) *MeioPgDataService {
	return &MeioPgDataService{
		mdb: mongo_connection,
	}
}

func (mpg *MeioPgDataService) Create(ctx context.Context, meiopg model.MeioPagamento) (*model.MeioPagamento, error) {
	collection := mpg.mdb.GetCollection("meiospamentos")

	dt := time.Now().Format(time.RFC3339)

	meiopg.Enabled = true
	meiopg.CreatedAt = dt
	meiopg.UpdatedAt = dt

	result, err := collection.InsertOne(ctx, meiopg)
	if err != nil {
		logger.Error("erro salvar meio de pagamento", err)
		return &meiopg, err
	}

	meiopg.ID = result.InsertedID.(primitive.ObjectID)

	return &meiopg, nil
}

func (mpg *MeioPgDataService) Update(ctx context.Context, ID string, meiopg model.MeioPagamento) (bool, error) {
	collection := mpg.mdb.GetCollection("meiospamentos")

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
			{Key: "meio_pg", Value: meiopg.Meiopg},
			{Key: "enabled", Value: meiopg.Enabled},
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

func (mpg *MeioPgDataService) GetByID(ctx context.Context, ID string) (*model.MeioPagamento, error) {

	collection := mpg.mdb.GetCollection("meiospamentos")

	meiopg := &model.MeioPagamento{}

	objectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {

		logger.Error("Error to parse ObjectIDFromHex", err)
		return nil, err
	}

	filter := bson.D{
		{Key: "_id", Value: objectID},
	}

	err = collection.FindOne(ctx, filter).Decode(meiopg)
	if err != nil {
		logger.Error("erro ao consultar meio de pagamento", err)
		return nil, err
	}

	return meiopg, nil
}

func (mpg *MeioPgDataService) GetAll(ctx context.Context, filters model.FilterMeioPg, limit, page int64) (*model.Paginate, error) {
	collection := mpg.mdb.GetCollection("meiospamentos")

	query := bson.M{}

	if filters.Meiopg != "" || filters.Enabled != "" {
		if filters.Meiopg != "" {
			query["meio_pg"] = bson.M{"$regex": fmt.Sprintf(".*%s.*", filters.Meiopg), "$options": "i"}
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
		logger.Error("erro ao consulta todos meios de pg", err)
		return nil, err
	}

	pagination := model.NewPaginate(limit, page, count)

	curr, err := collection.Find(ctx, query, pagination.GetPaginatedOpts())
	if err != nil {
		return nil, err
	}

	result := make([]*model.MeioPagamento, 0)
	for curr.Next(ctx) {
		mpg := &model.MeioPagamento{}
		if err := curr.Decode(mpg); err != nil {
			logger.Error("erro ao consulta todos meios de pg", err)
		}
		result = append(result, mpg)
	}

	pagination.Paginate(result)

	return pagination, nil
}
