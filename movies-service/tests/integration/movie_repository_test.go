package integration

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/movie-microservice/movies-service/internal/adapters/database"
	"github.com/movie-microservice/movies-service/internal/core/domain"
)

func TestMovieRepository_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration tests")
	}

	// Connect to test database
	mongoURI := getEnv("MONGODB_TEST_URI", "mongodb://admin:password@localhost:27018/?authSource=admin")
	testDB := "movies_test_db"

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		t.Skipf("MongoDB not available for integration tests: %v", err)
	}
	defer client.Disconnect(context.Background())

	// Ping the database
	if err := client.Ping(ctx, nil); err != nil {
		t.Skipf("MongoDB not reachable for integration tests: %v", err)
	}

	// Clean up test database
	defer client.Database(testDB).Drop(context.Background())

	// Create repository
	repo := database.NewMongoMovieRepository(client, testDB, logger)

	t.Run("CreateAndFindMovie", func(t *testing.T) {
		// Create test movie
		movie, err := domain.NewMovie(1, "Integration Test Movie", "2023")
		if err != nil {
			t.Fatalf("Failed to create movie: %v", err)
		}

		// Save movie
		createdMovie, err := repo.Create(context.Background(), movie)
		if err != nil {
			t.Fatalf("Failed to create movie in database: %v", err)
		}

		if createdMovie.ID != movie.ID {
			t.Errorf("Created movie ID = %v, want %v", createdMovie.ID, movie.ID)
		}

		// Find movie by ID
		foundMovie, err := repo.FindByID(context.Background(), movie.ID)
		if err != nil {
			t.Fatalf("Failed to find movie by ID: %v", err)
		}

		if !foundMovie.IsEqual(movie) {
			t.Errorf("Found movie doesn't match created movie")
		}
	})

	t.Run("FindAllMovies", func(t *testing.T) {
		// Create multiple test movies
		movies := []*domain.Movie{
			{ID: 2, Title: "Movie 2", Year: "2022"},
			{ID: 3, Title: "Movie 3", Year: "2021"},
			{ID: 4, Title: "Movie 4", Year: "2020"},
		}

		for _, movie := range movies {
			if _, err := repo.Create(context.Background(), movie); err != nil {
				t.Fatalf("Failed to create test movie: %v", err)
			}
		}

		// Find all movies
		filter := domain.MovieFilter{Page: 1, Limit: 10}
		foundMovies, err := repo.FindAll(context.Background(), filter)
		if err != nil {
			t.Fatalf("Failed to find all movies: %v", err)
		}

		if len(foundMovies) < 3 {
			t.Errorf("Expected at least 3 movies, got %d", len(foundMovies))
		}
	})

	t.Run("DeleteMovie", func(t *testing.T) {
		// Create test movie
		movie := &domain.Movie{ID: 5, Title: "Movie to Delete", Year: "2023"}
		if _, err := repo.Create(context.Background(), movie); err != nil {
			t.Fatalf("Failed to create movie to delete: %v", err)
		}

		// Delete movie
		if err := repo.Delete(context.Background(), movie.ID); err != nil {
			t.Fatalf("Failed to delete movie: %v", err)
		}

		// Try to find deleted movie
		_, err := repo.FindByID(context.Background(), movie.ID)
		if err != domain.ErrMovieNotFound {
			t.Errorf("Expected ErrMovieNotFound, got %v", err)
		}
	})

	t.Run("Count", func(t *testing.T) {
		count, err := repo.Count(context.Background())
		if err != nil {
			t.Fatalf("Failed to count movies: %v", err)
		}

		if count < 0 {
			t.Errorf("Count should not be negative, got %d", count)
		}
	})
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
