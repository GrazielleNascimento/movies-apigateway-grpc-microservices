package unit

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"testing"

	"github.com/movie-microservice/movies-service/internal/core/domain"
	"github.com/movie-microservice/movies-service/internal/core/services"
)

// Mock repository for testing
type MockMovieRepository struct {
	movies   map[int32]*domain.Movie
	nextID   int32
	findFail bool
}

func NewMockMovieRepository() *MockMovieRepository {
	return &MockMovieRepository{
		movies: make(map[int32]*domain.Movie),
		nextID: 1,
	}
}

func (m *MockMovieRepository) FindAll(ctx context.Context, filter domain.MovieFilter) ([]*domain.Movie, error) {
	if m.findFail {
		return nil, errors.New("database error")
	}

	var movies []*domain.Movie
	count := 0
	skip := (filter.Page - 1) * filter.Limit

	for _, movie := range m.movies {
		if count >= int(skip) && len(movies) < int(filter.Limit) {
			movies = append(movies, movie.Copy())
		}
		count++
	}

	return movies, nil
}

func (m *MockMovieRepository) FindByID(ctx context.Context, id int32) (*domain.Movie, error) {
	if m.findFail {
		return nil, errors.New("database error")
	}

	movie, exists := m.movies[id]
	if !exists {
		return nil, domain.ErrMovieNotFound
	}

	return movie.Copy(), nil
}

func (m *MockMovieRepository) Create(ctx context.Context, movie *domain.Movie) (*domain.Movie, error) {
	if m.findFail {
		return nil, errors.New("database error")
	}

	if _, exists := m.movies[movie.ID]; exists {
		return nil, domain.ErrMovieAlreadyExists
	}

	m.movies[movie.ID] = movie.Copy()
	return movie.Copy(), nil
}

func (m *MockMovieRepository) Delete(ctx context.Context, id int32) error {
	if m.findFail {
		return errors.New("database error")
	}

	if _, exists := m.movies[id]; !exists {
		return domain.ErrMovieNotFound
	}

	delete(m.movies, id)
	return nil
}

func (m *MockMovieRepository) Count(ctx context.Context) (int32, error) {
	if m.findFail {
		return 0, errors.New("database error")
	}

	return int32(len(m.movies)), nil
}

func (m *MockMovieRepository) ExistsByID(ctx context.Context, id int32) (bool, error) {
	if m.findFail {
		return false, errors.New("database error")
	}

	_, exists := m.movies[id]
	return exists, nil
}

func (m *MockMovieRepository) GetNextID(ctx context.Context) (int32, error) {
	if m.findFail {
		return 0, errors.New("database error")
	}

	id := m.nextID
	m.nextID++
	return id, nil
}

func TestMovieService_CreateMovie(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	mockRepo := NewMockMovieRepository()
	service := services.NewMovieService(mockRepo, logger)

	tests := []struct {
		name    string
		title   string
		year    string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid movie",
			title:   "Test Movie",
			year:    "2023",
			wantErr: false,
		},
		{
			name:    "empty title",
			title:   "",
			year:    "2023",
			wantErr: true,
			errMsg:  "invalid movie data: title cannot be empty",
		},
		{
			name:    "empty year",
			title:   "Test Movie",
			year:    "",
			wantErr: true,
			errMsg:  "invalid movie data: year cannot be empty",
		},
		{
			name:    "invalid year format",
			title:   "Test Movie",
			year:    "abc",
			wantErr: true,
			errMsg:  "invalid movie data: invalid year format",
		},
		{
			name:    "year too old",
			title:   "Test Movie",
			year:    "1700",
			wantErr: true,
			errMsg:  "invalid movie data: year must be between 1800 and current year + 10",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			movie, err := service.CreateMovie(context.Background(), tt.title, tt.year)

			if tt.wantErr {
				if err == nil {
					t.Errorf("CreateMovie() expected error but got none")
				}
				if tt.errMsg != "" && err != nil {
					if len(err.Error()) > 0 && err.Error()[:min(len(err.Error()), len(tt.errMsg))] != tt.errMsg[:min(len(err.Error()), len(tt.errMsg))] {
						t.Errorf("CreateMovie() error = %v, expected to contain %v", err, tt.errMsg)
					}
				}
			} else {
				if err != nil {
					t.Errorf("CreateMovie() unexpected error = %v", err)
				}
				if movie == nil {
					t.Errorf("CreateMovie() expected movie but got nil")
				}
				if movie != nil && movie.Title != tt.title {
					t.Errorf("CreateMovie() title = %v, want %v", movie.Title, tt.title)
				}
			}
		})
	}
}

func TestMovieService_GetMovie(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	mockRepo := NewMockMovieRepository()
	service := services.NewMovieService(mockRepo, logger)

	// Create a test movie
	testMovie, _ := domain.NewMovie(1, "Test Movie", "2023")
	mockRepo.movies[1] = testMovie

	tests := []struct {
		name    string
		id      int32
		wantErr bool
		wantID  int32
	}{
		{
			name:    "existing movie",
			id:      1,
			wantErr: false,
			wantID:  1,
		},
		{
			name:    "non-existing movie",
			id:      999,
			wantErr: true,
		},
		{
			name:    "invalid ID",
			id:      0,
			wantErr: true,
		},
		{
			name:    "negative ID",
			id:      -1,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			movie, err := service.GetMovie(context.Background(), tt.id)

			if tt.wantErr {
				if err == nil {
					t.Errorf("GetMovie() expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("GetMovie() unexpected error = %v", err)
				}
				if movie == nil {
					t.Errorf("GetMovie() expected movie but got nil")
				}
				if movie != nil && movie.ID != tt.wantID {
					t.Errorf("GetMovie() ID = %v, want %v", movie.ID, tt.wantID)
				}
			}
		})
	}
}

func TestMovieService_DeleteMovie(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	mockRepo := NewMockMovieRepository()
	service := services.NewMovieService(mockRepo, logger)

	// Create a test movie
	testMovie, _ := domain.NewMovie(1, "Test Movie", "2023")
	mockRepo.movies[1] = testMovie

	tests := []struct {
		name    string
		id      int32
		wantErr bool
	}{
		{
			name:    "existing movie",
			id:      1,
			wantErr: false,
		},
		{
			name:    "non-existing movie",
			id:      999,
			wantErr: true,
		},
		{
			name:    "invalid ID",
			id:      0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.DeleteMovie(context.Background(), tt.id)

			if tt.wantErr {
				if err == nil {
					t.Errorf("DeleteMovie() expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("DeleteMovie() unexpected error = %v", err)
				}
			}
		})
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
