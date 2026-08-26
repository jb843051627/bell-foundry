package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jb843051627/bell-foundry/internal/model"
	"github.com/jb843051627/bell-foundry/internal/notify"
)

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, err error) {
	status := http.StatusBadRequest
	if errors.Is(err, model.ErrNotFound) {
		status = http.StatusNotFound
	} else if errors.Is(err, model.ErrPreconditionFailed) || errors.Is(err, model.ErrBadTransition) {
		status = http.StatusConflict
	} else if errors.Is(err, notify.ErrDeliveryFailed) {
		status = http.StatusBadGateway
	}
	writeJSON(w, status, map[string]any{"error": err.Error()})
}

func decodeJSON(r *http.Request, target any) error {
	defer r.Body.Close()
	decoder := json.NewDecoder(http.MaxBytesReader(nil, r.Body, 1<<20))
	decoder.DisallowUnknownFields()
	return decoder.Decode(target)
}

func requireMethod(w http.ResponseWriter, r *http.Request, method string) bool {
	if r.Method != method {
		w.Header().Set("Allow", method)
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return false
	}
	return true
}
