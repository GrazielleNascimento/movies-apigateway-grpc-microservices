package services

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/movie-microservice/movies-service/internal/core/domain"
	"github.com/movie-microservice/movies-service/internal/core/ports"
)

type MovieService struct {
	repo   ports.MovieRepository
	logger *slog.Logger
}

func NewMovieService(repo ports.MovieRepository, logger *slog.Logger) ports.MovieService {
	return &MovieService{
		repo:   repo,
		logger: logger,
	}
}

func (s *MovieService) GetMovies(ctx context.Context, filter domain.MovieFilter) ([]*domain.Movie, int32, error) {
	s.logger.Info("Getting movies with filter", "page", filter.Page, "limit", filter.Limit)

	// Validate filter
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 || filter.Limit > 100 {
		filter.Limit = 10
	}

	movies, err := s.repo.FindAll(ctx, filter)
	if err != nil {
		s.logger.Error("Failed to get movies", "error", err)
		return nil, 0, fmt.Errorf("failed to get movies: %w", err)
	}

	total, err := s.repo.Count(ctx)
	if err != nil {
		s.logger.Error("Failed to count movies", "error", err)
		return movies, 0, nil // Return movies even if count fails
	}

	s.logger.Info("Successfully retrieved movies", "count", len(movies), "total", total)
	return movies, total, nil
}

func (s *MovieService) GetMovie(ctx context.Context, id int32) (*domain.Movie, error) {
	s.logger.Info("Getting movie by ID", "id", id)

	if id <= 0 {
		return nil, domain.ErrInvalidMovieData
	}

	movie, err := s.repo.FindByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get movie", "id", id, "error", err)
		return nil, fmt.Errorf("failed to get movie with id %d: %w", id, err)
	}

	s.logger.Info("Successfully retrieved movie", "id", id, "title", movie.Title)
	return movie, nil
}

func (s *MovieService) CreateMovie(ctx context.Context, title, year string) (*domain.Movie, error) {
	s.logger.Info("Creating new movie", "title", title, "year", year)

	// Get next available ID
	nextID, err := s.repo.GetNextID(ctx)
	if err != nil {
		s.logger.Error("Failed to get next ID", "error", err)
		return nil, fmt.Errorf("failed to generate movie ID: %w", err)
	}

	// Create and validate movie
	movie, err := domain.NewMovie(nextID, title, year)
	if err != nil {
		s.logger.Error("Invalid movie data", "title", title, "year", year, "error", err)
		return nil, fmt.Errorf("invalid movie data: %w", err)
	}

	// Check if movie with same ID already exists
	exists, err := s.repo.ExistsByID(ctx, movie.ID)
	if err != nil {
		s.logger.Error("Failed to check movie existence", "id", movie.ID, "error", err)
		return nil, fmt.Errorf("failed to check movie existence: %w", err)
	}
	if exists {
		return nil, domain.ErrMovieAlreadyExists
	}

	// Save movie
	createdMovie, err := s.repo.Create(ctx, movie)
	if err != nil {
		s.logger.Error("Failed to create movie", "movie", movie, "error", err)
		return nil, fmt.Errorf("failed to create movie: %w", err)
	}

	s.logger.Info("Successfully created movie", "id", createdMovie.ID, "title", createdMovie.Title)
	return createdMovie, nil
}

func (s *MovieService) DeleteMovie(ctx context.Context, id int32) error {
	s.logger.Info("Deleting movie", "id", id)

	if id <= 0 {
		return domain.ErrInvalidMovieData
	}

	// Check if movie exists
	exists, err := s.repo.ExistsByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to check movie existence", "id", id, "error", err)
		return fmt.Errorf("failed to check movie existence: %w", err)
	}
	if !exists {
		return domain.ErrMovieNotFound
	}

	// Delete movie
	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error("Failed to delete movie", "id", id, "error", err)
		return fmt.Errorf("failed to delete movie with id %d: %w", id, err)
	}

	s.logger.Info("Successfully deleted movie", "id", id)
	return nil
}
