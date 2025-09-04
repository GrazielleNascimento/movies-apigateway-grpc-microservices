package grpc

import (
	"context"
	
	"log/slog"

	pb "github.com/movie-microservice/proto/movies"
	"github.com/movie-microservice/movies-service/internal/core/domain"
	"github.com/movie-microservice/movies-service/internal/core/ports"
)

type MovieServer struct {
	pb.UnimplementedMovieServiceServer
	service ports.MovieService
	logger  *slog.Logger
}

func NewMovieServer(service ports.MovieService, logger *slog.Logger) *MovieServer {
	return &MovieServer{
		service: service,
		logger:  logger,
	}
}

func (s *MovieServer) GetMovies(ctx context.Context, req *pb.GetMoviesRequest) (*pb.GetMoviesResponse, error) {
	s.logger.Info("gRPC GetMovies called", "page", req.Page, "limit", req.Limit)

	filter := domain.MovieFilter{
		Page:  req.Page,
		Limit: req.Limit,
	}

	movies, total, err := s.service.GetMovies(ctx, filter)
	if err != nil {
		s.logger.Error("Failed to get movies", "error", err)
		return &pb.GetMoviesResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	// Convert domain movies to protobuf movies
	pbMovies := make([]*pb.Movie, len(movies))
	for i, movie := range movies {
		pbMovies[i] = &pb.Movie{
			Id:    movie.ID,
			Title: movie.Title,
			Year:  movie.Year,
		}
	}

	s.logger.Info("Successfully retrieved movies via gRPC", "count", len(movies))
	return &pb.GetMoviesResponse{
		Movies:  pbMovies,
		Total:   total,
		Success: true,
	}, nil
}

func (s *MovieServer) GetMovie(ctx context.Context, req *pb.GetMovieRequest) (*pb.GetMovieResponse, error) {
	s.logger.Info("gRPC GetMovie called", "id", req.Id)

	if req.Id <= 0 {
		s.logger.Warn("Invalid movie ID", "id", req.Id)
		return &pb.GetMovieResponse{
			Success: false,
			Error:   "invalid movie ID",
		}, nil
	}

	movie, err := s.service.GetMovie(ctx, req.Id)
	if err != nil {
		s.logger.Error("Failed to get movie", "id", req.Id, "error", err)
		
		if err == domain.ErrMovieNotFound {
			return &pb.GetMovieResponse{
				Success: false,
				Error:   "movie not found",
			}, nil
		}

		return &pb.GetMovieResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	s.logger.Info("Successfully retrieved movie via gRPC", "id", req.Id)
	return &pb.GetMovieResponse{
		Movie: &pb.Movie{
			Id:    movie.ID,
			Title: movie.Title,
			Year:  movie.Year,
		},
		Success: true,
	}, nil
}

func (s *MovieServer) CreateMovie(ctx context.Context, req *pb.CreateMovieRequest) (*pb.CreateMovieResponse, error) {
	s.logger.Info("gRPC CreateMovie called", "title", req.Title, "year", req.Year)

	if req.Title == "" || req.Year == "" {
		s.logger.Warn("Invalid movie data", "title", req.Title, "year", req.Year)
		return &pb.CreateMovieResponse{
			Success: false,
			Error:   "title and year are required",
		}, nil
	}

	movie, err := s.service.CreateMovie(ctx, req.Title, req.Year)
	if err != nil {
		s.logger.Error("Failed to create movie", "title", req.Title, "year", req.Year, "error", err)
		return &pb.CreateMovieResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	s.logger.Info("Successfully created movie via gRPC", "id", movie.ID)
	return &pb.CreateMovieResponse{
		Movie: &pb.Movie{
			Id:    movie.ID,
			Title: movie.Title,
			Year:  movie.Year,
		},
		Success: true,
	}, nil
}

func (s *MovieServer) DeleteMovie(ctx context.Context, req *pb.DeleteMovieRequest) (*pb.DeleteMovieResponse, error) {
	s.logger.Info("gRPC DeleteMovie called", "id", req.Id)

	if req.Id <= 0 {
		s.logger.Warn("Invalid movie ID", "id", req.Id)
		return &pb.DeleteMovieResponse{
			Success: false,
			Error:   "invalid movie ID",
		}, nil
	}

	err := s.service.DeleteMovie(ctx, req.Id)
	if err != nil {
		s.logger.Error("Failed to delete movie", "id", req.Id, "error", err)
		
		if err == domain.ErrMovieNotFound {
			return &pb.DeleteMovieResponse{
				Success: false,
				Error:   "movie not found",
			}, nil
		}

		return &pb.DeleteMovieResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	s.logger.Info("Successfully deleted movie via gRPC", "id", req.Id)
	return &pb.DeleteMovieResponse{
		Success: true,
	}, nil
}
