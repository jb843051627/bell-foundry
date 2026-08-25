package handler

import (
	"net/http"
	"strings"

	"github.com/jb843051627/bell-foundry/internal/model"
)

type resolveRequest struct {
	Status string `json:"status"`
}

func (rt *Router) handleQuality(w http.ResponseWriter, r *http.Request, path string) {
	if strings.HasPrefix(path, "/api/defects") {
		rt.handleDefects(w, r, path)
		return
	}
	rt.handleAlerts(w, r, path)
}

func (rt *Router) handleDefects(w http.ResponseWriter, r *http.Request, path string) {
	parts := pathParts(path)
	if len(parts) == 2 && r.Method == http.MethodGet {
		items, err := rt.lab.ListOpenDefects(r.Context())
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, items)
		return
	}
	if len(parts) == 2 && r.Method == http.MethodPost {
		var input model.Defect
		if err := decodeJSON(r, &input); err != nil {
			writeError(w, err)
			return
		}
		item, err := rt.lab.FileDefect(r.Context(), input)
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusCreated, item)
		return
	}
	if len(parts) == 4 && parts[3] == "resolve" && r.Method == http.MethodPost {
		var input resolveRequest
		if err := decodeJSON(r, &input); err != nil {
			writeError(w, err)
			return
		}
		item, err := rt.lab.ResolveDefect(r.Context(), parts[2], input.Status)
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, item)
		return
	}
	http.NotFound(w, r)
}

func (rt *Router) handleAlerts(w http.ResponseWriter, r *http.Request, path string) {
	parts := pathParts(path)
	if len(parts) == 2 && r.Method == http.MethodGet {
		openOnly := r.URL.Query().Get("open") != "false"
		items, err := rt.lab.ListAlerts(r.Context(), openOnly)
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, items)
		return
	}
	if len(parts) == 4 && parts[3] == "ack" && r.Method == http.MethodPost {
		item, err := rt.lab.AcknowledgeAlert(r.Context(), parts[2])
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, item)
		return
	}
	http.NotFound(w, r)
}
