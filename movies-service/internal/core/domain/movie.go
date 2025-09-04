package domain

import (
	"errors"
	"strconv"
	"time"
)

var (
	ErrMovieNotFound     = errors.New("movie not found")
	ErrInvalidMovieData  = errors.New("invalid movie data")
	ErrMovieAlreadyExists = errors.New("movie already exists")
	ErrInvalidYear       = errors.New("invalid year format")
)

type Movie struct {
	ID    int32  `json:"id" bson:"_id"`
	Title string `json:"title" bson:"title"`
	Year  string `json:"year" bson:"year"`
}

type MovieFilter struct {
	Page  int32
	Limit int32
}

// NewMovie creates a new movie with validation
func NewMovie(id int32, title, year string) (*Movie, error) {
	if title == "" {
		return nil, errors.New("title cannot be empty")
	}
	
	if year == "" {
		return nil, errors.New("year cannot be empty")
	}

	// Validate year format (should be 4 digits)
	if len(year) != 4 {
		return nil, ErrInvalidYear
	}

	if _, err := strconv.Atoi(year); err != nil {
		return nil, ErrInvalidYear
	}

	// Validate year range (1800 to current year + 10)
	currentYear := time.Now().Year()
	yearInt, _ := strconv.Atoi(year)
	if yearInt < 1800 || yearInt > currentYear+10 {
		return nil, errors.New("year must be between 1800 and current year + 10")
	}

	return &Movie{
		ID:    id,
		Title: title,
		Year:  year,
	}, nil
}

// Validate validates movie data
func (m *Movie) Validate() error {
	if m.Title == "" {
		return errors.New("title cannot be empty")
	}
	
	if m.Year == "" {
		return errors.New("year cannot be empty")
	}

	if len(m.Year) != 4 {
		return ErrInvalidYear
	}

	if _, err := strconv.Atoi(m.Year); err != nil {
		return ErrInvalidYear
	}

	return nil
}

// Update updates movie fields with validation
func (m *Movie) Update(title, year string) error {
	if title != "" {
		m.Title = title
	}
	
	if year != "" {
		if len(year) != 4 {
			return ErrInvalidYear
		}
		if _, err := strconv.Atoi(year); err != nil {
			return ErrInvalidYear
		}
		m.Year = year
	}

	return m.Validate()
}

// IsEqual checks if two movies are equal
func (m *Movie) IsEqual(other *Movie) bool {
	return m.ID == other.ID && m.Title == other.Title && m.Year == other.Year
}

// Copy creates a copy of the movie
func (m *Movie) Copy() *Movie {
	return &Movie{
		ID:    m.ID,
		Title: m.Title,
		Year:  m.Year,
	}
}