package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/malayegor9-dotcom/job-market-analytics/services/api/internal/market"
	"net/http"
)

type Handler struct{ market market.Reader }

func New(reader market.Reader) *Handler { return &Handler{reader} }
func (h *Handler) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")
	api.GET("/vacancies", h.ListVacancies)
	api.GET("/vacancies/:id", h.GetVacancy)
	api.GET("/vacancies/search", h.SearchVacancies)
	api.GET("/analytics/skills", h.TopSkills)
	api.GET("/analytics/companies", h.Companies)
	api.GET("/analytics/salary", h.SalaryStats)
	api.GET("/analytics/overview", h.Overview)
}
func respond(c *gin.Context, data any, err error) {
	switch {
	case err == nil:
		c.JSON(http.StatusOK, data)
	case errors.Is(err, market.ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "vacancy not found"})
	default:
		_ = c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not read market data"})
	}
}
