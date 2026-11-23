package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/malayegor9-dotcom/job-market-analytics/services/api/internal/market"
	"strconv"
	"strings"
)

func (h *Handler) ListVacancies(c *gin.Context) {
	page, err := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 32)
	limit, e := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil || e != nil || page < 1 || limit < 1 || limit > 100 {
		c.JSON(400, gin.H{"error": "page must be positive and limit must be 1–100"})
		return
	}
	f := market.Filter{Page: int(page), Limit: limit, Location: c.Query("location"), Source: c.Query("source"), Query: strings.TrimSpace(c.Query("q"))}
	if raw := c.Query("remote"); raw != "" {
		v, err := strconv.ParseBool(raw)
		if err != nil {
			c.JSON(400, gin.H{"error": "invalid remote filter"})
			return
		}
		f.Remote = &v
	}
	data, err := h.market.Vacancies(c.Request.Context(), f)
	respond(c, data, err)
}
func (h *Handler) GetVacancy(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id < 1 {
		c.JSON(400, gin.H{"error": "invalid vacancy id"})
		return
	}
	data, err := h.market.Vacancy(c.Request.Context(), id)
	respond(c, data, err)
}
func (h *Handler) SearchVacancies(c *gin.Context) {
	q := strings.TrimSpace(c.Query("q"))
	if q == "" {
		c.JSON(400, gin.H{"error": "query parameter 'q' is required"})
		return
	}
	data, err := h.market.Search(c.Request.Context(), q)
	respond(c, gin.H{"data": data, "query": q}, err)
}
