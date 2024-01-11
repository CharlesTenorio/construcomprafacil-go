package orcamento

import (
	"encoding/json"
	"strconv"

	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/katana/back-end/orcafacil-go/internal/config/logger"
	"github.com/katana/back-end/orcafacil-go/pkg/service/orcamento"

	"github.com/katana/back-end/orcafacil-go/pkg/model"
)

func createOrcamento(service orcamento.OrcamentoServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orcamento := &model.Orcamento{}

		type Response struct {
			Message string `json:"message"`
		}
		var msg Response
		err := json.NewDecoder(r.Body).Decode(&orcamento)

		if err != nil {
			logger.Error("error decoding request body", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		_, err = service.Create(r.Context(), *orcamento)
		if err != nil {
			logger.Error("erro ao acessar a camada de service do mpg", err)
			http.Error(w, "Error ou salvar orcamento"+err.Error(), http.StatusInternalServerError)
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

func updateOrcamento(service orcamento.OrcamentoServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		orcamento := &model.Orcamento{}
		err := json.NewDecoder(r.Body).Decode(&orcamento)

		if err != nil {
			logger.Error("error decoding request body", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		idp := orcamento.ID.String()
		logger.Info("PEGANDO O PARAMENTRO")

		_, err = service.GetByID(r.Context(), idp)
		if err != nil {
			http.Error(w, "orcamento nao encontrada", http.StatusNotFound)
			return
		}

		_, err = service.Update(r.Context(), idp, *&orcamento)
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

func getByIdOrcamento(service orcamento.OrcamentoServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		idp := chi.URLParam(r, "id")
		logger.Info("PEGANDO O PARAMENTRO NA CONSULTA")
		result, err := service.GetByID(r.Context(), idp)
		if err != nil {
			logger.Error("erro ao acessar a camada de service da orcamento no por id", err)
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"MSG": "orcamento não encontrada", "codigo": 404}`))
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

func getAllOrcamento(service orcamento.OrcamentoServiceInterface) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		filters := model.FilterOrcamento{
			DataInical: chi.URLParam(r, "data_inicial"),
			DataFinal:  chi.URLParam(r, "data_final"),
			Enabled:    chi.URLParam(r, "enable"),
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
