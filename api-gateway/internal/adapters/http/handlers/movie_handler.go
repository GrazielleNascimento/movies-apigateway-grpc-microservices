package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/movie-microservice/api-gateway/internal/core/domain"
	"github.com/movie-microservice/api-gateway/internal/core/ports"
)

type MovieHandler struct {
	movieService ports.MovieServicePort
	logger       *slog.Logger
}

func NewMovieHandler(movieService ports.MovieServicePort, logger *slog.Logger) *MovieHandler {
	return &MovieHandler{
		movieService: movieService,
		logger:       logger,
	}
}

func (h *MovieHandler) GetMovies(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")

	pageNum, _ := strconv.ParseInt(page, 10, 32)
	limitNum, _ := strconv.ParseInt(limit, 10, 32)

	if pageNum < 1 {
		pageNum = 1
	}
	if limitNum < 1 {
		limitNum = 10
	}

	h.logger.Info("fetching movies", "page", pageNum, "limit", limitNum)
	movies, total, err := h.movieService.GetMovies(r.Context(), int32(pageNum), int32(limitNum))
	if err != nil {
		h.logger.Error("failed to get movies", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		Movies []*domain.Movie `json:"movies"`
		Total  int32           `json:"total"`
	}{
		Movies: movies,
		Total:  total,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *MovieHandler) GetMovie(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid movie ID", http.StatusBadRequest)
		return
	}

	h.logger.Info("fetching movie", "id", id)
	movie, err := h.movieService.GetMovie(r.Context(), int32(id))
	if err != nil {
		h.logger.Error("failed to get movie", "error", err, "id", id)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movie)
}

func (h *MovieHandler) CreateMovie(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title string `json:"title"`
		Year  string `json:"year"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.logger.Error("failed to decode create movie request", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	h.logger.Info("creating movie", "title", input.Title, "year", input.Year)
	movie, err := h.movieService.CreateMovie(r.Context(), input.Title, input.Year)
	if err != nil {
		h.logger.Error("failed to create movie", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(movie)
}

func (h *MovieHandler) DeleteMovie(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		h.logger.Error("invalid movie id format", "id", idStr)
		http.Error(w, "Invalid movie ID", http.StatusBadRequest)
		return
	}

	h.logger.Info("deleting movie", "id", id)
	if err := h.movieService.DeleteMovie(r.Context(), int32(id)); err != nil {
		h.logger.Error("failed to delete movie", "error", err, "id", id)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
