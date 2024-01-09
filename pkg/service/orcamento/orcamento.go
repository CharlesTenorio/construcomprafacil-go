package orcamento

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/katana/back-end/orcafacil-go/internal/config/logger"
	"github.com/katana/back-end/orcafacil-go/internal/dto"
	"github.com/katana/back-end/orcafacil-go/pkg/adapter/mongodb"
	"github.com/katana/back-end/orcafacil-go/pkg/adapter/rabbitmq"
	"github.com/katana/back-end/orcafacil-go/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type OrcamentoServiceInterface interface {
	Create(ctx context.Context, Orcamento model.Orcamento) (*model.Orcamento, error)
	Update(ctx context.Context, ID string, OrcamentoToChange *model.Orcamento) (bool, error)
	GetByID(ctx context.Context, ID string) (*model.Orcamento, error)
	GetAll(ctx context.Context, filters model.FilterOrcamento, limit, page int64) (*model.Paginate, error)
	GetToQueue(ctx context.Context, ID string) (*dto.OrcamentoFila, error)
	SendToQuequeFil(ctx context.Context, fila *dto.OrcamentoFila) (bool, error)
}

type OrcamentoDataService struct {
	mdb       mongodb.MongoDBInterface
	rabbit_mq rabbitmq.RabbitInterface
}

func NewOrcamentoervice(mongo_connection mongodb.MongoDBInterface, rabbit_connection rabbitmq.RabbitInterface) *OrcamentoDataService {
	return &OrcamentoDataService{
		mdb:       mongo_connection,
		rabbit_mq: rabbit_connection,
	}
}

func (fornec *OrcamentoDataService) Create(ctx context.Context, Orcamento model.Orcamento) (*model.Orcamento, error) {
	collection := fornec.mdb.GetCollection("Orcamentoes")

	dt := time.Now().Format(time.RFC3339)

	Orcamento.Enabled = true
	Orcamento.CreatedAt = dt
	Orcamento.UpdatedAt = dt
	Orcamento.ID = primitive.NewObjectID()

	result, err := collection.InsertOne(ctx, Orcamento)
	if err != nil {
		logger.Error("erro salvar  Orcamento", err)
		return &Orcamento, err
	}

	Orcamento.ID = result.InsertedID.(primitive.ObjectID)

	return &Orcamento, nil
}

func (fornec *OrcamentoDataService) Update(ctx context.Context, ID string, OrcamentoToChange *model.Orcamento) (bool, error) {
	collection := fornec.mdb.GetCollection("Orcamentoes")

	opts := options.Update().SetUpsert(true)

	objectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		logger.Error("Error to parse ObjectIDFromHex", err)
		return false, err
	}

	filter := bson.D{
		{Key: "_id", Value: objectID},
	}

	// Criando um mapa para armazenar as atualizações dinâmicas
	updateFields := make(map[string]interface{})

	// Adicionando todos os campos que você deseja atualizar
	updateFields["enabled"] = OrcamentoToChange.Enabled
	updateFields["updated_at"] = time.Now().Format(time.RFC3339)
	updateFields["sugestaoprazoEntrega"] = OrcamentoToChange.SugestaoPrazoEntrega
	updateFields["finalizado"] = OrcamentoToChange.Finalizado
	updateFields["pegarEstabelecimento"] = OrcamentoToChange.PegarEstabelecimento
	updateFields["listaCliente"] = OrcamentoToChange.GrupoDeCliente
	updateFields["enderecoCliente"] = OrcamentoToChange.EnderecoCliente
	updateFields["status"] = OrcamentoToChange.StatusOrecamento
	updateFields["fornecedor"] = OrcamentoToChange.Fornecedores
	// Adicione mais campos conforme necessário...

	// Criando a atualização dinâmica com os campos fornecidos
	var updateFieldsDoc bson.D
	for key, value := range updateFields {
		updateFieldsDoc = append(updateFieldsDoc, bson.E{Key: key, Value: value})
	}

	// Criando a atualização final
	update := bson.D{{Key: "$set", Value: updateFieldsDoc}}

	// Executando a atualização
	_, err = collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		logger.Error("Error while updating data", err)
		return false, err
	}

	return true, nil
}

func (fornec *OrcamentoDataService) GetByID(ctx context.Context, ID string) (*model.Orcamento, error) {

	collection := fornec.mdb.GetCollection("Orcamentoes")

	Orcamento := &model.Orcamento{}

	objectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {

		logger.Error("Error to parse ObjectIDFromHex", err)
		return nil, err
	}

	filter := bson.D{
		{Key: "_id", Value: objectID},
	}

	err = collection.FindOne(ctx, filter).Decode(Orcamento)
	if err != nil {
		logger.Error("erro ao consultar Orcamento", err)
		return nil, err
	}

	return Orcamento, nil
}

func (fornec *OrcamentoDataService) GetAll(ctx context.Context, filters model.FilterOrcamento, limit, page int64) (*model.Paginate, error) {
	collection := fornec.mdb.GetCollection("Orcamentoes")

	query := bson.M{}

	if filters.DataInical != "" && filters.DataFinal == "" {
		query["data_solicitacao"] = bson.M{
			"$gte": filters.DataInical,
			"$lte": filters.DataFinal,
		}
	}

	if filters.Enabled != "" {
		enable, err := strconv.ParseBool(filters.Enabled)
		if err != nil {
			logger.Error("erro converter campo enabled", err)
			return nil, err
		}
		query["enabled"] = enable
	}

	count, err := collection.CountDocuments(ctx, query, &options.CountOptions{})

	if err != nil {
		logger.Error("erro ao consultar  Orcamentos", err)
		return nil, err
	}

	pagination := model.NewPaginate(limit, page, count)

	curr, err := collection.Find(ctx, query, pagination.GetPaginatedOpts())
	if err != nil {
		return nil, err
	}

	result := make([]*model.Orcamento, 0)
	for curr.Next(ctx) {
		fornec := &model.Orcamento{}
		if err := curr.Decode(fornec); err != nil {
			logger.Error("erro ao consultar todas as Orcamentos", err)
		}
		result = append(result, fornec)
	}

	pagination.Paginate(result)

	return pagination, nil
}
func (fornec *OrcamentoDataService) GetByCnpj(ctx context.Context, Cnpj string) bool {

	collection := fornec.mdb.GetCollection("Orcamentoes")

	// Utilizando o método CountDocuments para verificar a existência
	filter := bson.D{{Key: "cnpj", Value: Cnpj}}
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		logger.Error("erro ao consultar Orcamento", err)
		return false
	}

	// Se count for maior que zero, o Orcamento existe
	return count > 0
}

func (fornec *OrcamentoDataService) EnviarParaFila(ctx context.Context, fila *dto.OrcamentoFila) (bool, error) {

	data, err := json.Marshal(fila)
	if err != nil {
		logger.Error("Erro ao converter OrcamentoFila para []byte: ", err)
		return false, err
	}

	msg := &rabbitmq.Message{
		Data:        data,
		ContentType: "application/json; charset=utf-8",
	}

	err = fornec.rabbit_mq.SenderRb(ctx, "QUEUE_PRD_PARA_FORNECEDORES", msg)
	if err != nil {
		logger.Error("Erro ao envair produtos do orcamento pra faila ", err)
		return false, err
	}

	return true, nil
}
