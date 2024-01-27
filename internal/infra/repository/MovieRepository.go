package repository

import (
	"encoding/hex"
	"fmt"
	"movies-suggestion-service/internal/domain/entity"
	"movies-suggestion-service/internal/infra/db"
	"time"

	"github.com/gocql/gocql"
)

type MovieRepository struct {
	CassandraSession *gocql.Session
}

func NewMovieRepository() *MovieRepository {
	cassandra := db.NewCassandraSession()
	return &MovieRepository{CassandraSession: cassandra}
}

func (r *MovieRepository) GetMovieById(movieId string) *entity.Movie {
	var movie entity.Movie
	m := map[string]interface{}{}
	if err := r.CassandraSession.Query(`SELECT * FROM catalog.movie WHERE movieid = ? LIMIT 1`, movieId).MapScan(m); err != nil {
		fmt.Println(err)
	}
	movie.ID = m["movieid"].(string)
	movie.Title = m["title"].(string)
	movie.OriginalTitle = m["originaltitle"].(string)
	movie.TitleType = m["titletype"].(string)
	movie.StartYear = m["startyear"].(time.Time).Year()
	movie.Genres = m["genres"].([]string)
	movie.Score = m["score"].(float32)
	return &movie
}

func (r *MovieRepository) GetMoviesByGenres(genres []string) []entity.Movie {
	var movies []entity.Movie
	var movie entity.Movie
	m := map[string]interface{}{}
	for _, g := range genres {
		iter := r.CassandraSession.Query(`SELECT * FROM catalog.movie WHERE genres CONTAINS ? LIMIT 10`, g).Iter()
		for iter.MapScan(m) {
			movie.ID = m["movieid"].(string)
			movie.Title = m["title"].(string)
			movie.OriginalTitle = m["originaltitle"].(string)
			movie.TitleType = m["titletype"].(string)
			movie.StartYear = m["startyear"].(time.Time).Year()
			movie.Genres = m["genres"].([]string)
			movie.Score = m["score"].(float32)
			movies = append(movies, movie)
			m = map[string]interface{}{}
		}
	}
	return movies
}

func (r *MovieRepository) AddSuggestions(userId string, suggestions []entity.Suggestion) error {
	batch := r.CassandraSession.NewBatch(gocql.LoggedBatch)
	for _, s := range suggestions {
		batch.Query(`INSERT INTO catalog.suggestion (userid, title, movies) VALUES (?, ?, ?)`, userId, s.Title, s.Movies)
	}
	if err := r.CassandraSession.ExecuteBatch(batch); err != nil {
		return err
	}
	return nil
}

func (r *MovieRepository) GetSuggestions(userId string) []entity.Suggestion {
	var suggestions []entity.Suggestion
	var suggestion entity.Suggestion
	m := map[string]interface{}{}
	iter := r.CassandraSession.Query(`SELECT * FROM catalog.suggestion WHERE userid = ?`, userId).Iter()
	for iter.MapScan(m) {
		mList := m["movies"].([]map[string]interface{})
		var movies []entity.MovieSuggestion
		for _, s := range mList {
			movies = append(movies, entity.MovieSuggestion{
				ID:    s["id"].(string),
				Title: s["title"].(string),
				Score: s["score"].(float32),
			})
		}
		suggestion.UserId = m["userid"].(string)
		suggestion.Title = m["title"].(string)
		suggestion.Movies = movies
		suggestions = append(suggestions, suggestion)
		m = map[string]interface{}{}
	}
	return suggestions
}

func (r *MovieRepository) GetPaginatedMovies(page string) ([]entity.Movie, string) {
	hexPageState, _ := hex.DecodeString(page)
	var movies []entity.Movie
	var movie entity.Movie
	m := map[string]interface{}{}
	iter := r.CassandraSession.Query(`SELECT * FROM catalog.movie`).PageSize(50).PageState(hexPageState).Iter()
	for iter.MapScan(m) {
		movie.ID = m["movieid"].(string)
		movie.Title = m["title"].(string)
		movie.OriginalTitle = m["originaltitle"].(string)
		movie.TitleType = m["titletype"].(string)
		movie.StartYear = m["startyear"].(time.Time).Year()
		movie.Genres = m["genres"].([]string)
		movie.Score = m["score"].(float32)
		movies = append(movies, movie)
		m = map[string]interface{}{}
	}
	newPageState := hex.EncodeToString(iter.PageState())
	return movies, newPageState
}
