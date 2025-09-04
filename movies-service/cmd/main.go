package main

import (
	"context"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/movie-microservice/proto/movies"
	"github.com/movie-microservice/movies-service/internal/adapters/database"
	grpcAdapter "github.com/movie-microservice/movies-service/internal/adapters/grpc"
	"github.com/movie-microservice/movies-service/internal/config"
	"github.com/movie-microservice/movies-service/internal/core/services"
)

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

	logger.Info("Starting movies service", "grpc_port", cfg.GRPC.Port)

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	mongoClient, err := database.Connect(ctx, cfg.Database.ConnectionString, logger)
	if err != nil {
		logger.Error("Failed to connect to MongoDB", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := database.Disconnect(context.Background(), mongoClient, logger); err != nil {
			logger.Error("Failed to disconnect from MongoDB", "error", err)
		}
	}()

	// Initialize repository
	movieRepo := database.NewMongoMovieRepository(mongoClient, cfg.Database.DatabaseName, logger)

	// Initialize service
	movieService := services.NewMovieService(movieRepo, logger)

	// Initialize gRPC server
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(unaryInterceptor(logger)),
	)

	// Register movie service
	movieGRPCService := grpcAdapter.NewMovieServer(movieService, logger)
	pb.RegisterMovieServiceServer(grpcServer, movieGRPCService)

	// Enable reflection for grpcurl testing
	reflection.Register(grpcServer)

	// Start gRPC server
	lis, err := net.Listen("tcp", ":"+cfg.GRPC.Port)
	if err != nil {
		logger.Error("Failed to listen", "port", cfg.GRPC.Port, "error", err)
		os.Exit(1)
	}

	// Channel to listen for interrupt signal to terminate server
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		logger.Info("gRPC server listening", "address", lis.Addr().String())
		if err := grpcServer.Serve(lis); err != nil {
			logger.Error("Failed to serve gRPC", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal
	<-stop
	logger.Info("Shutting down gRPC server...")

	// Graceful shutdown
	grpcServer.GracefulStop()
	logger.Info("Server stopped")
}

// Unary interceptor for logging
func unaryInterceptor(logger *slog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		
		resp, err := handler(ctx, req)
		
		duration := time.Since(start)
		
		if err != nil {
			logger.Error("gRPC request failed",
				"method", info.FullMethod,
				"duration", duration,
				"error", err,
			)
		} else {
			logger.Info("gRPC request completed",
				"method", info.FullMethod,
				"duration", duration,
			)
		}
		
		return resp, err
	}
}
