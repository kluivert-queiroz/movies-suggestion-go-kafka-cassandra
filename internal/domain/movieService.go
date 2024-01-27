package domain

import "movies-suggestion-service/internal/domain/entity"

type MovieService interface {
	WatchMovie(movieId string, userId string) error
	GetSuggestions(userId string) ([]entity.Suggestion, error)
	GetPaginatedMovies(page string) ([]entity.Movie, string)
}
