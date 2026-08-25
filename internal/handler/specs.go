package handler

import (
	"net/http"

	"github.com/jb843051627/bell-foundry/internal/model"
)

func (rt *Router) handleSpecs(w http.ResponseWriter, r *http.Request, path string) {
	parts := pathParts(path)
	if len(parts) == 2 {
		switch r.Method {
		case http.MethodGet:
			items, err := rt.lab.ListSpecs(r.Context())
			if err != nil {
				writeError(w, err)
				return
			}
			writeJSON(w, http.StatusOK, items)
		case http.MethodPost:
			var input model.AlloySpec
			if err := decodeJSON(r, &input); err != nil {
				writeError(w, err)
				return
			}
			item, err := rt.lab.CreateSpec(r.Context(), input)
			if err != nil {
				writeError(w, err)
				return
			}
			writeJSON(w, http.StatusCreated, item)
		default:
			writeJSON(w, http.StatusMethodNotAllowed, nil)
		}
		return
	}
	if len(parts) == 3 && r.Method == http.MethodGet {
		item, err := rt.lab.GetSpec(r.Context(), parts[2])
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, item)
		return
	}
	http.NotFound(w, r)
}
