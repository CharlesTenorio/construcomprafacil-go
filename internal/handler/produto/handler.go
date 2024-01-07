package produto

import (
	"encoding/json"
	"strconv"

	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/katana/back-end/orcafacil-go/internal/config/logger"
	"github.com/katana/back-end/orcafacil-go/internal/dto"
	"github.com/katana/back-end/orcafacil-go/pkg/service/categoria"
	"github.com/katana/back-end/orcafacil-go/pkg/service/produto"

	"github.com/katana/back-end/orcafacil-go/pkg/model"
)

func createProduto(service produto.ProdutoServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		produto := &model.Produto{}

		type Response struct {
			Message string `json:"message"`
		}
		var msg Response
		err := json.NewDecoder(r.Body).Decode(&produto)

		if err != nil {
			logger.Error("error decoding request body", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		if produto.Nome == "" {
			msg = Response{
				Message: "Nome e obrigatorio",
			}
			http.Error(w, msg.Message, http.StatusBadRequest)
			return
		}

		_, err = service.Create(r.Context(), *produto)
		if err != nil {
			logger.Error("erro ao acessar a camada de service do mpg", err)
			http.Error(w, "Error ou salvar categoria"+err.Error(), http.StatusInternalServerError)
			return
		}

		// Crie uma instância da estrutura com a mensagem desejada.
		msg = Response{
			Message: "Dados gravados com sucesso",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(msg)
	}
}

func updateProduto(service produto.ProdutoServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		type Response struct {
			Message string `json:"message"`
		}
		var msg Response
		clientID := chi.URLParam(r, "id")
		_, err := service.GetByID(r.Context(), clientID)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"MSG": "Client Not Found", "codigo": 404}`))
			return
		}
		ProdutoToChang := &model.Produto{}
		err = json.NewDecoder(r.Body).Decode(&ProdutoToChang)

		if err != nil {
			logger.Error("error decoding request body", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		if clientID == "" {
			msg = Response{
				Message: "id e obrigatorio",
			}
			http.Error(w, msg.Message, http.StatusBadRequest)
			return
		}

		if ProdutoToChang.Nome == "" {
			msg = Response{
				Message: "Nome e obrigatorio",
			}
			http.Error(w, msg.Message, http.StatusBadRequest)
			return
		}

		_, err = service.Update(r.Context(), clientID, *&ProdutoToChang)
		if err != nil {
			logger.Error("erro ao acessar a camada de service do mpg no upd", err)
			msg = Response{
				Message: "Erro ao atualizar Produto" + err.Error(),
			}
			http.Error(w, msg.Message, http.StatusBadRequest)
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

		result, err := service.GetByID(r.Context(), idp)
		if err != nil {
			logger.Error("erro ao acessar a camada de service da categoria no por id", err)
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"MSG": "Produto não encontrada", "codigo": 404}`))
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
		nome := r.URL.Query().Get("nome")
		enable := r.URL.Query().Get("enable")

		filters := model.FilterProduto{
			Nome:    nome,
			Enabled: enable,
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

func getListaProdutos(service categoria.CategoriaServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		idp := chi.URLParam(r, "id")
		logger.Info("passando ID CAT No handle")
		logger.Info(idp)

		limit, _ := strconv.ParseInt(r.URL.Query().Get("limit"), 10, 64)
		page, _ := strconv.ParseInt(r.URL.Query().Get("page"), 10, 64)

		result, err := service.ListProduto(r.Context(), idp, limit, page)
		if err != nil {
			logger.Error("erro ao acessar a camada de service da categoria no por id", err)
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"MSG": "Categoria não encontrada", "codigo": 404}`))
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

func updateProdutoAddFornecedor(service produto.ProdutoServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		type Response struct {
			Message string `json:"message"`
		}
		var msg Response
		clientID := chi.URLParam(r, "id")
		_, err := service.GetByID(r.Context(), clientID)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"MSG": "Client Not Found", "codigo": 404}`))
			return
		}
		fornecedor := []dto.FornecedoresEmPrd{}
		err = json.NewDecoder(r.Body).Decode(&fornecedor)

		if err != nil {
			logger.Error("error decoding request body", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		if clientID == "" {
			msg = Response{
				Message: "id e obrigatorio",
			}
			http.Error(w, msg.Message, http.StatusBadRequest)
			return
		}

		_, err = service.AddFornecedroes(r.Context(), clientID, &fornecedor)
		if err != nil {
			logger.Error("erro ao acessar a camada de service do mpg no upd", err)
			msg = Response{
				Message: "Erro ao atualizar Produto" + err.Error(),
			}
			http.Error(w, msg.Message, http.StatusBadRequest)
			return

		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"MSG": "Success", "codigo": 1})
	}
}
