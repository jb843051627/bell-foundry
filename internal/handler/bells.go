package handler

import "net/http"

type partialRequest struct {
	Measured []float64 `json:"measured"`
}

type tuningRequest struct {
	TuneLimit   float64 `json:"tune_limit"`
	RetuneLimit float64 `json:"retune_limit"`
}

func (rt *Router) handleBells(w http.ResponseWriter, r *http.Request, path string) {
	if path == "/api/bells/retune" {
		if !requireMethod(w, r, http.MethodGet) {
			return
		}
		items, err := rt.lab.ListBellsNeedingRetune(r.Context())
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, items)
		return
	}
	parts := pathParts(path)
	if len(parts) < 3 {
		http.NotFound(w, r)
		return
	}
	id := parts[2]
	if len(parts) == 3 && r.Method == http.MethodGet {
		item, err := rt.lab.GetBell(r.Context(), id)
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
	case "partials":
		var input partialRequest
		if err := decodeJSON(r, &input); err != nil {
			writeError(w, err)
			return
		}
		item, err := rt.lab.RecordPartials(r.Context(), id, input.Measured)
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, item)
	case "evaluate":
		var input tuningRequest
		if err := decodeJSON(r, &input); err != nil {
			writeError(w, err)
			return
		}
		item, err := rt.lab.EvaluateBell(r.Context(), id, input.TuneLimit, input.RetuneLimit)
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, item)
	case "plan":
		item, err := rt.lab.TuningPlan(r.Context(), id)
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, item)
	default:
		http.NotFound(w, r)
	}
}
