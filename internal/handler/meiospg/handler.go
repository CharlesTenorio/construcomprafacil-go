package meiospg

import (
	"encoding/json"
	"strconv"

	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/katana/back-end/orcafacil-go/internal/config/logger"
	"github.com/katana/back-end/orcafacil-go/pkg/service/meiospg"
	"github.com/katana/back-end/orcafacil-go/pkg/service/validation"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/katana/back-end/orcafacil-go/pkg/model"
)

func createMeioPg(service meiospg.MeiosServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mpg := &model.MeioPagamento{}

		type Response struct {
			Message string `json:"message"`
		}
		var msg Response
		err := json.NewDecoder(r.Body).Decode(&mpg)

		if err != nil {
			logger.Error("error decoding request body", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		if mpg.Meiopg == "" {
			msg = Response{
				Message: "meio_pg e obrigatorio",
			}
			http.Error(w, msg.Message, http.StatusBadRequest)
			return
		}
		mpg.Meiopg = validation.CareString(mpg.Meiopg)

		if service.CheckExists(r.Context(), mpg.Meiopg) {
			msg = Response{
				Message: "meio_pg ja existe",
			}
			http.Error(w, msg.Message, http.StatusConflict)
			return
		}
		mpagamento := model.NewMeioPG(*mpg)

		_, err = service.Create(r.Context(), *mpagamento)
		if err != nil {
			logger.Error("erro ao acessar a camada de service do mpg", err)
			http.Error(w, "Error ou salvar Meio pg"+err.Error(), http.StatusInternalServerError)
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

func updateMeioPg(service meiospg.MeiosServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		idp := chi.URLParam(r, "id")
		logger.Info("PEGANDO O PARAMENTRO")

		_, err := service.GetByID(r.Context(), idp)
		if err != nil {
			http.Error(w, "Meio encontrada", http.StatusNotFound)
			return
		}

		mpg := &model.MeioPagamento{}
		nome := chi.URLParam(r, "nome")
		logger.Info("PEGANDO O NOME")
		logger.Info(nome)
		if nome == "" {
			http.Error(w, "o Nome do curso e obrigatório", http.StatusBadRequest)
			return
		}

		mpg.Meiopg = nome
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

func getByIdMeioPg(service meiospg.MeiosServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		idp := chi.URLParam(r, "id")
		logger.Info("PEGANDO O PARAMENTRO NA CONSULTA")
		result, err := service.GetByID(r.Context(), idp)
		if err != nil {
			logger.Error("erro ao acessar a camada de service do mpg no por id", err)
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

func getAllMeioPg(service meiospg.MeiosServiceInterface) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		filters := model.FilterMeioPg{
			Meiopg:  chi.URLParam(r, "nome"),
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
		err = json.NewEncoder(w).Encode(result)
		if err != nil {
			logger.Error("erro ao converto para json", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"MSG": "Error to parse User to JSON", "codigo": 500}`))
			return
		}
	})
}
