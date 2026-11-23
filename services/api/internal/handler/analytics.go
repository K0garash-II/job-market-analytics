package handler

import "github.com/gin-gonic/gin"

func (h *Handler) TopSkills(c *gin.Context) {
	data, err := h.market.Skills(c.Request.Context())
	respond(c, gin.H{"data": data}, err)
}
func (h *Handler) Companies(c *gin.Context) {
	data, err := h.market.Companies(c.Request.Context())
	respond(c, gin.H{"data": data}, err)
}
func (h *Handler) Overview(c *gin.Context) {
	data, err := h.market.Overview(c.Request.Context())
	respond(c, data, err)
}
func (h *Handler) SalaryStats(c *gin.Context) {
	data, err := h.market.Salary(c.Request.Context())
	respond(c, data, err)
}
