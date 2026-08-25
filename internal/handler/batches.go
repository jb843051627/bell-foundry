package handler

import "net/http"

type batchPlanRequest struct {
	SpecID   string  `json:"spec_id"`
	TargetKg float64 `json:"target_kg"`
}

type weighRequest struct {
	Weighed map[string]float64 `json:"weighed"`
}

func (rt *Router) handleBatches(w http.ResponseWriter, r *http.Request, path string) {
	parts := pathParts(path)
	if len(parts) == 2 {
		if r.Method == http.MethodGet {
			items, err := rt.lab.ListBatches(r.Context())
			if err != nil {
				writeError(w, err)
				return
			}
			writeJSON(w, http.StatusOK, items)
			return
		}
		if r.Method == http.MethodPost {
			var input batchPlanRequest
			if err := decodeJSON(r, &input); err != nil {
				writeError(w, err)
				return
			}
			item, err := rt.lab.PlanBatch(r.Context(), input.SpecID, input.TargetKg)
			if err != nil {
				writeError(w, err)
				return
			}
			writeJSON(w, http.StatusCreated, item)
			return
		}
	}
	if len(parts) >= 3 {
		id := parts[2]
		if len(parts) == 3 && r.Method == http.MethodGet {
			item, err := rt.lab.GetBatch(r.Context(), id)
			if err != nil {
				writeError(w, err)
				return
			}
			writeJSON(w, http.StatusOK, item)
			return
		}
		if len(parts) == 4 && parts[3] == "weigh" && r.Method == http.MethodPost {
			var input weighRequest
			if err := decodeJSON(r, &input); err != nil {
				writeError(w, err)
				return
			}
			item, err := rt.lab.RecordWeighIn(r.Context(), id, input.Weighed)
			if err != nil {
				writeError(w, err)
				return
			}
			writeJSON(w, http.StatusOK, item)
			return
		}
	}
	http.NotFound(w, r)
}
