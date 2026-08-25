package handler

import (
	"net/http"
)

type heatRequest struct {
	FurnaceNo   int     `json:"furnace_no"`
	SpecID      string  `json:"spec_id"`
	BatchID     string  `json:"batch_id"`
	ChargeKg    float64 `json:"charge_kg"`
	TargetTempC float64 `json:"target_temp_c"`
	WindowC     float64 `json:"window_c"`
}

type temperatureRequest struct {
	TemperatureC float64 `json:"temperature_c"`
}

type abortRequest struct {
	Reason string `json:"reason"`
}

func (rt *Router) handleHeats(w http.ResponseWriter, r *http.Request, path string) {
	parts := pathParts(path)
	if len(parts) == 2 {
		if r.Method == http.MethodGet {
			items, err := rt.lab.ListHeats(r.Context())
			if err != nil {
				writeError(w, err)
				return
			}
			writeJSON(w, http.StatusOK, items)
			return
		}
		if r.Method == http.MethodPost {
			var input heatRequest
			if err := decodeJSON(r, &input); err != nil {
				writeError(w, err)
				return
			}
			item, err := rt.lab.StartHeat(r.Context(), input.FurnaceNo, input.SpecID, input.BatchID, input.ChargeKg, input.TargetTempC, input.WindowC)
			if err != nil {
				writeError(w, err)
				return
			}
			writeJSON(w, http.StatusCreated, item)
			return
		}
	}
	if len(parts) < 3 {
		http.NotFound(w, r)
		return
	}
	id := parts[2]
	if len(parts) == 3 && r.Method == http.MethodGet {
		item, err := rt.lab.GetHeat(r.Context(), id)
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
	case "temperature":
		var input temperatureRequest
		if err := decodeJSON(r, &input); err != nil {
			writeError(w, err)
			return
		}
		item, err := rt.lab.RecordTemperature(r.Context(), id, input.TemperatureC)
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, item)
	case "ready":
		var input temperatureRequest
		if err := decodeJSON(r, &input); err != nil {
			writeError(w, err)
			return
		}
		item, err := rt.lab.MarkHeatReady(r.Context(), id, input.TemperatureC)
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, item)
	case "abort":
		var input abortRequest
		if err := decodeJSON(r, &input); err != nil {
			writeError(w, err)
			return
		}
		item, err := rt.lab.AbortHeat(r.Context(), id, input.Reason)
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, item)
	default:
		http.NotFound(w, r)
	}
}
