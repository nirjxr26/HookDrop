package handler

import (
  "net/http"

  "github.com/rs/zerolog"
)

type DashboardHandler struct{}

func NewDashboardHandler() *DashboardHandler {
	return &DashboardHandler{}
}

func (h *DashboardHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  zerolog.Ctx(r.Context()).Info().Msg("dashboard requested")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(`<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>HookDrop</title>
  <style>
    body { font-family: system-ui, sans-serif; margin: 40px; line-height: 1.5; }
    code { background: #f4f4f4; padding: 0.2rem 0.35rem; border-radius: 4px; }
  </style>
</head>
<body>
  <h1>HookDrop</h1>
  <p>Mock webhook receiver with in-memory storage and SSE.</p>
  <ul>
    <li>POST /h/&lt;bucket-id&gt;</li>
    <li>GET /h/&lt;bucket-id&gt;</li>
    <li>GET /h/&lt;bucket-id&gt;/stream</li>
    <li>GET /healthz</li>
    <li>GET /readyz</li>
  </ul>
</body>
</html>`))
}
