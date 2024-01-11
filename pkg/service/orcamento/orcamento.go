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
	GetToQueuePrdToFornec(ctx context.Context, ID string) (*dto.OrcamentoFilaPrdFornecedor, error)
	//SendToQuequePrdFornecedor(ctx context.Context, fila *dto.OrcamentoFilaPrdFornecedor) (bool, error)
	GetPrdContacao(ctx context.Context, Orcamento *model.Orcamento) (dto.ProdutoEnviadoParaFilaDeOrcamentoDTO, error)
	SenderePrdOrcamentoFila(ctx context.Context, orcamentoPrdsFila *dto.ProdutoEnviadoParaFilaDeOrcamentoDTO) bool
}

type OrcamentoDataService struct {
	mdb       mongodb.MongoDBInterface
	rabbit_mq rabbitmq.RabbitInterface
}

func NewOrcamentoService(mongo_connection mongodb.MongoDBInterface, rabbit_connection rabbitmq.RabbitInterface) *OrcamentoDataService {
	return &OrcamentoDataService{
		mdb:       mongo_connection,
		rabbit_mq: rabbit_connection,
	}
}

func (fornec *OrcamentoDataService) Create(ctx context.Context, Orcamento model.Orcamento) (*model.Orcamento, error) {
	collection := fornec.mdb.GetCollection("orcamentos")

	dt := time.Now().Format(time.RFC3339)

	Orcamento.Enabled = true
	Orcamento.CreatedAt = dt
	Orcamento.UpdatedAt = dt
	Orcamento.ID = primitive.NewObjectID()
	var err error
	result, err := collection.InsertOne(ctx, Orcamento)
	if err != nil {
		logger.Error("erro salvar  Orcamento", err)
		return &Orcamento, err
	}

	Orcamento.ID = result.InsertedID.(primitive.ObjectID)

	prds, err := fornec.GetPrdContacao(ctx, &Orcamento)
	if err != nil {
		logger.Error("Erro ao pegar produtos do orcamento", err)
		return &Orcamento, err
	}
	productsJSON, err := json.Marshal(prds.Produtos)
	if err != nil {
		//fmt.Println("Error marshalling products to JSON:", err)
		return &Orcamento, err
	}

	// Log the JSON string
	logger.Info(string(productsJSON))
	send := fornec.SenderePrdOrcamentoFila(ctx, &prds)
	if send == false {
		logger.Error("Erro ao enviar produtos do orcamento para fila", err)
		return &Orcamento, err
	}

	return &Orcamento, nil
}

func (fornec *OrcamentoDataService) Update(ctx context.Context, ID string, OrcamentoToChange *model.Orcamento) (bool, error) {
	collection := fornec.mdb.GetCollection("orcamentos")

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

	collection := fornec.mdb.GetCollection("orcamentos")

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
	collection := fornec.mdb.GetCollection("orcamentos")

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

	collection := fornec.mdb.GetCollection("orcamentos")

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

func (fornec *OrcamentoDataService) EnviarParaFila(ctx context.Context, fila *dto.OrcamentoFilaPrdFornecedor) (bool, error) {

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

func (fornec *OrcamentoDataService) GetToQueuePrdToFornec(ctx context.Context, ID string) (*dto.OrcamentoFilaPrdFornecedor, error) {
	fetchedOrcamento, err := fornec.GetByID(ctx, ID)
	if err != nil {
		return nil, err
	}

	// Populate OrcamentoFila based on fetchedOrcamento
	orcamentoFila := dto.OrcamentoFilaPrdFornecedor{
		IdOrcamento: fetchedOrcamento.ID,
		StatusFila:  fetchedOrcamento.StatusOrecamento,
		// You may need to populate other fields based on your requirements
	}

	// Fetch other necessary information from the fetchedOrcamento
	// For example, populate Fornecedor field
	for _, fornecedor := range fetchedOrcamento.Fornecedores {
		orcamentoFila.Fornecedor = fornecedor
		break // Assuming you only want the first Fornecedor
	}

	return &orcamentoFila, nil
}
func (fornec *OrcamentoDataService) GetPrdContacao(ctx context.Context, Orcamento *model.Orcamento) (dto.ProdutoEnviadoParaFilaDeOrcamentoDTO, error) {
	var prd dto.ProdutoContacaoDTO
	var prds []dto.ProdutoContacaoDTO
	for _, produto := range Orcamento.ProdutosContacao {
		prd.ProdutoID = produto.ProdutoID
		prd.Quantidade = produto.Quantidade
		prd.Nome = produto.Nome
		prds = append(prds, prd)
	}
	orcamentoPrdFila := dto.ProdutoEnviadoParaFilaDeOrcamentoDTO{
		IdOrcamento: Orcamento.ID,
		DataEnvio:   Orcamento.DataSolicitacao,
		Produtos:    prds,
	}
	// You can return nil as the second parameter if there is no error
	return orcamentoPrdFila, nil
}
func (fornec *OrcamentoDataService) SenderePrdOrcamentoFila(ctx context.Context, orcamentoPrdsFila *dto.ProdutoEnviadoParaFilaDeOrcamentoDTO) bool {
	data, err := json.Marshal(orcamentoPrdsFila)
	if err != nil {
		logger.Error("Erro ao converter OrcamentoFila para []byte: ", err)
		return false
	}

	msg := &rabbitmq.Message{
		Data:        data,
		ContentType: "application/json; charset=utf-8",
	}

	err = fornec.rabbit_mq.SenderRb(ctx, "QUEUE_PRDS_PARA_COTACAO", msg)
	if err != nil {
		logger.Error("Erro ao envair produtos do orcamento pra faila ", err)
		return false
	}

	return true
}
