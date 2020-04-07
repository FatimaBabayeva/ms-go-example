package handler

import (
	"github.com/gorilla/mux"
	"net/http"
)

type healthHandler struct{}

// HandleHealthRequest is for handling requests of kubernetes health and readiness checks
func HandleHealthRequest(router *mux.Router) {
	h := &healthHandler{}

	router.HandleFunc("/readiness", h.Health)
	router.HandleFunc("/health", h.Health)
}

// Health is a function that stands behind the health/readiness endpoint call
func (*healthHandler) Health(w http.ResponseWriter, r *http.Request) {
}
