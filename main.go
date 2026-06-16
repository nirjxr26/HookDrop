package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/your-username/hookdrop/handler"
	"github.com/your-username/hookdrop/store"
)

const version = "1.0.0"

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	configureLogger()
	otelShutdown := initTelemetry(context.Background())
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = otelShutdown(ctx)
	}()

	port := getEnv("PORT", "8080")
	address := ":" + port

	webhookStore := store.New()
	seedMockData(webhookStore)
	webhookHandler := handler.NewWebhookHandler(webhookStore)
	dashboardHandler := handler.NewDashboardHandler()

	mux := http.NewServeMux()
	mux.Handle("/healthz", requestLogger(http.HandlerFunc(handler.Healthz)))
	mux.Handle("/readyz", requestLogger(http.HandlerFunc(handler.Readyz)))
	mux.Handle("/h/", requestLogger(webhookHandler))
	mux.HandleFunc("/logs/stream", handleLogStream)
	mux.Handle("/", requestLogger(dashboardHandler))

	server := &http.Server{
		Addr:         address,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	logger := log.Logger
	logger.Info().Str("port", port).Str("version", version).Msg("starting hookdrop")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal().Err(err).Msg("server failed")
		}
	}()

	<-stop
	logger.Info().Msg("shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Error().Err(err).Msg("graceful shutdown failed")
	}
}

func requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceID := r.Header.Get("X-Trace-Id")
		if traceID == "" {
			traceID = uuid.NewString()
		}
		logger := log.Logger.With().Str("trace_id", traceID).Logger()
		ctx := logger.WithContext(r.Context())
		ctx, span := startRequestSpan(ctx, r.Method, r.URL.Path, traceID)
		defer span.End()
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func configureLogger() {
	level := getEnv("LOG_LEVEL", "info")
	parsedLevel, err := zerolog.ParseLevel(level)
	if err != nil {
		parsedLevel = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(parsedLevel)
	log.Logger = zerolog.New(globalLogBroker).With().Timestamp().Logger()
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

type LogBroker struct {
	mu          sync.Mutex
	subscribers map[chan string]bool
}

var globalLogBroker = &LogBroker{
	subscribers: make(map[chan string]bool),
}

func (b *LogBroker) Write(p []byte) (n int, err error) {
	n, err = os.Stdout.Write(p)
	b.mu.Lock()
	defer b.mu.Unlock()
	msg := string(p)
	for ch := range b.subscribers {
		select {
		case ch <- msg:
		default:
		}
	}
	return
}

func (b *LogBroker) Subscribe() chan string {
	b.mu.Lock()
	defer b.mu.Unlock()
	ch := make(chan string, 100)
	b.subscribers[ch] = true
	return ch
}

func (b *LogBroker) Unsubscribe(ch chan string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.subscribers, ch)
	close(ch)
}

func handleLogStream(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	ctx := r.Context()
	ch := globalLogBroker.Subscribe()
	defer globalLogBroker.Unsubscribe(ch)

	// Keep-alive ticker
	keepAlive := time.NewTicker(15 * time.Second)
	defer keepAlive.Stop()

	flusher.Flush()
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-ch:
			_, _ = fmt.Fprintf(w, "data: %s\n\n", msg)
			flusher.Flush()
		case <-keepAlive.C:
			_, _ = w.Write([]byte(": keepalive\n\n"))
			flusher.Flush()
		}
	}
}

