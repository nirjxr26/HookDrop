package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
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
	webhookHandler := handler.NewWebhookHandler(webhookStore)
	dashboardHandler := handler.NewDashboardHandler()

	mux := http.NewServeMux()
	mux.Handle("/healthz", requestLogger(http.HandlerFunc(handler.Healthz)))
	mux.Handle("/readyz", requestLogger(http.HandlerFunc(handler.Readyz)))
	mux.Handle("/h/", requestLogger(webhookHandler))
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
	log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

