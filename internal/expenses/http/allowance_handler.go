package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/rrhhumand/api/internal/expenses/domain"
)

func (h *Handler) CreateAllowanceRule(c *gin.Context) {
	var req domain.DailyAllowanceRule
	if !bindJSON(c, &req) {
		return
	}
	rule, err := h.AllowanceSvc.CreateRule(c.Request.Context(), companyID(c), userID(c), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	created(c, rule)
}

func (h *Handler) ListAllowanceRules(c *gin.Context) {
	rules, err := h.AllowanceSvc.ListRules(c.Request.Context(), companyID(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, rules)
}

func (h *Handler) GetAllowanceRule(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	rule, err := h.AllowanceSvc.GetRule(c.Request.Context(), companyID(c), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	success(c, rule)
}

func (h *Handler) UpdateAllowanceRule(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req domain.DailyAllowanceRule
	if !bindJSON(c, &req) {
		return
	}
	req.ID = id
	rule, err := h.AllowanceSvc.UpdateRule(c.Request.Context(), companyID(c), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, rule)
}

func (h *Handler) CalculateAllowance(c *gin.Context) {
	var req struct {
		Destination      string   `json:"destination" binding:"required"`
		EmployeeCategory string   `json:"employee_category"`
		Days             int      `json:"days"`
		MealsProvided    []string `json:"meals_provided"`
	}
	if !bindJSON(c, &req) {
		return
	}
	dailyAmount, totalAmount, err := h.AllowanceSvc.CalculateAllowance(c.Request.Context(), companyID(c), req.Destination, req.Days, req.EmployeeCategory)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, gin.H{"daily_amount": dailyAmount, "total_amount": totalAmount})
}
