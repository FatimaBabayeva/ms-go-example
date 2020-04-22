package handler

import (
	"encoding/json"
	"github.com/FatimaBabayeva/ms-go-example/ctmerror"
	"github.com/FatimaBabayeva/ms-go-example/middleware"
	"github.com/FatimaBabayeva/ms-go-example/model"
	"github.com/FatimaBabayeva/ms-go-example/properties"
	"github.com/FatimaBabayeva/ms-go-example/repo"
	"github.com/FatimaBabayeva/ms-go-example/service"
	mid "github.com/go-chi/chi/middleware"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type messageHandler struct {
	service service.MessageService
}

var messageService = service.MessageServiceImpl{
	MsgRepo: &repo.MessageRepoImpl{},
}

// NewMessageHandler returns new message handler with predefined configuration
func NewMessageHandler(router *mux.Router) *mux.Router {
	router.Use(mid.Recoverer)
	router.Use(middleware.RequestParamsMiddleware)

	h := &messageHandler{service: &messageService}

	router.HandleFunc(properties.RootPath+"/message", h.saveMessage).Methods("POST")
	router.HandleFunc(properties.RootPath+"/message/{id}", h.getMessage).Methods("GET")
	router.HandleFunc(properties.RootPath+"/message/{id}", h.editMessage).Methods("PUT")
	router.HandleFunc(properties.RootPath+"/message/{id}", h.deleteMessage).Methods("DELETE")
	return router
}

func (h *messageHandler) saveMessage(w http.ResponseWriter, r *http.Request) {
	var m model.Message
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := h.service.SaveMessage(r.Context(), m)
	if err != nil {
		http.Error(w, err.Error(), err.(*ctmerror.MessageError).HttpCode())
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(result)
}

func (h *messageHandler) getMessage(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := h.service.GetMessageById(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), err.(*ctmerror.MessageError).HttpCode())
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func (h *messageHandler) editMessage(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var m model.Message
	err = json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := h.service.UpdateMessageById(r.Context(), id, m)
	if err != nil {
		http.Error(w, err.Error(), err.(*ctmerror.MessageError).HttpCode())
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func (h *messageHandler) deleteMessage(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.service.DeleteMessageById(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), err.(*ctmerror.MessageError).HttpCode())
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
