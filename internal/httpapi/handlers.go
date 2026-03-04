package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/digocelo/account-api/internal/account"
	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	svc *account.Service
}

func NewHandlers(svc *account.Service) *Handlers {
	return &Handlers{svc: svc}
}

func (h *Handlers) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var in account.CreateInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeErr(w, validationf("invalid json"))
		return
	}
	acc, created, err := h.svc.Create(r.Context(), in)
	if err != nil {
		writeErr(w, err)
		return
	}

	if created {
		writeJSON(w, http.StatusCreated, acc)
		return
	}

	writeJSON(w, http.StatusOK, acc)
}

func (h *Handlers) GetAccount(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	acc, err := h.svc.Get(r.Context(), id)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, acc)
}

func (h *Handlers) UpdateAccount(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var in account.UpdateInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeErr(w, validationf("invalid json"))
		return
	}
	acc, err := h.svc.Update(r.Context(), id, in)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, acc)
}
