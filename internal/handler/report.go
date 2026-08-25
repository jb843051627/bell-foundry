package handler

import (
	"net/http"
	"time"
)

func (rt *Router) handleDailyReport(w http.ResponseWriter, r *http.Request) {
	if !requireMethod(w, r, http.MethodGet) {
		return
	}
	day := time.Now().UTC()
	if raw := r.URL.Query().Get("date"); raw != "" {
		parsed, err := time.Parse("2006-01-02", raw)
		if err != nil {
			writeError(w, err)
			return
		}
		day = parsed
	}
	item, err := rt.lab.DailyReport(r.Context(), day)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}
