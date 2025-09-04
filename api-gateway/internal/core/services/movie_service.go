package services

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/movie-microservice/api-gateway/internal/core/domain"
	"github.com/movie-microservice/api-gateway/internal/core/ports"
)

type MovieService struct {
	moviePort ports.MovieServicePort
	logger    *slog.Logger
}

func NewMovieService(moviePort ports.MovieServicePort, logger *slog.Logger) *MovieService {
	return &MovieService{
		moviePort: moviePort,
		logger:    logger,
	}
}

func (s *MovieService) GetMovies(ctx context.Context, page, limit int32) ([]*domain.Movie, int32, error) {
	s.logger.Info("API Gateway: Getting movies", "page", page, "limit", limit)

	// Validate parameters
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	movies, total, err := s.moviePort.GetMovies(ctx, page, limit)
	if err != nil {
		s.logger.Error("API Gateway: Failed to get movies", "error", err)
		return nil, 0, fmt.Errorf("failed to get movies: %w", err)
	}

	s.logger.Info("API Gateway: Successfully retrieved movies", "count", len(movies), "total", total)
	return movies, total, nil
}

func (s *MovieService) GetMovie(ctx context.Context, id int32) (*domain.Movie, error) {
	s.logger.Info("API Gateway: Getting movie by ID", "id", id)

	if id <= 0 {
		return nil, fmt.Errorf("invalid movie ID: %d", id)
	}

	movie, err := s.moviePort.GetMovie(ctx, id)
	if err != nil {
		s.logger.Error("API Gateway: Failed to get movie", "id", id, "error", err)
		return nil, fmt.Errorf("failed to get movie: %w", err)
	}

	s.logger.Info("API Gateway: Successfully retrieved movie", "id", id, "title", movie.Title)
	return movie, nil
}

func (s *MovieService) CreateMovie(ctx context.Context, title, year string) (*domain.Movie, error) {
	s.logger.Info("API Gateway: Creating movie", "title", title, "year", year)

	if title == "" || year == "" {
		return nil, fmt.Errorf("title and year are required")
	}

	movie, err := s.moviePort.CreateMovie(ctx, title, year)
	if err != nil {
		s.logger.Error("API Gateway: Failed to create movie", "title", title, "year", year, "error", err)
		return nil, fmt.Errorf("failed to create movie: %w", err)
	}

	s.logger.Info("API Gateway: Successfully created movie", "id", movie.ID, "title", movie.Title)
	return movie, nil
}

func (s *MovieService) DeleteMovie(ctx context.Context, id int32) error {
	s.logger.Info("API Gateway: Deleting movie", "id", id)

	if id <= 0 {
		return fmt.Errorf("invalid movie ID: %d", id)
	}

	if err := s.moviePort.DeleteMovie(ctx, id); err != nil {
		s.logger.Error("API Gateway: Failed to delete movie", "id", id, "error", err)
		return fmt.Errorf("failed to delete movie: %w", err)
	}

	s.logger.Info("API Gateway: Successfully deleted movie", "id", id)
	return nil
}