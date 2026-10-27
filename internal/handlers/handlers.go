package handlers

import (
	"encoding/json"
	"net/http"

	"setho-pymnt/internal/models"
	"setho-pymnt/internal/service"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	Service *service.Service
}

func New(service *service.Service) *Handler {
	return &Handler{Service: service}
}

func (h *Handler) ListBusinesses(w http.ResponseWriter, r *http.Request) {
	businesses, err := h.Service.ListBusinesses()
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, businesses)
}

func (h *Handler) CreateBusiness(w http.ResponseWriter, r *http.Request) {
	var input models.BusinessInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid business payload")
		return
	}

	business, err := h.Service.CreateBusiness(input)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, business)
}

func (h *Handler) GetBusiness(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	business, err := h.Service.GetBusiness(id)
	if err != nil {
		h.writeError(w, http.StatusNotFound, "business not found")
		return
	}
	writeJSON(w, http.StatusOK, business)
}

func (h *Handler) ListTransactions(w http.ResponseWriter, r *http.Request) {
	businessID := r.URL.Query().Get("business_id")
	transactions, err := h.Service.ListTransactions(businessID)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, transactions)
}

func (h *Handler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var input models.TransactionInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid payment payload")
		return
	}

	transaction, err := h.Service.CreateTransaction(input)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, transaction)
}

func (h *Handler) GetTransaction(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	transaction, err := h.Service.GetTransaction(id)
	if err != nil {
		h.writeError(w, http.StatusNotFound, "transaction not found")
		return
	}
	writeJSON(w, http.StatusOK, transaction)
}

func (h *Handler) Summary(w http.ResponseWriter, r *http.Request) {
	summary, err := h.Service.Summary()
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, summary)
}

func (h *Handler) writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
