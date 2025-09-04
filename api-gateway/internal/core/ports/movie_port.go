package ports

import (
	"context"
	"net/http"
	"github.com/movie-microservice/api-gateway/internal/core/domain"
)

// MovieServicePort defines the contract for external movie service communication
type MovieServicePort interface {
	GetMovies(ctx context.Context, page, limit int32) ([]*domain.Movie, int32, error)
	GetMovie(ctx context.Context, id int32) (*domain.Movie, error)
	CreateMovie(ctx context.Context, title, year string) (*domain.Movie, error)
	DeleteMovie(ctx context.Context, id int32) error
}

// MovieHandler defines HTTP handler contract
type MovieHandler interface {
	GetMovies(w http.ResponseWriter, r *http.Request)
	GetMovie(w http.ResponseWriter, r *http.Request)
	CreateMovie(w http.ResponseWriter, r *http.Request)
	DeleteMovie(w http.ResponseWriter, r *http.Request)
}