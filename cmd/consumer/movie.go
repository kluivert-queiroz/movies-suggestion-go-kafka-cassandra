package main

import (
	"fmt"
	"log"
	"math"
	"movies-suggestion-service/client"
	"movies-suggestion-service/internal/domain/entity"
	"movies-suggestion-service/internal/infra/repository"
	"strings"

	"github.com/IBM/sarama"
)

func main() {
	fmt.Println("Starting consumer")
	r := repository.NewMovieRepository()

	handler := &WatchedMovieConsumerHandler{MovieRepository: *r}
	err := client.ConsumeWatchedMovies(handler)
	if err != nil {
		log.Fatalf("failed to consume messages: %v", err)
	}
}

type WatchedMovieConsumerHandler struct {
	MovieRepository repository.MovieRepository
}

func (h *WatchedMovieConsumerHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (h *WatchedMovieConsumerHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (h *WatchedMovieConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	fmt.Println("Consuming messages")
	for message := range claim.Messages() {
		fmt.Println("Received message: ", message)
		args := strings.Split(string(message.Value), ":")
		userId := args[0]
		movieId := args[1]
		watchedMovie := h.MovieRepository.GetMovieById(movieId)

		// Don't recommend movies with score less than the watched movie, 'cause they are probably worse :)
		movies := filterBetterOrEqualMovies(h.MovieRepository.GetMoviesByGenres(watchedMovie.Genres), watchedMovie)

		moviesGroupedByGenre := groupMoviesByGenre(movies)
		suggestions := getSuggestions(moviesGroupedByGenre, userId)
		if err := h.MovieRepository.AddSuggestions(userId, suggestions); err != nil {
			log.Fatalf("failed to add suggestions: %v", err)
		} else {
			session.MarkMessage(message, "")
		}
		// Mark the message as processed in the session.
	}
	fmt.Println("Finished consuming messages")
	return nil
}

func getSuggestions(moviesGroupedByGenre map[string][]entity.Movie, userId string) []entity.Suggestion {
	suggestions := make([]entity.Suggestion, len(moviesGroupedByGenre))
	for genre, movies := range moviesGroupedByGenre {
		if len(movies) == 0 {
			continue
		}
		movieSuggestions := make([]entity.MovieSuggestion, len(movies))
		for _, m := range movies {
			movieSuggestions = append(movieSuggestions, entity.MovieSuggestion{
				ID:    m.ID,
				Title: m.Title,
				Score: m.Score,
			})
		}
		suggestion := entity.Suggestion{
			UserId: userId,
			Title:  "Suggestion for " + genre,
			Movies: movieSuggestions,
		}
		suggestions = append(suggestions, suggestion)
	}
	return suggestions
}

func filterBetterOrEqualMovies(movies []entity.Movie, watchedMovie *entity.Movie) []entity.Movie {
	// We'll recommend movies with score at least 2 point less than the watched movie
	// round bottom
	minimumScore := float32(math.Floor(float64(watchedMovie.Score - float32(2.0))))

	var filteredMovies []entity.Movie
	for _, m := range movies {
		if m.Score >= -minimumScore {
			filteredMovies = append(filteredMovies, m)
		}
	}
	return filteredMovies
}

func groupMoviesByGenre(movies []entity.Movie) map[string][]entity.Movie {
	moviesGroupedByGenre := make(map[string][]entity.Movie)
	for _, m := range movies {
		genres := m.Genres
		for len(genres) > 0 {
			g := genres[0]
			moviesGroupedByGenre[g] = append(moviesGroupedByGenre[m.Genres[0]], m)
			genres = genres[1:]
		}
	}
	return moviesGroupedByGenre
}
