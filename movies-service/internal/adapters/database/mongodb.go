package database

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/movie-microservice/movies-service/internal/core/domain"
	"github.com/movie-microservice/movies-service/internal/core/ports"
)

const (
	moviesCollection = "movies"
	defaultTimeout   = 10 * time.Second
)

type MongoMovieRepository struct {
	client   *mongo.Client
	database *mongo.Database
	logger   *slog.Logger
}

func NewMongoMovieRepository(client *mongo.Client, databaseName string, logger *slog.Logger) ports.MovieRepository {
	database := client.Database(databaseName)

	return &MongoMovieRepository{
		client:   client,
		database: database,
		logger:   logger,
	}
}

func (r *MongoMovieRepository) FindAll(ctx context.Context, filter domain.MovieFilter) ([]*domain.Movie, error) {
	collection := r.database.Collection(moviesCollection)

	// Calculate skip value
	skip := (filter.Page - 1) * filter.Limit

	// Set up options
	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(filter.Limit)).
		SetSort(bson.D{{Key: "_id", Value: 1}})

	cursor, err := collection.Find(ctx, bson.D{}, opts)
	if err != nil {
		r.logger.Error("Failed to find movies", "error", err)
		return nil, fmt.Errorf("failed to find movies: %w", err)
	}
	defer func() {
		if err := cursor.Close(ctx); err != nil {
			r.logger.Warn("Failed to close cursor", "error", err)
		}
	}()

	var movies []*domain.Movie
	if err := cursor.All(ctx, &movies); err != nil {
		r.logger.Error("Failed to decode movies", "error", err)
		return nil, fmt.Errorf("failed to decode movies: %w", err)
	}

	r.logger.Info("Successfully found movies", "count", len(movies), "page", filter.Page, "limit", filter.Limit)
	return movies, nil
}

func (r *MongoMovieRepository) FindByID(ctx context.Context, id int32) (*domain.Movie, error) {
	collection := r.database.Collection(moviesCollection)

	var movie domain.Movie
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&movie)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			r.logger.Info("Movie not found", "id", id)
			return nil, domain.ErrMovieNotFound
		}
		r.logger.Error("Failed to find movie by ID", "id", id, "error", err)
		return nil, fmt.Errorf("failed to find movie by ID: %w", err)
	}

	r.logger.Info("Successfully found movie", "id", id, "title", movie.Title)
	return &movie, nil
}

func (r *MongoMovieRepository) Create(ctx context.Context, movie *domain.Movie) (*domain.Movie, error) {
	collection := r.database.Collection(moviesCollection)

	// Validate movie before insertion
	if err := movie.Validate(); err != nil {
		return nil, fmt.Errorf("invalid movie data: %w", err)
	}

	_, err := collection.InsertOne(ctx, movie)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			r.logger.Warn("Movie with ID already exists", "id", movie.ID)
			return nil, domain.ErrMovieAlreadyExists
		}
		r.logger.Error("Failed to create movie", "movie", movie, "error", err)
		return nil, fmt.Errorf("failed to create movie: %w", err)
	}

	r.logger.Info("Successfully created movie", "id", movie.ID, "title", movie.Title)
	return movie, nil
}

func (r *MongoMovieRepository) Delete(ctx context.Context, id int32) error {
	collection := r.database.Collection(moviesCollection)

	result, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		r.logger.Error("Failed to delete movie", "id", id, "error", err)
		return fmt.Errorf("failed to delete movie: %w", err)
	}

	if result.DeletedCount == 0 {
		r.logger.Info("Movie not found for deletion", "id", id)
		return domain.ErrMovieNotFound
	}

	r.logger.Info("Successfully deleted movie", "id", id)
	return nil
}

func (r *MongoMovieRepository) Count(ctx context.Context) (int32, error) {
	collection := r.database.Collection(moviesCollection)

	count, err := collection.CountDocuments(ctx, bson.D{})
	if err != nil {
		r.logger.Error("Failed to count movies", "error", err)
		return 0, fmt.Errorf("failed to count movies: %w", err)
	}

	r.logger.Debug("Successfully counted movies", "count", count)
	return int32(count), nil
}

func (r *MongoMovieRepository) ExistsByID(ctx context.Context, id int32) (bool, error) {
	collection := r.database.Collection(moviesCollection)

	count, err := collection.CountDocuments(ctx, bson.M{"_id": id})
	if err != nil {
		r.logger.Error("Failed to check movie existence", "id", id, "error", err)
		return false, fmt.Errorf("failed to check movie existence: %w", err)
	}

	exists := count > 0
	r.logger.Debug("Checked movie existence", "id", id, "exists", exists)
	return exists, nil
}

func (r *MongoMovieRepository) GetNextID(ctx context.Context) (int32, error) {
	collection := r.database.Collection(moviesCollection)

	// Find the movie with the highest ID
	opts := options.FindOne().SetSort(bson.D{{Key: "_id", Value: -1}})
	var movie domain.Movie

	err := collection.FindOne(ctx, bson.D{}, opts).Decode(&movie)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// No movies exist, start with ID 1
			r.logger.Info("No movies found, starting with ID 1")
			return 1, nil
		}
		r.logger.Error("Failed to get max movie ID", "error", err)
		return 0, fmt.Errorf("failed to get max movie ID: %w", err)
	}

	nextID := movie.ID + 1
	r.logger.Debug("Generated next movie ID", "nextID", nextID)
	return nextID, nil
}

// Connect creates a new MongoDB connection
func Connect(ctx context.Context, connectionString string, logger *slog.Logger) (*mongo.Client, error) {
	clientOptions := options.Client().
		ApplyURI(connectionString).
		SetConnectTimeout(defaultTimeout).
		SetServerSelectionTimeout(defaultTimeout)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		logger.Error("Failed to connect to MongoDB", "error", err)
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Ping the database
	if err := client.Ping(ctx, nil); err != nil {
		logger.Error("Failed to ping MongoDB", "error", err)
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	logger.Info("Successfully connected to MongoDB")
	return client, nil
}

// Disconnect closes the MongoDB connection
func Disconnect(ctx context.Context, client *mongo.Client, logger *slog.Logger) error {
	if err := client.Disconnect(ctx); err != nil {
		logger.Error("Failed to disconnect from MongoDB", "error", err)
		return fmt.Errorf("failed to disconnect from MongoDB: %w", err)
	}

	logger.Info("Successfully disconnected from MongoDB")
	return nil
}
