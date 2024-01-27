package main

import (
	"log"
	"movies-suggestion-service/client"
	"movies-suggestion-service/internal/application/handler"
	"movies-suggestion-service/internal/infra/repository"
	"movies-suggestion-service/internal/infra/service"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	producer, err := client.SetupProducer()
	if err != nil {
		log.Fatalf("failed to initialize producer: %v", err)
	}
	MovieRepository := repository.NewMovieRepository()
	client := client.Client{KafkaSamaraProducer: producer}
	movieService := service.MovieService{Client: client, MovieRepository: *MovieRepository}
	handler := &handler.Handler{MovieService: &movieService}
	r.POST("/users/:userId/movies/:movieId/watched", handler.WatchMovie)
	r.GET("/users/:userId/suggestions", handler.GetSuggestions)
	r.GET("/movies", handler.GetPaginatedMovies)
	r.Run() // listen and serve on 0.0.0.0:8080
}
