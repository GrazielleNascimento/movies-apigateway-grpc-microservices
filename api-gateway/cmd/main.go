package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/movie-microservice/api-gateway/docs"
	httpSwagger "github.com/swaggo/http-swagger"

	grpcAdapter "github.com/movie-microservice/api-gateway/internal/adapters/grpc"
	"github.com/movie-microservice/api-gateway/internal/adapters/http/handlers"
	"github.com/movie-microservice/api-gateway/internal/adapters/http/middleware"
	"github.com/movie-microservice/api-gateway/internal/config"
	"github.com/movie-microservice/api-gateway/internal/core/services"
)

// @title Movie API Gateway
// @version 1.0
// @description API Gateway for Movie Microservice
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1

func main() {
	// Initialize logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	// Load configuration
	cfg := config.Load()
	if err := cfg.Validate(); err != nil {
		logger.Error("Invalid configuration", "error", err)
		os.Exit(1)
	}

	logger.Info("Starting API Gateway", "port", cfg.Server.Port)

	// Initialize gRPC client for movie service
	movieGRPCClient, err := grpcAdapter.NewMovieGRPCClient(cfg.MovieService.GRPCAddress, logger)
	if err != nil {
		logger.Error("Failed to connect to movie service", "error", err)
		os.Exit(1)
	}
	defer func() {
		if client, ok := movieGRPCClient.(*grpcAdapter.MovieGRPCClient); ok {
			client.Close()
		}
	}()

	// Initialize services
	movieService := services.NewMovieService(movieGRPCClient, logger)

	// Initialize handlers
	movieHandler := handlers.NewMovieHandler(movieService, logger)

	// Setup router
	router := mux.NewRouter()

	// Add middleware
	router.Use(middleware.CORS(logger))
	router.Use(middleware.Logging(logger))

	// API routes
	api := router.PathPrefix("/api/v1").Subrouter()

	// Movie routes
	api.HandleFunc("/movies", movieHandler.GetMovies).Methods("GET")
	api.HandleFunc("/movies/{id:[0-9]+}", movieHandler.GetMovie).Methods("GET")
	api.HandleFunc("/movies", movieHandler.CreateMovie).Methods("POST")
	api.HandleFunc("/movies/{id:[0-9]+}", movieHandler.DeleteMovie).Methods("DELETE")

	// Health check
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"healthy","timestamp":"%s"}`, time.Now().UTC().Format(time.RFC3339))
	}).Methods("GET")

	// Swagger documentation
	router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
	))

	// Create HTTP server
	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Channel to listen for interrupt signal to terminate server
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		logger.Info("HTTP server listening", "address", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Failed to start server", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal
	<-stop
	logger.Info("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", "error", err)
		os.Exit(1)
	}

	logger.Info("Server stopped")
}
