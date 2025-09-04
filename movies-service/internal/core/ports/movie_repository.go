package ports

import (
	"context"
	"github.com/movie-microservice/movies-service/internal/core/domain"
)

// MovieRepository defines the contract for movie data access
type MovieRepository interface {
	FindAll(ctx context.Context, filter domain.MovieFilter) ([]*domain.Movie, error)
	FindByID(ctx context.Context, id int32) (*domain.Movie, error)
	Create(ctx context.Context, movie *domain.Movie) (*domain.Movie, error)
	Delete(ctx context.Context, id int32) error
	Count(ctx context.Context) (int32, error)
	ExistsByID(ctx context.Context, id int32) (bool, error)
	GetNextID(ctx context.Context) (int32, error)
}

// MovieService defines the contract for movie business logic
type MovieService interface {
	GetMovies(ctx context.Context, filter domain.MovieFilter) ([]*domain.Movie, int32, error)
	GetMovie(ctx context.Context, id int32) (*domain.Movie, error)
	CreateMovie(ctx context.Context, title, year string) (*domain.Movie, error)
	DeleteMovie(ctx context.Context, id int32) error
}
