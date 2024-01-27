package handler

import "movies-suggestion-service/internal/domain"

type (
	Handler struct {
		MovieService domain.MovieService
	}
)
