package handler

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog"
)

func Healthz(w http.ResponseWriter, r *http.Request) {
	zerolog.Ctx(r.Context()).Info().Msg("health check")
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func Readyz(w http.ResponseWriter, r *http.Request) {
	zerolog.Ctx(r.Context()).Info().Msg("readiness check")
	writeJSON(w, http.StatusOK, map[string]string{"status": "ready"})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
