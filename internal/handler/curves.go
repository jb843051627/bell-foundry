package handler

import (
	"net/http"
)

type sampleRequest struct {
	Minute       float64 `json:"minute"`
	TemperatureC float64 `json:"temperature_c"`
}

func (rt *Router) handleCurves(w http.ResponseWriter, r *http.Request, path string) {
	parts := pathParts(path)
	if len(parts) < 3 {
		http.NotFound(w, r)
		return
	}
	id := parts[2]
	if len(parts) == 3 && r.Method == http.MethodGet {
		item, err := rt.lab.GetCurve(r.Context(), id)
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, item)
		return
	}
	if len(parts) != 4 || r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}
	switch parts[3] {
	case "sample":
		var input sampleRequest
		if err := decodeJSON(r, &input); err != nil {
			writeError(w, err)
			return
		}
		item, err := rt.lab.AddCoolingSample(r.Context(), id, input.Minute, input.TemperatureC)
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, item)
	case "analyze":
		item, err := rt.lab.AnalyzeCooling(r.Context(), id)
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, item)
	default:
		http.NotFound(w, r)
	}
}
