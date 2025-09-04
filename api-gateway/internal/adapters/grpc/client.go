package grpc

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
      
	pb "github.com/movie-microservice/proto/movies"
	"github.com/movie-microservice/api-gateway/internal/core/domain"
	"github.com/movie-microservice/api-gateway/internal/core/ports"
)

type MovieGRPCClient struct {
	client pb.MovieServiceClient
	conn   *grpc.ClientConn
	logger *slog.Logger
}

func NewMovieGRPCClient(serverAddress string, logger *slog.Logger) (ports.MovieServicePort, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, serverAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		logger.Error("Failed to connect to movie service", "address", serverAddress, "error", err)
		return nil, fmt.Errorf("failed to connect to movie service: %w", err)
	}

	client := pb.NewMovieServiceClient(conn)
	logger.Info("Successfully connected to movie service", "address", serverAddress)

	return &MovieGRPCClient{
		client: client,
		conn:   conn,
		logger: logger,
	}, nil
}

func (c *MovieGRPCClient) GetMovies(ctx context.Context, page, limit int32) ([]*domain.Movie, int32, error) {
	c.logger.Info("gRPC client: Getting movies", "page", page, "limit", limit)

	req := &pb.GetMoviesRequest{
		Page:  page,
		Limit: limit,
	}

	resp, err := c.client.GetMovies(ctx, req)
	if err != nil {
		c.logger.Error("gRPC client: Failed to get movies", "error", err)
		return nil, 0, fmt.Errorf("failed to get movies: %w", err)
	}

	if !resp.Success {
		c.logger.Error("gRPC client: Movie service returned error", "error", resp.Error)
		return nil, 0, fmt.Errorf("movie service error: %s", resp.Error)
	}

	// Convert protobuf movies to domain movies
	movies := make([]*domain.Movie, len(resp.Movies))
	for i, pbMovie := range resp.Movies {
		movies[i] = &domain.Movie{
			ID:    pbMovie.Id,
			Title: pbMovie.Title,
			Year:  pbMovie.Year,
		}
	}

	c.logger.Info("gRPC client: Successfully retrieved movies", "count", len(movies))
	return movies, resp.Total, nil
}

func (c *MovieGRPCClient) GetMovie(ctx context.Context, id int32) (*domain.Movie, error) {
	c.logger.Info("gRPC client: Getting movie", "id", id)

	req := &pb.GetMovieRequest{Id: id}

	resp, err := c.client.GetMovie(ctx, req)
	if err != nil {
		c.logger.Error("gRPC client: Failed to get movie", "id", id, "error", err)
		return nil, fmt.Errorf("failed to get movie: %w", err)
	}

	if !resp.Success {
		c.logger.Error("gRPC client: Movie service returned error", "id", id, "error", resp.Error)
		return nil, fmt.Errorf("movie service error: %s", resp.Error)
	}

	movie := &domain.Movie{
		ID:    resp.Movie.Id,
		Title: resp.Movie.Title,
		Year:  resp.Movie.Year,
	}

	c.logger.Info("gRPC client: Successfully retrieved movie", "id", id)
	return movie, nil
}

func (c *MovieGRPCClient) CreateMovie(ctx context.Context, title, year string) (*domain.Movie, error) {
	c.logger.Info("gRPC client: Creating movie", "title", title, "year", year)

	req := &pb.CreateMovieRequest{
		Title: title,
		Year:  year,
	}

	resp, err := c.client.CreateMovie(ctx, req)
	if err != nil {
		c.logger.Error("gRPC client: Failed to create movie", "title", title, "year", year, "error", err)
		return nil, fmt.Errorf("failed to create movie: %w", err)
	}

	if !resp.Success {
		c.logger.Error("gRPC client: Movie service returned error", "title", title, "year", year, "error", resp.Error)
		return nil, fmt.Errorf("movie service error: %s", resp.Error)
	}

	movie := &domain.Movie{
		ID:    resp.Movie.Id,
		Title: resp.Movie.Title,
		Year:  resp.Movie.Year,
	}

	c.logger.Info("gRPC client: Successfully created movie", "id", movie.ID)
	return movie, nil
}

func (c *MovieGRPCClient) DeleteMovie(ctx context.Context, id int32) error {
	c.logger.Info("gRPC client: Deleting movie", "id", id)

	req := &pb.DeleteMovieRequest{Id: id}

	resp, err := c.client.DeleteMovie(ctx, req)
	if err != nil {
		c.logger.Error("gRPC client: Failed to delete movie", "id", id, "error", err)
		return fmt.Errorf("failed to delete movie: %w", err)
	}

	if !resp.Success {
		c.logger.Error("gRPC client: Movie service returned error", "id", id, "error", resp.Error)
		return fmt.Errorf("movie service error: %s", resp.Error)
	}

	c.logger.Info("gRPC client: Successfully deleted movie", "id", id)
	return nil
}

func (c *MovieGRPCClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}