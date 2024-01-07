package fornecedor

import (
	"encoding/json"
	"strconv"

	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/katana/back-end/orcafacil-go/internal/config/logger"
	"github.com/katana/back-end/orcafacil-go/internal/dto"
	"github.com/katana/back-end/orcafacil-go/pkg/service/fornecedor"
	"github.com/katana/back-end/orcafacil-go/pkg/service/validation"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/katana/back-end/orcafacil-go/pkg/model"
)

func createFornecedor(service fornecedor.FornecedorServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		Fornecedor := &model.Fornecedor{}

		err := json.NewDecoder(r.Body).Decode(&Fornecedor)

		if err != nil {
			logger.Error("error decoding request body", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		if !validation.IsCNPJValid(Fornecedor.CNPJ) {
			http.Error(w, "Cnpj incorreto", http.StatusBadRequest)
			return
		}

		if service.GetByCnpj(r.Context(), Fornecedor.CNPJ) {
			http.Error(w, "CNPJ já exite", http.StatusBadRequest)
			return

		}

		_, err = service.Create(r.Context(), *Fornecedor)
		if err != nil {
			logger.Error("erro ao acessar a camada de service do mpg", err)
			http.Error(w, "Error ou salvar Fornecedor"+err.Error(), http.StatusInternalServerError)
			return
		}

		type Response struct {
			Message string `json:"message"`
		}

		// Crie uma instância da estrutura com a mensagem desejada.
		msg := Response{
			Message: "Dados gravados com sucesso",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(msg)
	}
}

func updateFornecedor(service fornecedor.FornecedorServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		idp := chi.URLParam(r, "id")
		logger.Info("PEGANDO O PARAMENTRO")

		_, err := service.GetByID(r.Context(), idp)
		if err != nil {
			http.Error(w, "Fornecedor nao encontrada", http.StatusNotFound)
			return
		}

		mpg := &model.Fornecedor{}
		nome := chi.URLParam(r, "nome")
		logger.Info("PEGANDO O NOME")
		logger.Info(nome)
		if nome == "" {
			http.Error(w, "o Nome do curso e obrigatório", http.StatusBadRequest)
			return
		}

		mpg.Nome = nome
		id, err := primitive.ObjectIDFromHex(idp)
		if err != nil {
			http.Error(w, "erro ao converter id", http.StatusBadRequest)

			return
		}

		mpg.ID = id
		_, err = service.Update(r.Context(), idp, *&mpg)
		if err != nil {
			logger.Error("erro ao acessar a camada de service do mpg no upd", err)
			http.Error(w, "Error ao atualizar meio de pagamento", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"MSG": "Success", "codigo": 1})
	}
}

func getByIdFornecedor(service fornecedor.FornecedorServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		idp := chi.URLParam(r, "id")
		logger.Info("PEGANDO O PARAMENTRO NA CONSULTA")
		result, err := service.GetByID(r.Context(), idp)
		if err != nil {
			logger.Error("erro ao acessar a camada de service da Fornecedor no por id", err)
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"MSG": "Fornecedor não encontrada", "codigo": 404}`))
			return
		}

		err = json.NewEncoder(w).Encode(result)
		if err != nil {
			logger.Error("erro ao converter em json", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"MSG": "Error to parse Bot to JSON", "codigo": 500}`))
			return
		}
	}
}

func getAllFornecedor(service fornecedor.FornecedorServiceInterface) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		filters := model.FilterFornecedor{
			Nome:    chi.URLParam(r, "nome"),
			Enabled: chi.URLParam(r, "enable"),
		}

		limit, _ := strconv.ParseInt(r.URL.Query().Get("limit"), 10, 64)
		page, _ := strconv.ParseInt(r.URL.Query().Get("page"), 10, 64)

		result, err := service.GetAll(r.Context(), filters, limit, page)
		if err != nil {
			logger.Error("erro ao acessar a camada de service do mpg no upd", err)
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"MSG": "User not found", "codigo": 404}`))
			return
		}

		// Configurando o cabeçalho para resposta JSON usando o middleware
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		// Escrevendo a resposta JSON
		err = json.NewEncoder(w).Encode(result)
		if err != nil {
			logger.Error("erro ao converter para json", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"MSG": "Error to parse User to JSON", "codigo": 500}`))
			return
		}
	})
}

func addProduto(service fornecedor.FornecedorServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		idp := chi.URLParam(r, "id")
		_, err := service.GetByID(r.Context(), idp)
		if err != nil {
			http.Error(w, "Fornecedor nao encontrada", http.StatusNotFound)
			return
		}

		// Decodifique os dados do corpo da requisição
		// Decodifique os dados do corpo da requisição
		var payload dto.ProdutosPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			logger.Error("Erro ao decodificar dados do corpo da requisição", err)
			http.Error(w, "Erro ao decodificar dados do corpo da requisição", http.StatusBadRequest)
			return
		}

		// Adicione a lista de produtos ao fornecedor
		_, err = service.AddProdutos(r.Context(), idp, payload.Produtos)
		if err != nil {
			logger.Error("Erro ao adicionar produtos ao fornecedor", err)
			http.Error(w, "Erro ao adicionar produtos ao fornecedor", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"MSG": "Produtos adicionados com sucesso", "codigo": 1})

	}
}
