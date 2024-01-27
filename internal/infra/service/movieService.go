package service

import (
	"movies-suggestion-service/client"
	"movies-suggestion-service/internal/domain/entity"
	"movies-suggestion-service/internal/infra/repository"
)

type (
	MovieService struct {
		Client client.Client
		MovieRepository repository.MovieRepository
	}
)

func (s *MovieService) WatchMovie(movieId string, userId string) error {
	return s.Client.SendMovieWatchedMessage(movieId, userId)
}

func (s *MovieService) GetSuggestions(userId string) ([]entity.Suggestion, error) {
	return s.MovieRepository.GetSuggestions(userId), nil
}

func (s *MovieService) GetPaginatedMovies(page string) ([]entity.Movie, string) {
	return s.MovieRepository.GetPaginatedMovies(page)
}