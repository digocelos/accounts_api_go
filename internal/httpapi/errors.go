package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/digocelo/account-api/internal/account"
)

type apiError struct {
	Error struct {
		Code    string                 `json:"code"`
		Message string                 `json:"message"`
		Details map[string]interface{} `json:"details,omitempty"`
	} `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeErr(w http.ResponseWriter, err error) {
	var ae apiError

	status := http.StatusInternalServerError
	ae.Error.Code = "internal"
	ae.Error.Message = "internal error"

	if errors.Is(err, account.ErrValidation) {
		status = http.StatusBadRequest
		ae.Error.Code = "validation_error"
		ae.Error.Message = err.Error()
	} else if errors.Is(err, account.ErrNotFound) {
		status = http.StatusNotFound
		ae.Error.Code = "not_found"
		ae.Error.Message = "resource not found"
	} else if errors.Is(err, account.ErrConflict) {
		status = http.StatusConflict
		ae.Error.Code = "conflict"
		ae.Error.Message = "version conflict"
	}
	writeJSON(w, status, ae)
}
