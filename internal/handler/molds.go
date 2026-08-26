package handler

import (
	"net/http"

	"github.com/jb843051627/bell-foundry/internal/service"
)

type moldRequest struct {
	ProfileCode string `json:"profile_code"`
}

type moistureRequest struct {
	MoisturePct float64 `json:"moisture_pct"`
	Hours       float64 `json:"hours"`
}

type scrapRequest struct {
	Reason string `json:"reason"`
}

func (rt *Router) handleMolds(w http.ResponseWriter, r *http.Request, path string) {
	parts := pathParts(path)
	if len(parts) == 2 {
		switch r.Method {
		case http.MethodGet:
			items, err := rt.lab.ListMolds(r.Context())
			if err != nil {
				writeError(w, err)
				return
			}
			writeJSON(w, http.StatusOK, items)
		case http.MethodPost:
			var input moldRequest
			if err := decodeJSON(r, &input); err != nil {
				writeError(w, err)
				return
			}
			item, err := rt.lab.CreateMold(r.Context(), input.ProfileCode)
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
	if len(parts) < 3 {
		http.NotFound(w, r)
		return
	}
	id := parts[2]
	if len(parts) == 3 && r.Method == http.MethodGet {
		item, err := rt.lab.GetMold(r.Context(), id)
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
	var input moistureRequest
	var item any
	var err error
	switch parts[3] {
	case "drying-start":
		item, err = rt.lab.StartDrying(r.Context(), id)
	case "moisture":
		if decodeErr := decodeJSON(r, &input); decodeErr != nil {
			writeError(w, decodeErr)
			return
		}
		item, err = rt.lab.RecordMoisture(r.Context(), id, input.MoisturePct, input.Hours)
	case "drying-complete":
		if decodeErr := decodeJSON(r, &input); decodeErr != nil {
			writeError(w, decodeErr)
			return
		}
		item, err = rt.lab.CompleteDrying(r.Context(), id, input.Hours, input.MoisturePct)
	case "close":
		item, err = rt.lab.CloseMold(r.Context(), id)
	case "scrap":
		var request scrapRequest
		if decodeErr := decodeJSON(r, &request); decodeErr != nil {
			writeError(w, decodeErr)
			return
		}
		item, err = rt.lab.ScrapMold(r.Context(), id, request.Reason)
	default:
		http.NotFound(w, r)
		return
	}
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

// 保留一个显式引用，防止后续扩展时误删 service 依赖契约。
var _ *service.Lab
