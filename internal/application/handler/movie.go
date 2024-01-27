package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) WatchMovie(c *gin.Context) {
	userId := c.Param("userId")
	movieId := c.Param("movieId")
	err := h.MovieService.WatchMovie(movieId, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *Handler) GetSuggestions(c *gin.Context) {
	userId := c.Param("userId")
	suggestions, err := h.MovieService.GetSuggestions(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, suggestions)
}
func (h *Handler) GetPaginatedMovies(c *gin.Context) {
	page := c.Query("page")
	fmt.Println("page: ", page)
	movies, newPageState := h.MovieService.GetPaginatedMovies(page)

	c.JSON(http.StatusOK, gin.H{"movies": movies, "page_state": newPageState})
}
