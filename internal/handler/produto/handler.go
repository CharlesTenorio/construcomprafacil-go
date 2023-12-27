package produto

import (
	"encoding/json"
	"strconv"

	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/katana/back-end/orcafacil-go/internal/config/logger"
	"github.com/katana/back-end/orcafacil-go/pkg/service/produto"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/katana/back-end/orcafacil-go/pkg/model"
)

func createProduto(service produto.ProdutoServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prd := &model.Produto{}
		err := json.NewDecoder(r.Body).Decode(&prd)

		if err != nil {
			logger.Error("error decoding request body", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		if prd.Nome == "" {
			http.Error(w, "o Nome do meio de pagamento obrigatorio", http.StatusBadRequest)
			return
		}

		_, err = service.CreateProduto(r.Context(), *prd)
		if err != nil {
			logger.Error("erro ao acessar a camada de service do prd", err)
			http.Error(w, "Error ou salvar Meio pg"+err.Error(), http.StatusInternalServerError)
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

func updateProduto(service produto.ProdutoServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prd := &model.Produto{}
		idp := chi.URLParam(r, "id")
		err := json.NewDecoder(r.Body).Decode(&prd)

		if err != nil {
			logger.Error("error decoding request body", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		if prd.Nome == "" {
			http.Error(w, "o Nome do meio de pagamento obrigatorio", http.StatusBadRequest)
			return
		}

		_, err = service.GetProdutoByID(r.Context(), idp)
		if err != nil {
			http.Error(w, "Meio encontrada", http.StatusNotFound)
			return
		}

		id, err := primitive.ObjectIDFromHex(idp)
		if err != nil {
			http.Error(w, "erro ao converter id", http.StatusBadRequest)

			return
		}

		prd.ID = id
		_, err = service.UpdateProduto(r.Context(), idp, *&prd)
		if err != nil {
			logger.Error("erro ao acessar a camada de service do prd no upd", err)
			http.Error(w, "Error ao atualizar meio de pagamento", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"MSG": "Success", "codigo": 1})
	}
}

func getByIdProduto(service produto.ProdutoServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		idp := chi.URLParam(r, "id")
		logger.Info("PEGANDO O PARAMENTRO NA CONSULTA")
		result, err := service.GetProdutoByID(r.Context(), idp)
		if err != nil {
			logger.Error("erro ao acessar a camada de service do prd no por id", err)
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"MSG": "Meio de pagamento não encontrado", "codigo": 404}`))
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

func getAllProduto(service produto.ProdutoServiceInterface) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		filters := model.FilterProduto{
			Nome:    chi.URLParam(r, "nome"),
			Enabled: chi.URLParam(r, "enable"),
		}

		limit, _ := strconv.ParseInt(r.URL.Query().Get("limit"), 10, 64)
		page, _ := strconv.ParseInt(r.URL.Query().Get("page"), 10, 64)

		result, err := service.GetAllProdutos(r.Context(), filters, limit, page)
		if err != nil {
			logger.Error("erro ao acessar a camada de service do prd no upd", err)
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"MSG": "User not found", "codigo": 404}`))
			return
		}
		err = json.NewEncoder(w).Encode(result)
		if err != nil {
			logger.Error("erro ao converto para json", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"MSG": "Error to parse User to JSON", "codigo": 500}`))
			return
		}
	})
}
