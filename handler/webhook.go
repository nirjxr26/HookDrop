package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"github.com/your-username/hookdrop/store"
)

type WebhookHandler struct {
	Store *store.Store
}

func NewWebhookHandler(st *store.Store) *WebhookHandler {
	return &WebhookHandler{Store: st}
}

func (h *WebhookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/h/")
	if path == "" || path == r.URL.Path {
		http.NotFound(w, r)
		return
	}

	bucketID := path
	if idx := strings.IndexByte(path, '/'); idx >= 0 {
		bucketID = path[:idx]
	}

	switch {
	case r.Method == http.MethodPost && !strings.HasSuffix(r.URL.Path, "/stream"):
		h.handlePost(w, r, bucketID)
	case r.Method == http.MethodGet && strings.HasSuffix(r.URL.Path, "/stream"):
		h.handleStream(w, r, bucketID)
	case r.Method == http.MethodGet:
		h.handleList(w, r, bucketID)
	default:
		http.NotFound(w, r)
	}
}

func (h *WebhookHandler) handlePost(w http.ResponseWriter, r *http.Request, bucketID string) {
	traceID := uuid.NewString()
	ctx := withTraceLogger(r.Context(), traceID)
	logger := zerolog.Ctx(ctx)

	body, err := readBody(r)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "unable to read body"})
		logger.Error().Err(err).Str("bucket_id", bucketID).Msg("failed to read webhook body")
		return
	}

	event := store.WebhookEvent{
		TraceID:   traceID,
		BucketID:  bucketID,
		Method:    r.Method,
		Headers:   cloneHeaders(r.Header),
		Body:      body,
		SourceIP:  sourceIP(r),
		Timestamp: time.Now().UTC(),
	}

	h.Store.Add(bucketID, event)
	logger.Info().Str("bucket_id", bucketID).Str("method", event.Method).Str("source_ip", event.SourceIP).Time("timestamp", event.Timestamp).Int("body_size", len(body)).Msg("webhook received")

	writeJSON(w, http.StatusOK, map[string]string{"trace_id": traceID, "status": "received"})
}


func (h *WebhookHandler) handleList(w http.ResponseWriter, r *http.Request, bucketID string) {
	zerolog.Ctx(r.Context()).Info().Str("bucket_id", bucketID).Msg("webhook bucket listed")
	writeJSON(w, http.StatusOK, h.Store.List(bucketID))
}

func (h *WebhookHandler) handleStream(w http.ResponseWriter, r *http.Request, bucketID string) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	ctx := r.Context()
	logger := zerolog.Ctx(ctx)
	logger.Info().Str("bucket_id", bucketID).Msg("sse client connected")

	events := h.Store.Subscribe(bucketID)
	defer h.Store.Unsubscribe(bucketID, events)

	keepAlive := time.NewTicker(15 * time.Second)
	defer keepAlive.Stop()

	flusher.Flush()
	for {
		select {
		case <-ctx.Done():
			return
		case event := <-events:
			payload, err := json.Marshal(event)
			if err != nil {
				logger.Error().Err(err).Str("bucket_id", bucketID).Msg("failed to marshal sse event")
				continue
			}
			_, _ = fmt.Fprintf(w, "data: %s\n\n", payload)
			flusher.Flush()
		case <-keepAlive.C:
			_, _ = w.Write([]byte(": keepalive\n\n"))
			flusher.Flush()
		}
	}
}

func readBody(r *http.Request) (string, error) {
	defer func() {
		_ = r.Body.Close()
	}()
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func cloneHeaders(headers http.Header) map[string][]string {
	copied := make(map[string][]string, len(headers))
	for key, values := range headers {
		copied[key] = append([]string(nil), values...)
	}
	return copied
}

func sourceIP(r *http.Request) string {
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		parts := strings.Split(forwarded, ",")
		return strings.TrimSpace(parts[0])
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		return host
	}
	return r.RemoteAddr
}

func withTraceLogger(ctx context.Context, traceID string) context.Context {
	logger := zerolog.Ctx(ctx).With().Str("trace_id", traceID).Logger()
	return logger.WithContext(ctx)
}
