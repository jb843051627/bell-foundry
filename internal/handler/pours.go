package handler

import "net/http"

type pourRequest struct {
	HeatID      string  `json:"heat_id"`
	MoldID      string  `json:"mold_id"`
	PourTempC   float64 `json:"pour_temp_c"`
	FlowSeconds int     `json:"flow_seconds"`
}

func (rt *Router) handlePours(w http.ResponseWriter, r *http.Request, path string) {
	parts := pathParts(path)
	if len(parts) == 2 && r.Method == http.MethodGet {
		items, err := rt.lab.ListPours(r.Context())
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, items)
		return
	}
	if len(parts) == 2 && r.Method == http.MethodPost {
		var input pourRequest
		if err := decodeJSON(r, &input); err != nil {
			writeError(w, err)
			return
		}
		item, err := rt.lab.ExecutePour(r.Context(), input.HeatID, input.MoldID, input.PourTempC, input.FlowSeconds)
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusCreated, item)
		return
	}
	if len(parts) == 3 && r.Method == http.MethodGet {
		item, err := rt.lab.GetPour(r.Context(), parts[2])
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, item)
		return
	}
	http.NotFound(w, r)
}
